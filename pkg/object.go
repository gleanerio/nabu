package pkg

import (
	"fmt"
	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/internal/sparqlapi"
	"github.com/gleanerio/nabu/pkg/config"
	"github.com/spf13/viper"
	"log"

	"github.com/minio/minio-go/v7"
)

func NabuObject(v1 *viper.Viper, bucket string, object string) error {
	mc, err := objects.MinioConnection(v1)
	if err != nil {
		log.Fatal("cannot connect to minio: %s", err)
	}
	return Object(v1, mc, bucket, object)
}
func Object(v1 *viper.Viper, mc *minio.Client, bucket string, object string) error {
	fmt.Println("Load graph object to triplestore")
	//spql := v1.GetStringMapString("sparql")
	spql, _ := config.GetSparqlConfig(v1)
	if bucket == "" {
		bucket, _ = config.GetBucketName(v1)
	}
	//s, err := sparqlapi.PipeLoad(v1, mc, bucket, object, spql["endpoint"])
	s, err := sparqlapi.PipeLoad(v1, mc, bucket, object, spql.Endpoint)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(s))
	return err
}
