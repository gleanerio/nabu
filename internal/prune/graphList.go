package prune

import (
	"bytes"
	"fmt"
	"github.com/gleanerio/nabu/internal/graph"
	"github.com/gleanerio/nabu/pkg/config"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func graphList(v1 *viper.Viper, mc *minio.Client, prefix string) ([]string, error) {
	ga := []string{}

	spql, err := config.GetSparqlConfig(v1)
	if err != nil {
		log.Println(err)
		return ga, err
	}
	bucketName, err := config.GetBucketName(v1)
	if err != nil {
		log.Println(err)
		return ga, err
	}

	gp, err := graph.MakeURNPrefix(v1, bucketName, prefix)
	if err != nil {
		log.Println(err)
		return ga, err
	}

	log.Printf("Pattern: %s\n", gp)

	d := fmt.Sprintf("SELECT DISTINCT ?g WHERE {GRAPH ?g {?s ?p ?o} FILTER regex(str(?g), \"^%s\")}", gp)
	//d := fmt.Sprintf("SELECT DISTINCT ?g WHERE {GRAPH ?g {?s ?p ?o} FILTER regex(str(?g), \"^%s\")}", gp)

	log.Printf("SPARQL: %s\n", d)

	pab := []byte("")
	params := url.Values{}
	params.Add("query", d)
	//req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", spql["endpoint"], params.Encode()), bytes.NewBuffer(pab))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", spql.Endpoint, params.Encode()), bytes.NewBuffer(pab))
	if err != nil {
		log.Println(err)
	}

	// These headers
	req.Header.Set("Accept", "application/sparql-results+json")
	//req.Header.Add("Accept", "application/sparql-update")
	//req.Header.Add("Accept", "application/n-quads")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(strings.Repeat("ERROR", 5))
		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
		log.Println("response Body:", string(body))
	}

	//fmt.Println("response Body:", string(body))
	err = ioutil.WriteFile("myfile.txt", body, 0644)
	if err != nil {
		fmt.Println("An error occurred:", err)
		return ga, err
	}

	result := gjson.Get(string(body), "results.bindings.#.g.value")
	result.ForEach(func(key, value gjson.Result) bool {
		ga = append(ga, value.String())
		return true // keep iterating
	})

	log.Println(len(ga))

	return ga, nil
}
