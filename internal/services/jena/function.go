package jena

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gleanerio/nabu/internal/objects"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/minio/minio-go/v7"
)

// curl -u admin:Complexpass#123 -XPUT -d '{"name":"Prabhat Sharma"}' http://localhost:4080/api/myshinynewindex/document
func docfunc(v1 *viper.Viper, mc *minio.Client, bucketName string, item string, endpoint string) (string, error) {
	jo2 := item

	b, _, err := objects.GetS3Bytes(mc, bucketName, jo2)
	if err != nil {
		return "", err
	}

	log.Printf("Jena call with %s/%s", bucketName, item)

	// TODO skolemize the RDF here..
	// unless bulk loading, in which case it needs to be done prior to here and this should be skipped
	// the "bulk" load function might be different too

	s2c := strings.Replace(item, "/", ":", -1)

	// build the URN for the graph context string we use
	var g string
	if strings.Contains(s2c, ".rdf") {
		g = fmt.Sprintf("urn:%s:%s", bucketName, strings.TrimSuffix(s2c, ".rdf"))
	} else if strings.Contains(s2c, ".jsonld") {
		g = fmt.Sprintf("urn:%s:%s", bucketName, strings.TrimSuffix(s2c, ".jsonld"))
	} else if strings.Contains(s2c, ".nq") {
		g = fmt.Sprintf("urn:%s:%s", bucketName, strings.TrimSuffix(s2c, ".nq"))
	} else {
		return "", errors.New("unable to generate graph URI")
	}

	url := fmt.Sprintf("http://localhost:3030/oih/data?graph=%s", g)
	log.Printf("Loding to %s\n", url)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(b))
	//req.Header.Set("Accept", "application/n-quads")
	req.Header.Set("Content-Type", "application/n-quads")
	req.Header.Set("User-Agent", "EarthCube_DataBot/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) // return body if you want to debugg test with it
	if err != nil {
		log.Println(string(body))
		return string(body), err
	}

	// TESTING
	//log.Printf("%s : %d  : %s\n", jo2, len(b), endpoint)

	return string(body), err
}
