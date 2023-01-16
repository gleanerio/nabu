package pkg

import (
	"github.com/gleanerio/nabu/internal/objects"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Prefix(v1 *viper.Viper, mc *minio.Client) error {
	log.Info("Nabu started with mode: prefix")
	err := objects.ObjectAssembly(v1, mc)

	if err != nil {
		log.Error(err)
	}
	return err
}
