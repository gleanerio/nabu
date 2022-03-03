package meili

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/pkg/config"
	"github.com/meilisearch/meilisearch-go"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/sjson"

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

	fmt.Println("MeiliSearch index started")

	msc := meilisearch.NewClient(meilisearch.ClientConfig{
		Host: "http://127.0.0.1:7700",
	})

	bar := progressbar.Default(int64(len(oa)))

	// Single threaded loop
	for item := range oa {
		_, err := docfunc(v1, mc, msc, bucketName, oa[item], "endpoint")
		if err != nil {
			log.Println(err)
		}
		bar.Add(1)
	}

	return err
}

// curl -u admin:Complexpass#123 -XPUT -d '{"name":"Prabhat Sharma"}' http://localhost:4080/api/myshinynewindex/document
func docfunc(v1 *viper.Viper, mc *minio.Client, msc *meilisearch.Client, bucketName string, item string, endpoint string) ([]byte, error) {
	// get item
	b, _, err := objects.GetS3Bytes(mc, bucketName, item)
	if err != nil {
		return nil, err
	}
	s := string(b)

	// Build ID entry for JSON
	fp := filepath.Base(item)
	nns := strings.TrimSuffix(fp, path.Ext(fp))

	value, _ := sjson.Set(s, "id", nns)
	log.Println("---------------------------------------------------------------")
	log.Println(value)
	log.Println("---------------------------------------------------------------")

	var doc map[string]interface{}
	json.Unmarshal([]byte(value), &doc)

	//r, err := msc.Index("testi").AddDocuments(movies)
	r, err := msc.Index("movies").AddDocuments(doc)
	if err != nil {
		log.Println(err)
	}

	log.Println(r)

	return nil, nil
}
