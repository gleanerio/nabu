package objects

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/minio/minio-go/v7"
)

// PipeCopyNG writes a new object based on an prefix, this function assumes the objects are valid when concatenated
// name:  name of the NEW object
// bucket:  source bucket  (and target bucket)
// prefix:  source prefix
// mc:  minio client pointer
func PipeCopyNG(name, bucket, prefix string, mc *minio.Client) error {
	log.Printf("Pipecopy with name: %s   bucket: %s  prefix: %s", name, bucket, prefix)

	pr, pw := io.Pipe()     // TeeReader of use?
	lwg := sync.WaitGroup{} // work group for the pipe writes...
	lwg.Add(2)

	// params for list objects calls
	doneCh := make(chan struct{}) // , N) Create a done channel to control 'ListObjectsV2' go routine.
	defer close(doneCh)           // Indicate to our routine to exit cleanly upon return.
	isRecursive := true

	go func() {
		defer lwg.Done()
		defer func(pw *io.PipeWriter) {
			err := pw.Close()
			if err != nil {
			}
		}(pw)

		objectCh := mc.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: isRecursive})

		// for object := range mc.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: isRecursive}, doneCh) {
		for object := range objectCh {
			fo, err := mc.GetObject(context.Background(), bucket, object.Key, minio.GetObjectOptions{})
			if err != nil {
				fmt.Println(err)
			}

			var b bytes.Buffer
			bw := bufio.NewWriter(&b)

			_, err = io.Copy(bw, fo)
			if err != nil {
				log.Println(err)
			}

			_, err = pw.Write(b.Bytes())
			if err != nil {
				return
			}
		}

	}()

	// go function to write to minio from pipe
	go func() {
		defer lwg.Done()
		_, err := mc.PutObject(context.Background(), bucket, fmt.Sprintf("%s/%s", prefix, name), pr, -1, minio.PutObjectOptions{})
		if err != nil {
			log.Println(err)
			return
		}
	}()

	lwg.Wait() // wait for the pipe read writes to finish
	err := pw.Close()
	if err != nil {
		return err
	}
	err = pr.Close()
	if err != nil {
		return err
	}

	return nil
}
