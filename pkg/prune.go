package pkg

import (
	"github.com/gleanerio/nabu/internal/prune"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Prune(v1 *viper.Viper, mc *minio.Client) error {
	log.Info("Prune graphs in triplestore not in objectVal store")
	err := prune.Snip(v1, mc)
	if err != nil {
		log.Error(err)
	}
	return err
}
