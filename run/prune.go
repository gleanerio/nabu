package run

import (
	"fmt"
	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/internal/prune"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
	"log"
)

func NabuPrune(v1 *viper.Viper) error {
	mc, err := objects.MinioConnection(v1)
	if err != nil {
		log.Fatal("cannot connect to minio: %s", err)
	}
	return Prune(v1, mc)
}
func Prune(v1 *viper.Viper, mc *minio.Client) error {
	fmt.Println("Prune graphs in triplestore not in object store")
	err := prune.Snip(v1, mc)
	if err != nil {
		log.Println(err)
	}
	return err
}
