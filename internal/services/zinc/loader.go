package zinc

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/pkg/config"
	"github.com/schollz/progressbar/v3"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

// ObjectAssembly collects the objects from a bucket to load
func ObjectAssembly(v1 *viper.Viper, mc *minio.Client) error {
	// log to custom file
	LOG_FILE := "zinc_log" // open log file
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
		return err
	}
	defer logFile.Close()

	log.SetOutput(logFile) // Set log out put and enjoy :)

	log.SetFlags(log.Lshortfile | log.LstdFlags) // optional: log date-time, filename, and line number
	log.Println("Logging to custom file")

	bucketName, _ := config.GetBucketName(v1)

	oa, err := GetObjects(v1, mc)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(oa)))
	for item := range oa {
		_, err := docfunc(v1, mc, bucketName, oa[item], "endpoint")
		if err != nil {
			log.Println(err)
		}
		bar.Add(1)
	}

	return err
}

// GetObjects is the generic object collection function
func GetObjects(v1 *viper.Viper, mc *minio.Client) ([]string, error) {
	oa := []string{}

	objs, err := config.GetObjectsConfig(v1)
	var pa = objs.Prefix

	log.Println(pa)

	for p := range pa {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		bucketName, _ := config.GetBucketName(v1)
		objectCh := mc.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Prefix: pa[p], Recursive: true})

		for object := range objectCh {
			if object.Err != nil {
				log.Println(object.Err)
				return oa, object.Err
			}
			// fmt.Println(object)
			oa = append(oa, object.Key)
		}
		log.Printf("%s:%s object count: %d\n", bucketName, pa[p], len(oa))
	}

	return oa, err
}

// curl -u admin:Complexpass#123 -XPUT -d '{"name":"Prabhat Sharma"}' http://localhost:4080/api/myshinynewindex/document
func docfunc(v1 *viper.Viper, mc *minio.Client, bucketName string, item string, endpoint string) ([]byte, error) {

	// s, err := loader(v1, mc, objs["bucket"], oa[item], spql["endpoint"])
	//jo := strings.Replace(item, "milled", "summoned", 1)
	//jo2 := strings.Replace(jo, ".rdf", ".jsonld", 1)
	jo2 := item

	//b, _, err := objects.GetS3Bytes(mc, objs["bucket"], jo2)
	b, _, err := objects.GetS3Bytes(mc, bucketName, jo2)
	if err != nil {
		//log.Printf("%s : %s \n", objs["bucket"], jo2)
		log.Printf("%s : %s \n", bucketName, jo2)
		log.Println(err)
	}

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
		log.Println(err)
	}
	defer resp.Body.Close()

	// fmt.Println("Tika Response Status:", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(body))

	// TESTING
	log.Printf("%s : %d  : %s\n", jo2, len(b), endpoint)

	return nil, nil
}
