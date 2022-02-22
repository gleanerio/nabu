package pkg

import (
	"fmt"
	"log"

	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/internal/services/zinc"
	"github.com/spf13/viper"

	"github.com/minio/minio-go/v7"
)

func NabuZinc(v1 *viper.Viper) error {
	mc, err := objects.MinioConnection(v1)
	if err != nil {
		log.Fatal("cannot connect to minio: %s", err)
	}
	return Zinc(v1, mc)
}

func Zinc(v1 *viper.Viper, mc *minio.Client) error {
	fmt.Println("Tika extract text from objects")

	err := zinc.ObjectAssembly(v1, mc)

	if err != nil {
		fmt.Println(err)
	}
	return err
}
