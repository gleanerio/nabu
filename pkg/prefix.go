package pkg

import (
	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/internal/sparqlapi"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Prefix(v1 *viper.Viper, mc *minio.Client) error {

	log.Println("Load graphs from prefix to triplestore")
	err := sparqlapi.ObjectAssembly(v1, mc)

	if err != nil {
		log.Println(err)
	}
	return err

}

func NabuPrefix(v1 *viper.Viper) error {
	mc, err := objects.MinioConnection(v1)
	if err != nil {
		log.Fatal("cannot connect to minio: %s", err)
	}
	return Prefix(v1, mc)
}
