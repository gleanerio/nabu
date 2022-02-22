package jena

import (
	"fmt"
	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/minio/minio-go/v7"
)

// BulkAssembly collects the objects from a bucket to load
func BulkAssembly(v1 *viper.Viper, mc *minio.Client) error {
	bucketName, _ := config.GetBucketName(v1)

	// TODO remove TESTING settings for bulk loading
	prefix := "milled/aquadocs"
	name := "testobject.rdf"

	err := objects.PipeCopyNG(name, bucketName, prefix, mc)
	if err != nil {
		return err
	}

	// will need a function call at some point to work with the new object
	r, err := docfunc(v1, mc, bucketName, fmt.Sprintf("%s/%s", prefix, name), "endpoint")
	if err != nil {
		log.Println(err)
	}

	log.Printf("Return from docfunc: %s", string(r))

	// TODO  remove the temporary object?

	return err
}
