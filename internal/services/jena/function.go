package jena

import (
	"bytes"
	"fmt"
	"github.com/gleanerio/nabu/internal/objects"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"

	"github.com/minio/minio-go/v7"
)

// curl -u admin:Complexpass#123 -XPUT -d '{"name":"Prabhat Sharma"}' http://localhost:4080/api/myshinynewindex/document
func docfunc(v1 *viper.Viper, mc *minio.Client, bucketName string, item string, endpoint string) ([]byte, error) {
	jo2 := item

	b, _, err := objects.GetS3Bytes(mc, bucketName, jo2)
	if err != nil {
		return nil, err
	}

	log.Printf("Jena call with %s/%s", bucketName, item)

	// TODO skolemize the RDF here..
	// unless bulk loading, in which case it needs to be done prior to here and this should be skipped
	// the "bulk" load function might be different too

	url := fmt.Sprintf("http://localhost:3030/testing/data?graph=urn:testing:testgraph2")
	req, err := http.NewRequest("PUT", url, bytes.NewReader(b))
	//req.Header.Set("Accept", "application/n-quads")
	req.Header.Set("Content-Type", "application/n-quads")
	req.Header.Set("User-Agent", "EarthCube_DataBot/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) // return body if you want to debugg test with it
	if err != nil {
		log.Println(string(body))
		return body, err
	}

	// TESTING
	//log.Printf("%s : %d  : %s\n", jo2, len(b), endpoint)

	return body, nil
}
