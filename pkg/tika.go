package pkg

import (
	"fmt"
	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/internal/services/tika"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/minio/minio-go/v7"
)

func NabuTika(v1 *viper.Viper) error {
	mc, err := objects.MinioConnection(v1)
	if err != nil {
		log.Fatal("cannot connect to minio: %s", err)
	}
	return Tika(v1, mc)
}
func Tika(v1 *viper.Viper, mc *minio.Client) error {
	fmt.Println("Tika extract text from objects")
	err := tika.SingleBuild(v1, mc)

	if err != nil {
		log.Println(err)
	}
	return err
}
