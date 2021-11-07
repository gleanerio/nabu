package pkg

import (
	"fmt"
	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/internal/semsearch"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
	"log"
)

func NabuTxtai(v1 *viper.Viper) error {
	mc, err := objects.MinioConnection(v1)
	if err != nil {
		log.Fatal("cannot connect to minio: %s", err)
	}
	return Txtai(v1, mc)
}

func Txtai(v1 *viper.Viper, mc *minio.Client) error {
	fmt.Println("Index descriptions to txtai")
	err := semsearch.ObjectAssembly(v1, mc)
	if err != nil {
		log.Println(err)
	}
	return err
}
