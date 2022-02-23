package jena

import (
	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/pkg/config"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"

	"github.com/minio/minio-go/v7"
)

// ThreadedLoader collects the objects from a bucket to load
func ThreadedLoader(v1 *viper.Viper, mc *minio.Client) error {
	bucketName, _ := config.GetBucketName(v1)

	oa, err := objects.GetObjects(v1, mc)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(oa)))
	// TODO Go func version
	semaphoreChan := make(chan struct{}, 15) // a blocking channel to keep concurrency under control
	defer close(semaphoreChan)
	wg := sync.WaitGroup{}

	log.Println("Threaded run testing")

	for item := range oa {
		wg.Add(1)
		go func(item int) {
			semaphoreChan <- struct{}{}

			_, err := docfunc(v1, mc, bucketName, oa[item], "endpoint")
			if err != nil {
				log.Println(err)
			}

			wg.Done()
			bar.Add(1)
			<-semaphoreChan // clear a spot in the semaphore channel for the next indexing event
		}(item)
	}
	wg.Wait()

	return err
}
