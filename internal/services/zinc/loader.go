package zinc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/pkg/config"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

// ObjectAssembly collects the objects from a bucket to load
func ObjectAssembly(v1 *viper.Viper, mc *minio.Client) error {
	bucketName, _ := config.GetBucketName(v1)

	oa, err := objects.GetObjects(v1, mc)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(oa)))

	// Single threaded loop
	//for item := range oa {
	//_, err := docfunc(v1, mc, bucketName, oa[item], "endpoint")
	//if err != nil {
	//log.Println(err)
	//}
	//bar.Add(1)
	//}

	// TODO Go func version
	semaphoreChan := make(chan struct{}, 15) // a blocking channel to keep concurrency under control
	defer close(semaphoreChan)
	wg := sync.WaitGroup{}

	log.Println("Threaded run testing")

	for item := range oa {
		wg.Add(1)
		go func(item int) {
			semaphoreChan <- struct{}{}

			_, err := docfunc(v1, mc, bucketName, oa[item], "endpoint")
			if err != nil {
				log.Println(err)
			}

			wg.Done()
			bar.Add(1)
			<-semaphoreChan // clear a spot in the semaphore channel for the next indexing event
		}(item)
	}
	wg.Wait()

	return err
}

// curl -u admin:Complexpass#123 -XPUT -d '{"name":"Prabhat Sharma"}' http://localhost:4080/api/myshinynewindex/document
func docfunc(v1 *viper.Viper, mc *minio.Client, bucketName string, item string, endpoint string) ([]byte, error) {
	jo2 := item

	b, _, err := objects.GetS3Bytes(mc, bucketName, jo2)
	if err != nil {
		return nil, err
	}

	// TODO skolemize the RDF here..

	// pulled fro the tika code which I need to review for examples
	// of modularity and also to clean up based on this new approach too.
	url := fmt.Sprintf("http://localhost:3030/testing/data?graph=urn:testing:testgraph")
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
		return nil, err
	}

	// TESTING
	//log.Printf("%s : %d  : %s\n", jo2, len(b), endpoint)

	return nil, nil
}
