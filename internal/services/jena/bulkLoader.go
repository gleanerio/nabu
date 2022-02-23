package jena

import (
	"context"
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
	objCfg, _ := config.GetObjectsConfig(v1)
	pa := objCfg.Prefix

	var err error

	name := "bulkobject.rdf"

	for p := range pa {
		err = objects.PipeCopyNG(name, bucketName, pa[p], mc)
		if err != nil {
			return err
		}
	}

	for p := range pa {
		// will need a function call at some point to work with the new object
		r, err := docfunc(v1, mc, bucketName, fmt.Sprintf("%s/%s", pa[p], name), "endpoint")
		if err != nil {
			log.Println(err)
		}

		log.Printf("Return from docfunc: %s", string(r))
	}

	// TODO  remove the temporary object?
	for p := range pa {
		opts := minio.RemoveObjectOptions{}
		err = mc.RemoveObject(context.Background(), bucketName, fmt.Sprintf("%s/%s", pa[p], name), opts)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return err
}
