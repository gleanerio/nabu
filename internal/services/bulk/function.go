package bulk

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gleanerio/nabu/internal/graph"

	"github.com/gleanerio/nabu/internal/objects"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/minio/minio-go/v7"
)

func docfunc(v1 *viper.Viper, mc *minio.Client, bucketName string, item string, endpoint string) (string, error) {
	log.Printf("Jena docfunc called with %s%s", bucketName, item)
	log.Println(endpoint)

	b, _, err := objects.GetS3Bytes(mc, bucketName, item)
	if err != nil {
		return "", err
	}

	// TODO skolemize the RDF here..
	// unless bulk loading, in which case it needs to be done prior to here and this should be skipped
	// the "bulk" load function might be different too

	bn := strings.Replace(bucketName, ".", ":", -1) //why is this here?
	g, err := graph.MakeURN(item, bn)
	if err != nil {
		log.Error("gets3Bytes %v\n", err)
		// should this just return. since on this error things are not good
	}

	// check if JSON-LD and convert to RDF
	if strings.Contains(item, ".jsonld") {
		nb, err := graph.JSONLDToNQ(v1, string(b))
		if err != nil {
			return "", err
		}
		b = []byte(nb)
	}

	url := fmt.Sprintf("%s?graph=%s", endpoint, g)

	req, err := http.NewRequest("PUT", url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/n-quads")
	req.Header.Set("User-Agent", "EarthCube_DataBot/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	log.Println(resp)
	body, err := ioutil.ReadAll(resp.Body) // return body if you want to debugg test with it
	if err != nil {
		log.Println(string(body))
		return string(body), err
	}

	// TESTING
	log.Println(string(body))
	log.Printf("success: %s : %d  : %s\n", item, len(b), endpoint)

	return string(body), err
}
