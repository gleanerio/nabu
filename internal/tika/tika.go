package tika

import (
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/minio/minio-go"
	"github.com/spf13/viper"
)

// Build launches a go func to build.   Needs to return a
func Build(v1 *viper.Viper, mc *minio.Client) error {
	// go builder(bucket, prefix, domain, mc)
	go builder(v1, mc)

	return nil
}

func builder(v1 *viper.Viper, mc *minio.Client) {

	objs := v1.GetStringMapString("objects")
	bucket := objs["bucket"]
	prefix := objs["prefix"]
	domain := objs["domain"]

	log.Printf("Bucket: %s  Prefix: %s   Domain: %s\n", bucket, prefix, domain)

	// Create a done channel.
	doneCh := make(chan struct{})
	defer close(doneCh)
	recursive := true

	// Pipecopy elements
	pr, pw := io.Pipe()     // TeeReader of use?
	lwg := sync.WaitGroup{} // work group for the pipe writes...
	lwg.Add(2)

	go func() {
		defer lwg.Done()
		defer pw.Close()

		// WARNING hard coded "prefix" here
		for message := range mc.ListObjectsV2(bucket, fmt.Sprintf("%s/csdco/do", prefix), recursive, doneCh) {

			if !strings.HasSuffix(message.Key, ".jsonld") {
				log.Println(message.Key)

				s, err := processObject(v1, mc, bucket, prefix, message)
				if err != nil {
					log.Println(err)
				}

				pw.Write([]byte(s))
			}
		}

	}()

	go func() {
		defer lwg.Done()
		var op string
		if prefix == "" {
			op = fmt.Sprint("website/fulltext.nq") // should this be website?
		} else {
			op = fmt.Sprintf("%s/website/fulltext.nq", prefix)
		}

		log.Println(op)

		_, err := mc.PutObject(bucket, op, pr, -1, minio.PutObjectOptions{}) // TODO  this is potentially dangerous..  it will over write this object at least
		if err != nil {
			log.Println(err)
		}
	}()

	lwg.Wait() // wait for the pipe read writes to finish
	pw.Close()
	pr.Close()

	log.Println("Builder call done")

}
