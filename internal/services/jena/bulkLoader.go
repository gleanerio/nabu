package jena

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

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

	//name := "bulkobject.rdf"

	for p := range pa {
		name := fmt.Sprintf("%s_bulk.rdf", baseName(path.Base(pa[p])))

		//err = objects.PipeCopyNG(name, bucketName, pa[p], mc) // have this function return the object name and path, easy to load and remove then
		//err = objects.PipeCopyJLD2NQ(name, bucketName, pa[p], mc) // have this function return the object name and path, easy to load and remove then
		err = objects.MillerNG(name, bucketName, pa[p], mc) // have this function return the object name and path, easy to load and remove then

		if err != nil {
			return err
		}
	}

	os.Exit(0)

	for p := range pa {
		// will need a function call at some point to work with the new object
		name := fmt.Sprintf("%s_bulk.rdf", pa[p])
		r, err := docfunc(v1, mc, bucketName, fmt.Sprintf("%s/%s", pa[p], name), "endpoint")
		if err != nil {
			log.Println(err)
		}
		log.Printf("docfunc: %s", string(r))
	}

	// TODO  remove the temporary object?
	for p := range pa {
		name := fmt.Sprintf("%s_bulk.rdf", pa[p])
		opts := minio.RemoveObjectOptions{}
		err = mc.RemoveObject(context.Background(), bucketName, fmt.Sprintf("%s/%s", pa[p], name), opts)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return err
}

func baseName(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n == -1 {
		return s
	}
	return s[:n]
}
