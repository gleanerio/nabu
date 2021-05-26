package prune

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/UFOKN/nabu/internal/flows"
	"github.com/minio/minio-go/v7"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

//Snip removes graphs in TS not in object store
func Snip(v1 *viper.Viper, mc *minio.Client) error {

	var pa []string
	err := v1.UnmarshalKey("objects.prefix", &pa)
	if err != nil {
		log.Println(err)
	}

	for p := range pa {

		// do the object assembly
		oa, err := ObjectList(v1, mc, pa[p])
		if err != nil {
			log.Println(err)
			return err
		}

		// collect all the graphs from triple store
		// TODO resolve issue with Graph and graphList vs graphListStatements
		//ga, err := graphListStatements(v1, mc, pa[p])
		ga, err := graphList(v1, mc, pa[p])
		if err != nil {
			log.Println(err)
			return err
		}

		objs := v1.GetStringMapString("objects")
		// convert the object names to the URN pattern used in the graph
		for x := range oa {
			s := strings.TrimSuffix(oa[x], ".rdf")
			s2 := strings.Replace(s, "/", ":", -1)
			g := fmt.Sprintf("urn:%s:%s", objs["bucket"], s2)
			oa[x] = g
		}

		//compare lists..   anything IN graph not in objects list should be removed
		d := difference(ga, oa) // return array of items in ga that are NOT in oa

		fmt.Printf("Graph items: %d  Object items: %d  difference: %d\n", len(ga), len(oa), len(d))

		// For each in d will delete that graph
		bar := progressbar.Default(int64(len(d)))
		for x := range d {
			log.Printf("Remove graph: %s\n", d[x])
			flows.Drop(v1, d[x])
			bar.Add(1)
		}
	}

	return nil
}

func graphList(v1 *viper.Viper, mc *minio.Client, prefix string) ([]string, error) {
	ga := []string{}

	spql := v1.GetStringMapString("sparql")
	objs := v1.GetStringMapString("objects")

	gp := fmt.Sprintf("urn:%s:%s", objs["bucket"], strings.Replace(prefix, "/", ":", -1))
	fmt.Printf("Pattern: %s\n", gp)

	d := fmt.Sprintf("SELECT DISTINCT ?g WHERE {GRAPH ?g {?s ?p ?o} FILTER regex(str(?g), \"^%s\")}", gp)

	//fmt.Println(d)

	pab := []byte("")
	params := url.Values{}
	params.Add("query", d)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", spql["endpoint"], params.Encode()), bytes.NewBuffer(pab))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Accept", "application/sparql-results+json")

	//req.Header.Add("Accept", "application/sparql-update")
	//req.Header.Add("Accept", "application/n-quads")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(strings.Repeat("ERROR", 5))
		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
		log.Println("response Body:", string(body))
	}

	//fmt.Println("response Body:", string(body))

	result := gjson.Get(string(body), "results.bindings.#.g.value")
	result.ForEach(func(key, value gjson.Result) bool {
		ga = append(ga, value.String())
		return true // keep iterating
	})

	// ask := Ask{}
	// json.Unmarshal(body, &ask)
	return ga, nil
}

func graphListStatements(v1 *viper.Viper, mc *minio.Client, prefix string) ([]string, error) {
	ga := []string{}

	spql := v1.GetStringMapString("sparql")
	objs := v1.GetStringMapString("objects")

	gp := fmt.Sprintf("urn:%s:%s", objs["bucket"], strings.Replace(prefix, "/", ":", -1))
	fmt.Printf("Pattern: %s\n", gp)

	d := fmt.Sprintf("SELECT DISTINCT ?g WHERE {GRAPH ?g {?s ?p ?o} FILTER regex(str(?g), \"^%s\")}", gp)

	fmt.Println(d)

	pab := []byte("")
	params := url.Values{}
	params.Add("query", d)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", spql["endpoint"], params.Encode()), bytes.NewBuffer(pab))
	if err != nil {
		log.Println(err)
	}
	// req.Header.Add("Accept", "application/sparql-update")
	req.Header.Add("Accept", "application/n-quads")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(strings.Repeat("ERROR", 5))
		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
		log.Println("response Body:", string(body))
	}

	// fmt.Println("response Body:", string(body))

	result := gjson.Get(string(body), "results.bindings.#.g.value")
	result.ForEach(func(key, value gjson.Result) bool {
		ga = append(ga, value.String())
		return true // keep iterating
	})

	// ask := Ask{}
	// json.Unmarshal(body, &ask)
	return ga, nil
}

func ObjectList(v1 *viper.Viper, mc *minio.Client, prefix string) ([]string, error) {
	objs := v1.GetStringMapString("objects")

	// My go func controller vars
	semaphoreChan := make(chan struct{}, 1) // a blocking channel to keep concurrency under control (1 == single thread)
	defer close(semaphoreChan)
	wg := sync.WaitGroup{} // a wait group enables the main process a wait for goroutines to finish

	// params for list objects calls
	doneCh := make(chan struct{}) // , N) Create a done channel to control 'ListObjectsV2' go routine.
	defer close(doneCh)           // Indicate to our routine to exit cleanly upon return.
	// isRecursive := true

	oa := []string{}

	// for object := range mc.ListObjectsV2(objs["bucket"], objs["prefix"], isRecursive, doneCh) {
	for object := range mc.ListObjects(context.Background(), objs["bucket"],
		minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {

		wg.Add(1)
		go func(object minio.ObjectInfo) {
			oa = append(oa, object.Key) // WARNING  append is not always thread safe..   wg of 1 till I address this
			wg.Done()                   // tell the wait group that we be done
			// log.Printf("Doc: %s error: %v ", name, err) // why print the status??
			<-semaphoreChan
		}(object)
		wg.Wait()
	}

	log.Printf("%s:%s object count: %d\n", objs["bucket"], prefix, len(oa))

	return oa, nil
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
