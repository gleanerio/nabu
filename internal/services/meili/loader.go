package meili

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/pkg/config"
	"github.com/meilisearch/meilisearch-go"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

// ObjectAssembly collects the objects from a bucket to load
func ObjectAssembly(v1 *viper.Viper, mc *minio.Client) error {
	fmt.Println("MeiliSearch index started")

	bucketName, _ := config.GetBucketName(v1)
	objCfg, _ := config.GetObjectsConfig(v1)
	pa := objCfg.Prefix

	msc := meilisearch.NewClient(meilisearch.ClientConfig{
		//Host: "http://127.0.0.1:7700",
		Host: "https://index.geoconnex.us",
	})

	var err error

	name := "bulkobject.json"

	for p := range pa {
		log.Printf("ToJSONArray for %s", pa[p])
		objects.ToJSONArray(name, bucketName, pa[p], mc)
		if err != nil {
			return err
		}
	}

	for p := range pa {
		log.Printf("dofunc for %s  %s", p, fmt.Sprintf("%s/%s", pa[p], name))
		// will need a function call at some point to work with the new object
		r, err := docfunc(v1, mc, msc, bucketName, fmt.Sprintf("%s/%s", pa[p], name), "endpoint")
		if err != nil {
			log.Println(err)
		}
		log.Printf("Return from docfunc: %s", string(r))
	}

	// TODO  remove the temporary object?
	for p := range pa {
		log.Printf("remove %s for %s", name, pa[p])
		opts := minio.RemoveObjectOptions{}
		err = mc.RemoveObject(context.Background(), bucketName, fmt.Sprintf("%s/%s", pa[p], name), opts)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return err
}

// ObjectAssembly collects the objects from a bucket to load
func ObjectAssemblyORIG(v1 *viper.Viper, mc *minio.Client) error {
	bucketName, _ := config.GetBucketName(v1)

	oa, err := objects.GetObjects(v1, mc)
	if err != nil {
		return err
	}

	fmt.Println("MeiliSearch index started")

	msc := meilisearch.NewClient(meilisearch.ClientConfig{
		//Host: "http://127.0.0.1:7700",
		Host: "https://index.geoconnex.us",
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

	// Build ID entry for JSON
	//fp := filepath.Base(item)
	//nns := strings.TrimSuffix(fp, path.Ext(fp))

	//s := string(b)
	//value, _ := sjson.Set(s, "id", nns)
	//log.Println("---------------------------------------------------------------")
	//log.Println(value)
	//log.Println("---------------------------------------------------------------")

	var doc interface{} // why was this a map[string]interface{} before?  bulk vs single..  doesn't seem so..
	//json.Unmarshal([]byte(value), &doc)
	json.Unmarshal(b, &doc)

	//r, err := msc.Index("testi").AddDocuments(movies)
	r, err := msc.Index("iow").AddDocuments(doc, "id")
	if err != nil {
		log.Println(err)
	}

	log.Println(r)

	return nil, nil
}
