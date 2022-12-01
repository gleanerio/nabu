package objects

import (
	"bufio"
	"bytes"
	"context"
	"github.com/gleanerio/nabu/internal/graph"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
)

func MillerNG(name, bucket, prefix string, mc *minio.Client) error {
	log.Printf("MillerNG: %s   bucket: %s  prefix: %s", name, bucket, prefix)

	opts := minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    prefix,
	}

	// for object := range mc.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: isRecursive}, doneCh) {
	for object := range mc.ListObjects(context.Background(), bucket, opts) {
		if object.Err != nil {
			log.Println(object.Err)
			log.Println(bucket)
			log.Println(opts.Prefix)
			return object.Err
		}

		fo, err := mc.GetObject(context.Background(), bucket, object.Key, minio.GetObjectOptions{})
		if err != nil {
			log.Println(err)
		}

		var b bytes.Buffer
		bw := bufio.NewWriter(&b)

		_, err = io.Copy(bw, fo)
		if err != nil {
			log.Println(err)
		}

		s := string(b.Bytes())

		log.Println("Calling JSONLDtoNQ")
		nq, err := graph.JSONLDToNQ(s)
		if err != nil {
			log.Println(err)
			return err
		}

		// TODO add the context into this fle  (then load to Jena withouth explicate graph)
		// g = fmt.Sprintf("urn:%s:%s", bucketName, strings.TrimSuffix(s2c, ".rdf"))
		// func NQNewGraph(inquads, newctx string) (string, string, error) {

		log.Println("Calling Skolemization")
		nqs, err := graph.Skolemization(nq)
		if err != nil {
			return err
		}

		contentType := "application/ld+json"
		prefixmod := strings.ReplaceAll(object.Key, "summoned", "scratch")
		keymod := strings.ReplaceAll(prefixmod, "jsonld", "rdf")

		//newkey := fmt.Sprintf("%s/%s", prefixmod, keymod)
		log.Println(keymod)

		i, err := PutS3Bytes(mc, bucket, keymod, contentType, []byte(nqs))
		if err != nil {
			return err
		}

		log.Printf("Put objectg len %d", i)
		// write the new nqs to a miller area..  or a temp dir

	}

	return nil
}
