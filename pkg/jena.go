package pkg

import (
	"fmt"
	"github.com/gleanerio/nabu/internal/services/jena"
	log "github.com/sirupsen/logrus"

	"github.com/gleanerio/nabu/internal/objects"
	"github.com/spf13/viper"

	"github.com/minio/minio-go/v7"
)

func NabuJena(v1 *viper.Viper) error {
	mc, err := objects.MinioConnection(v1)
	if err != nil {
		log.Fatal("cannot connect to minio: %s", err)
	}
	return Jena(v1, mc)
}

func Jena(v1 *viper.Viper, mc *minio.Client) error {
	//err := jena.ObjectAssembly(v1, mc)
	err := jena.BulkAssembly(v1, mc)

	if err != nil {
		fmt.Println(err)
	}
	return err
}
