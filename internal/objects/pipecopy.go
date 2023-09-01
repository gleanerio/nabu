package objects

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/spf13/viper"

	"github.com/gleanerio/nabu/internal/graph"
	log "github.com/sirupsen/logrus"

	"github.com/minio/minio-go/v7"
)

// PipeCopy writes a new object based on an prefix, this function assumes the objects are valid when concatenated
// v1:  viper config object
// mc:  minio client pointer
// name:  name of the NEW object
// bucket:  source bucket  (and target bucket)
// prefix:  source prefix
// destprefix:   destination prefix
// sf: boolean to declare if single file or not.   If so, skip skolimization since JSON-LD library output is enough
func PipeCopy(v1 *viper.Viper, mc *minio.Client, name, bucket, prefix, destprefix string) error {
	log.Printf("PipeCopy with name: %s   bucket: %s  prefix: %s", name, bucket, prefix)

	pr, pw := io.Pipe()     // TeeReader of use?
	lwg := sync.WaitGroup{} // work group for the pipe writes...
	lwg.Add(2)

	// params for list objects calls
	doneCh := make(chan struct{}) // , N) Create a done channel to control 'ListObjectsV2' go routine.
	defer close(doneCh)           // Indicate to our routine to exit cleanly upon return.
	isRecursive := true

	//log.Printf("Bulkfile name: %s_graph.nq", name)

	go func() {
		defer lwg.Done()
		defer func(pw *io.PipeWriter) {
			err := pw.Close()
			if err != nil {
			}
		}(pw)

		// Set and use a "single file flag" to bypass skolimaization since if it is a single file
		// the JSON-LD to RDF will correctly map blank nodes.
		// NOTE:  with a background context we can't get the len(channel) so we have to iterate it.
		// This is fast, but it means we have to do the ListObjects twice
		clen := 0
		sf := false
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		lenCh := mc.ListObjects(ctx, bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: isRecursive})
		for _ = range lenCh {
			clen = clen + 1
		}
		if clen == 1 {
			sf = true
		}
		log.Printf("\nChannel/object length: %d\n", clen)
		log.Printf("Single file mode set: %t", sf)

		objectCh := mc.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: isRecursive})

		// for object := range mc.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: isRecursive}, doneCh) {
		for object := range objectCh {
			fo, err := mc.GetObject(context.Background(), bucket, object.Key, minio.GetObjectOptions{})
			if err != nil {
				log.Errorf(" failed to read object %s %s ", object.Key, err)
				//fmt.Println(err)
			}
			log.Tracef(" processing  object %s  ", object.Key)
			var b bytes.Buffer
			bw := bufio.NewWriter(&b)

			_, err = io.Copy(bw, fo)
			if err != nil {
				log.Errorf(" failed to read object %s %s ", object.Key, err)
				log.Println(err)
			}

			s := string(b.Bytes())

			//log.Println("Calling JSONLDtoNQ")
			nq, err := graph.JSONLDToNQ(v1, s)
			if err != nil {
				log.Errorf(" failed to convert to NQ %s %s ", object.Key, err)
				continue
			}

			var snq string

			if sf {
				snq = nq //  just pass through the RDF without trying to Skolemize since we ar a single fil
			} else {
				snq, err = graph.Skolemization(nq, object.Key)
				if err != nil {
					log.Errorf(" failed Skolemization %s %s ", object.Key, err)
					continue
				}
			}

			// 1) get graph URI
			ctx, err := graph.MakeURN(object.Key, bucket)
			if err != nil {
				log.Errorf(" failed MakeURN %s %s ", object.Key, err)
				continue
			}
			// 2) convert NT to NQ
			csnq, err := graph.NtToNq(snq, ctx)
			if err != nil {
				log.Errorf(" failed NtToNq %s %s ", object.Key, err)
				continue
			}

			_, err = pw.Write([]byte(csnq))
			if err != nil {
				log.Errorf(" failed pipe write %s %s ", object.Key, err)
				continue
			}
		}
	}()

	// go function to write to minio from pipe
	go func() {
		defer lwg.Done()
		_, err := mc.PutObject(context.Background(), bucket, fmt.Sprintf("%s/%s", destprefix, name), pr, -1, minio.PutObjectOptions{})
		//_, err := mc.PutObject(context.Background(), bucket, fmt.Sprintf("%s/%s", prefix, name), pr, -1, minio.PutObjectOptions{})
		if err != nil {
			log.Errorf(" failed PutObject bucket: %s  %s/%s ", bucket, destprefix, name)
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
