package flows

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/UFOKN/nabu/internal/graph"
	"github.com/UFOKN/nabu/internal/objects"
	"github.com/schollz/progressbar/v3"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

// ObjectAssembly collects the objects from a bucket to load
func ObjectAssembly(v1 *viper.Viper, mc *minio.Client) error {
	objs := v1.GetStringMapString("objects")
	spql := v1.GetStringMapString("sparql")

	// My go func controller vars
	semaphoreChan := make(chan struct{}, 1) // a blocking channel to keep concurrency under control (1 == single thread)
	defer close(semaphoreChan)
	// wg := sync.WaitGroup{} // a wait group enables the main process a wait for goroutines to finish

	// params for list objects calls
	// doneCh := make(chan struct{}) // , N) Create a done channel to control 'ListObjectsV2' go routine.
	// defer close(doneCh)           // Indicate to our routine to exit cleanly upon return.

	var pa []string
	err := v1.UnmarshalKey("objects.prefix", &pa)
	if err != nil {
		log.Println(err)
	}
<<<<<<< HEAD

	fmt.Println(pa)

	for p := range pa {
		oa := []string{}

		// NEW
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		objectCh := mc.ListObjects(ctx, objs["bucket"],
			minio.ListObjectsOptions{Prefix: pa[p], Recursive: true})

		for object := range objectCh {
			if object.Err != nil {
				fmt.Println(object.Err)
				return object.Err
			}
			//fmt.Println(object.Key)
			oa = append(oa, object.Key)
		}

=======

	fmt.Println(pa)

	for p := range pa {
		oa := []string{}

		// NEW
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		objectCh := mc.ListObjects(ctx, objs["bucket"],
			minio.ListObjectsOptions{Prefix: pa[p], Recursive: true})

		for object := range objectCh {
			if object.Err != nil {
				fmt.Println(object.Err)
				return object.Err
			}
			// fmt.Println(object)
			oa = append(oa, object.Key)
		}

>>>>>>> multiprefix
		log.Printf("%s:%s object count: %d\n", objs["bucket"], pa[p], len(oa))
		bar := progressbar.Default(int64(len(oa)))
		for item := range oa {
			_, err := PipeLoad(v1, mc, objs["bucket"], oa[item], spql["endpoint"])
			if err != nil {
				log.Println(err)
			}
			bar.Add(1)
			// log.Println(string(s)) // get "s" on pipeload and send to a log file
		}
	}

	return nil
}

// PipeLoad reads from an object and loads directly into a triplestore
func PipeLoad(v1 *viper.Viper, mc *minio.Client, bucket, object, spql string) ([]byte, error) {
	// build our quad/graph from the object path
	log.Printf("Loading %s \n", object)
	s2c := strings.Replace(object, "/", ":", -1)

	var g string
	if strings.Contains(s2c, ".rdf") {
		g = fmt.Sprintf("urn:%s:%s", bucket, strings.TrimSuffix(s2c, ".rdf"))
	} else if strings.Contains(s2c, ".jsonld") {
		g = fmt.Sprintf("urn:%s:%s", bucket, strings.TrimSuffix(s2c, ".jsonld"))
	} else if strings.Contains(s2c, ".nq") {
		g = fmt.Sprintf("urn:%s:%s", bucket, strings.TrimSuffix(s2c, ".nq"))
	} else {
		return nil, errors.New("unable to generate graph URI")
	}

	log.Println(g)

	// Turn checking off while testing other parts of Nabu
	//c, err := gexists(spql, g)
	//if err != nil {
	//log.Println(err)
	//}
	//if c {
	//return nil, nil // our graph is loaded already..
	//}

	b, _, err := objects.GetS3Bytes(mc, bucket, object)
	if err != nil {
		log.Printf("gets3Bytes %v\n", err)
	}

	// check the object string

	// log.Println(mt) // application/ld+json

	// TODO, use the mimetype or suffix in general to select the
	// path to load    or overload from the config file?

	mt := mime.TypeByExtension(filepath.Ext(object))
	nt := ""

	if strings.Compare(mt, "application/ld+json") == 0 {
		log.Println("Convert JSON-LD file to nq")
		nt, err = graph.JSONLDToNQ(string(b))
		if err != nil {
			log.Printf("JSONLDToNQ err: %s", err)
		}
	} else {
		nt, _, err = graph.NQToNTCtx(string(b))
		if err != nil {
			log.Printf("nqToNTCtx err: %s", err)
		}
	}

	// drop any graph we are going to load..  we assume we are doing those due to an update...
	_, err = Drop(v1, g)
	if err != nil {
		log.Println(err)
	}

	// If the graph is a quad already..   we need to make it triples
	// so we can load with "our" context.
	// Note: We are tossing source prov for out prov

	log.Printf("Graph loading as: %s\n", g)

	// TODO if array is too large, need to split it and load parts
	// Let's declare 10k lines the largest we want to send in.
	log.Printf("Graph size: %d\n", len(nt))

	scanner := bufio.NewScanner(strings.NewReader(nt))
	lc := 0
	sg := []string{}
	for scanner.Scan() {
		lc = lc + 1
		sg = append(sg, scanner.Text())
		if lc == 10000 { // use line count, since byte len might break inside a triple statement..   it's an OK proxy
			log.Printf("Subgraph of %d lines", len(sg))
			// TODO..  upload what we have here, modify the call code to upload these sections
			_, err = Insert(g, strings.Join(sg, "\n"), spql) // convert []string to strings joined with new line to form a RDF NT set
			if err != nil {
				log.Printf("Insert err: %s", err)
			}
			sg = nil // clear the array
			lc = 0   // reset the counter
		}
	}
	if lc > 0 {
		log.Printf("Subgraph (out of scanner) of %d lines", len(sg))
		_, err = Insert(g, strings.Join(sg, "\n"), spql) // convert []string to strings joined with new line to form a RDF NT set
	}

	return []byte("remove me"), err
}

func Insert(g, nt, spql string) (string, error) {

	p := "INSERT DATA { "
	pab := []byte(p)
	gab := []byte(fmt.Sprintf(" graph <%s>  { ", g))
	u := " } }"
	uab := []byte(u)
	pab = append(pab, gab...)
	pab = append(pab, []byte(nt)...)
	pab = append(pab, uab...)

	req, err := http.NewRequest("POST", spql, bytes.NewBuffer(pab))
	if err != nil {
		log.Println(err)
	}
	// req.Header.Set("Content-Type", "application/sparql-results+xml")
	req.Header.Set("Content-Type", "application/sparql-update")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	// log.Println(string(body))
	if err != nil {
		log.Println("response Body:", string(body))
		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
	}

	return resp.Status, err
}

// Drop removes a graph
func Drop(v1 *viper.Viper, g string) ([]byte, error) {
	spql := v1.GetStringMapString("sparql")
	// d := fmt.Sprintf("DELETE { GRAPH <%s> {?s ?p ?o} } WHERE {GRAPH <%s> {?s ?p ?o}}", g, g)
	d := fmt.Sprintf("DROP GRAPH <%s> ", g)

	pab := []byte(d)

	req, err := http.NewRequest("POST", spql["endpoint"], bytes.NewBuffer(pab))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/sparql-update")
	// req.Header.Set("Content-Type", "application/sparql-results+xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("response Body:", string(body))
		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
	}

	// log.Println(string(body))

	return body, err
}

// DropGet removes a graph
func DropGet(v1 *viper.Viper, g string) ([]byte, error) {
	spql := v1.GetStringMapString("sparql")
	// d := fmt.Sprintf("DELETE { GRAPH <%s> {?s ?p ?o} } WHERE {GRAPH <%s> {?s ?p ?o}}", g, g)
	d := fmt.Sprintf("DROP GRAPH <%s> ", g)

	log.Println(d)
	pab := []byte("")
	//	fmt.Println(spql["endpoint"])

	// TODO try and GET with query set
	params := url.Values{}
	params.Add("query", d)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", spql["endpoint"], params.Encode()), bytes.NewBuffer(pab))
	if err != nil {
		log.Println(err)
	}
	// req.Header.Set("Content-Type", "application/sparql-results+json")
	req.Header.Set("Content-Type", "application/sparql-update")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("response Body:", string(body))
		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
	}

	log.Println(string(body))

	return body, err
}

// OLD
// for object := range mc.ListObjects(context.Background(), objs["bucket"],
// 	minio.ListObjectsOptions{Prefix: objs["prefix"], Recursive: true}) {

// 	wg.Add(1)
// 	go func(object minio.ObjectInfo) {
// 		// TEMP HACK for UFOKN
// 		// if objs["objectsuffix"] != "" {
// 		// 	if strings.HasSuffix(object.Key, objs["objectsuffix"]) {
// 		// 		// log.Println(object.Key)
// 		// 		oa = append(oa, object.Key) // WARNING  append is not always thread safe..   wg of 1 till I address this
// 		// 	}
// 		// } else {
// 		// log.Println(object.Key)

// 		oa = append(oa, object.Key) // WARNING  append is not always thread safe..   wg of 1 till I address this
// 		// }

// 		wg.Done() // tell the wait group that we be done
// 		// log.Printf("Doc: %s error: %v ", name, err) // why print the status??
// 		<-semaphoreChan
// 	}(object)
// 	wg.Wait()
// }
