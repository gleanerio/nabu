package pkg

import (
	"github.com/gleanerio/nabu/internal/services/bulk"
	log "github.com/sirupsen/logrus"

	"github.com/gleanerio/nabu/internal/objects"
	"github.com/spf13/viper"

	"github.com/minio/minio-go/v7"
)

func NabuGraphDB(v1 *viper.Viper) error {
	mc, err := objects.MinioConnection(v1)
	if err != nil {
		log.Fatal("cannot connect to minio: %s", err)
	}
	return Bulk(v1, mc)
}

func GraphDB(v1 *viper.Viper, mc *minio.Client) error {
	//err := bulk.ObjectAssembly(v1, mc)
	err := bulk.BulkAssembly(v1, mc)

	if err != nil {
		log.Error(err)
	}
	return err
}
