package flows

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gleanerio/nabu/run/config"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gleanerio/nabu/internal/graph"
	"github.com/gleanerio/nabu/internal/objects"
	"github.com/schollz/progressbar/v3"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

// ObjectAssembly collects the objects from a bucket to load
func ObjectAssembly(v1 *viper.Viper, mc *minio.Client) error {
	//objs := v1.GetStringMapString("objects")
	//spql := v1.GetStringMapString("sparql")
	objs, err := config.GetObjectsConfig(v1)
	spql, err := config.GetSparqlConfig(v1)

	var pa = objs.Prefix
	//var pa []string
	//err := v1.UnmarshalKey("objects.prefix", &pa)
	//if err != nil {
	//	log.Println(err)
	//}

	log.Println(pa)

	for p := range pa {
		oa := []string{}

		// NEW
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		bucketName, _ := config.GetBucketName(v1)
		objectCh := mc.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Prefix: pa[p], Recursive: true})

		for object := range objectCh {
			if object.Err != nil {
				log.Println(object.Err)
				return object.Err
			}
			// fmt.Println(object)
			oa = append(oa, object.Key)
		}

		log.Printf("%s:%s object count: %d\n", bucketName, pa[p], len(oa))
		bar := progressbar.Default(int64(len(oa)))
		for item := range oa {
			_, err := PipeLoad(v1, mc, bucketName, oa[item], spql.Endpoint)
			if err != nil {
				log.Println(err)
			}
			bar.Add(1)
			// log.Println(string(s)) // get "s" on pipeload and send to a log file
		}
	}

	return err
}

// PipeLoad reads from an object and loads directly into a triplestore
func PipeLoad(v1 *viper.Viper, mc *minio.Client, bucket, object, spql string) ([]byte, error) {
	// build our quad/graph from the object path
	log.Printf("Loading %s \n", object)
	s2c := strings.Replace(object, "/", ":", -1)

	// build the URN for the graph context string we use
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

	// TODO WARNING this needs to be addressed
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

	// TODO, use the mimetype or suffix in general to select the path to load    or overload from the config file?
	// check the object string
	mt := mime.TypeByExtension(filepath.Ext(object))
	log.Printf("Object: %s reads as mimetype: %s", object, mt) // application/ld+json
	nt := ""

	// if strings.Contains(object, ".jsonld") { // TODO explore why this hack is needed and the mimetype for JSON-LD is not returned
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
	//spql := v1.GetStringMapString("sparql")
	spql, _ := config.GetSparqlConfig(v1)
	// d := fmt.Sprintf("DELETE { GRAPH <%s> {?s ?p ?o} } WHERE {GRAPH <%s> {?s ?p ?o}}", g, g)
	d := fmt.Sprintf("DROP GRAPH <%s> ", g)

	pab := []byte(d)

	//req, err := http.NewRequest("POST", spql["endpoint"], bytes.NewBuffer(pab))
	req, err := http.NewRequest("POST", spql.Endpoint, bytes.NewBuffer(pab))
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
	//spql := v1.GetStringMapString("sparql")
	spql, _ := config.GetSparqlConfig(v1)
	// d := fmt.Sprintf("DELETE { GRAPH <%s> {?s ?p ?o} } WHERE {GRAPH <%s> {?s ?p ?o}}", g, g)
	d := fmt.Sprintf("DROP GRAPH <%s> ", g)

	log.Println(d)
	pab := []byte("")
	//	fmt.Println(spql["endpoint"])

	// TODO try and GET with query set
	params := url.Values{}
	params.Add("query", d)
	//req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", spql["endpoint"], params.Encode()), bytes.NewBuffer(pab))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", spql.Endpoint, params.Encode()), bytes.NewBuffer(pab))
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
