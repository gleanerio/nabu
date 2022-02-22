package prune

import (
	"bytes"
	"fmt"
	"github.com/gleanerio/nabu/internal/flows"
	"github.com/gleanerio/nabu/pkg/config"
	"github.com/minio/minio-go/v7"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//Snip removes graphs in TS not in object store
func Snip(v1 *viper.Viper, mc *minio.Client) error {

	var pa []string
	//err := v1.UnmarshalKey("objects.prefix", &pa)
	objs, err := config.GetObjectsConfig(v1)
	bucketName, _ := config.GetBucketName(v1)
	if err != nil {
		log.Println(err)
	}
	pa = objs.Prefix

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

		//objs := v1.GetStringMapString("objects") // from above
		// convert the object names to the URN pattern used in the graph
		for x := range oa {
			s := strings.TrimSuffix(oa[x], ".rdf")
			s2 := strings.Replace(s, "/", ":", -1)
			g := fmt.Sprintf("urn:%s:%s", bucketName, s2)
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

	//spql := v1.GetStringMapString("sparql")
	//objs := v1.GetStringMapString("objects")
	spql, _ := config.GetSparqlConfig(v1)
	//objs,_ := config.GetObjectsConfig(v1)
	bucketName, _ := config.GetBucketName(v1)
	//gp := fmt.Sprintf("urn:%s:%s", objs["bucket"], strings.Replace(prefix, "/", ":", -1))
	gp := fmt.Sprintf("urn:%s:%s", bucketName, strings.Replace(prefix, "/", ":", -1))
	fmt.Printf("Pattern: %s\n", gp)

	d := fmt.Sprintf("SELECT DISTINCT ?g WHERE {GRAPH ?g {?s ?p ?o} FILTER regex(str(?g), \"^%s\")}", gp)

	//fmt.Println(d)

	pab := []byte("")
	params := url.Values{}
	params.Add("query", d)
	//req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", spql["endpoint"], params.Encode()), bytes.NewBuffer(pab))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", spql.Endpoint, params.Encode()), bytes.NewBuffer(pab))
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

//
//func graphListStatements(v1 *viper.Viper, mc *minio.Client, prefix string) ([]string, error) {
//
//	ga := []string{}
//
//	//spql := v1.GetStringMapString("sparql")
//	//objs := v1.GetStringMapString("objects")
//	spql, _ := config.GetSparqlConfig(v1)
//	//objs,_ := config.GetObjectsConfig(v1)
//	bucketName, _ := config.GetBucketName(v1)
//
//	gp := fmt.Sprintf("urn:%s:%s", bucketName, strings.Replace(prefix, "/", ":", -1))
//	fmt.Printf("Pattern: %s\n", gp)
//
//	d := fmt.Sprintf("SELECT DISTINCT ?g WHERE {GRAPH ?g {?s ?p ?o} FILTER regex(str(?g), \"^%s\")}", gp)
//
//	fmt.Println(d)
//
//	pab := []byte("")
//	params := url.Values{}
//	params.Add("query", d)
//	//req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", spql["endpoint"], params.Encode()), bytes.NewBuffer(pab))
//	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", spql.Endpoint, params.Encode()), bytes.NewBuffer(pab))
//	if err != nil {
//		log.Println(err)
//	}
//	// req.Header.Add("Accept", "application/sparql-update")
//	req.Header.Add("Accept", "application/n-quads")
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Println(err)
//	}
//
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Println(strings.Repeat("ERROR", 5))
//		log.Println("response Status:", resp.Status)
//		log.Println("response Headers:", resp.Header)
//		log.Println("response Body:", string(body))
//	}
//
//	// fmt.Println("response Body:", string(body))
//
//	result := gjson.Get(string(body), "results.bindings.#.g.value")
//	result.ForEach(func(key, value gjson.Result) bool {
//		ga = append(ga, value.String())
//		return true // keep iterating
//	})
//
//	// ask := Ask{}
//	// json.Unmarshal(body, &ask)
//	return ga, nil
//}
