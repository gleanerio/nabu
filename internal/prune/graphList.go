package prune

import (
	"bytes"
	"fmt"
	"github.com/gleanerio/nabu/pkg/config"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func graphList(v1 *viper.Viper, mc *minio.Client, prefix string) ([]string, error) {
	ga := []string{}

	spql, _ := config.GetSparqlConfig(v1)
	bucketName, _ := config.GetBucketName(v1)

	// Reference ADR 0001 for why we are building the regex here for the graphs like this.
	s2c := strings.Replace(prefix, "summoned/", ":", -1)
	gp := fmt.Sprintf("urn:%s%s:", bucketName, s2c)
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
