package flows

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/UFOKN/nabu/internal/graph"
	"github.com/UFOKN/nabu/internal/objects"
	"github.com/schollz/progressbar/v3"

	"github.com/minio/minio-go"
	"github.com/spf13/viper"
)

// ObjectAssembly collects the objects from a bucket to load
func ObjectAssembly(v1 *viper.Viper, mc *minio.Client) error {

	objs := v1.GetStringMapString("objects")
	spql := v1.GetStringMapString("sparql")

	// My go func controller vars
	semaphoreChan := make(chan struct{}, 1) // a blocking channel to keep concurrency under control (1 == single thread)
	defer close(semaphoreChan)
	wg := sync.WaitGroup{} // a wait group enables the main process a wait for goroutines to finish

	// params for list objects calls
	doneCh := make(chan struct{}) // , N) Create a done channel to control 'ListObjectsV2' go routine.
	defer close(doneCh)           // Indicate to our routine to exit cleanly upon return.
	isRecursive := true

	oa := []string{}

	for object := range mc.ListObjectsV2(objs["bucket"], objs["prefix"], isRecursive, doneCh) {
		wg.Add(1)
		go func(object minio.ObjectInfo) {
			oa = append(oa, object.Key) // WARNING  append is not always thread safe..   wg of 1 till I address this
			wg.Done()                   // tell the wait group that we be done
			// log.Printf("Doc: %s error: %v ", name, err) // why print the status??
			<-semaphoreChan
		}(object)
		wg.Wait()
	}

	log.Printf("%s:%s object count: %d\n", objs["bucket"], objs["prefix"], len(oa))
	bar := progressbar.Default(int64(len(oa)))
	for item := range oa {
		_, err := PipeLoad(v1, mc, objs["bucket"], oa[item], spql["endpoint"])
		if err != nil {
			log.Println(err)
		}
		bar.Add(1)
		// log.Println(string(s)) // get "s" on pipeload and send to a log file
	}

	return nil
}

// PipeLoad reads from an object and loads directly into a triplestore
func PipeLoad(v1 *viper.Viper, mc *minio.Client, bucket, object, spql string) ([]byte, error) {
	// build our quad/graph from the object path

	s2c := strings.Replace(object, "/", ":", -1)
	g := fmt.Sprintf("urn:%s:%s", bucket, strings.TrimSuffix(s2c, ".rdf"))

	// Turn checking off while testing other parts of Nabu
	c, err := gexists(spql, g)
	if err != nil {
		log.Println(err)
	}
	if c {
		return nil, nil // our graph is loaded already..
	}

	b, _, err := objects.GetS3Bytes(mc, bucket, object)
	if err != nil {
		fmt.Printf("gets3Bytes %v\n", err)
	}

	nt, _, err := graph.NQToNTCtx(string(b))
	if err != nil {
		log.Printf("nqToNTCtx err: %s", err)
	}

	// drop any graph we are going to load..  we assume we are doing those due to an update...
	_, err = Drop(v1, g)
	if err != nil {
		log.Println(err)
	}
	// If the graph is a quad already..   we need to make it triples
	// so we can load with "our" context.
	// Note: We are tossing source prov for out prov

	// log.Printf("Graph loading as: %s\n", g)

	p := "INSERT DATA { "
	pab := []byte(p)
	gab := []byte(fmt.Sprintf(" graph <%s>  { ", g))
	u := " } }"
	uab := []byte(u)
	pab = append(pab, gab...)
	pab = append(pab, []byte(nt)...)
	pab = append(pab, uab...)

	req, err := http.NewRequest("POST", spql, bytes.NewBuffer(pab))
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

	return []byte(resp.Status), err
}

// Drop removes a graph
func Drop(v1 *viper.Viper, g string) ([]byte, error) {

	spql := v1.GetStringMapString("sparql")

	d := fmt.Sprintf("DELETE { GRAPH <%s> {?s ?p ?o} } WHERE {GRAPH <%s> {?s ?p ?o}}", g, g)

	pab := []byte(d)

	req, err := http.NewRequest("POST", spql["endpoint"], bytes.NewBuffer(pab))
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

	return body, err
}
