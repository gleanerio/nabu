package pkg

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/internal/services/meili"
	"github.com/spf13/viper"

	"github.com/minio/minio-go/v7"
)

func NabuMeili(v1 *viper.Viper) error {
	mc, err := objects.MinioConnection(v1)
	if err != nil {
		log.Fatal("cannot connect to minio: %s", err)
	}
	return Meili(v1, mc)
}

func Meili(v1 *viper.Viper, mc *minio.Client) error {
	fmt.Println("Index object into MeiliSearch")

	err := meili.ObjectAssembly(v1, mc)

	if err != nil {
		fmt.Println(err)
	}
	return err
}
