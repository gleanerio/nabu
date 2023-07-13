package pkg

import (
	"fmt"
	"github.com/gleanerio/nabu/internal/common"
	"github.com/gleanerio/nabu/internal/objects"

	"github.com/gleanerio/nabu/internal/services/meili"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Meili(v1 *viper.Viper, mc *minio.Client) error {
	fmt.Println("Index object into MeiliSearch")

	err := meili.ObjectAssembly(v1, mc)

	if err != nil {
		fmt.Println(err)
	}
	return err
}

// used by glcon in gleaner. Need to develop a more common config for the services (aka s3, graph, etc)
// cannot pass a nabu config to the gleaner code to create a minio client, and have it work
func NabuMeili(v1 *viper.Viper) error {
	common.InitLogging()
	mc, err := objects.MinioConnection(v1)
	if err != nil {
		log.Fatal("cannot connect to minio: %s", err)
	}
	return Meili(v1, mc)
}
