package pkg

import (
	"fmt"
	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

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
	spql, _ := config.GetSparqlConfig(v1)
	if bucket == "" {
		bucket, _ = config.GetBucketName(v1)
	}
	s, err := objects.PipeLoad(v1, mc, bucket, object, spql.Endpoint)
	if err != nil {
		log.Error(err)
	}
	log.Trace(string(s))
	return err
}
