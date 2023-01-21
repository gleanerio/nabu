package pkg

import (
	"github.com/gleanerio/nabu/internal/services/bulk"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/minio/minio-go/v7"
)

func Bulk(v1 *viper.Viper, mc *minio.Client) error {
	//err := bulk.ObjectAssembly(v1, mc)
	err := bulk.BulkAssembly(v1, mc)

	if err != nil {
		log.Error(err)
	}
	return err
}
