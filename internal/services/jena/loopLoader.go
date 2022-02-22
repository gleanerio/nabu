package jena

import (
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

	//Single threaded loop
	for item := range oa {
		_, err := docfunc(v1, mc, bucketName, oa[item], "endpoint")
		if err != nil {
			log.Println(err)
		}
		bar.Add(1)
	}

	return err
}
