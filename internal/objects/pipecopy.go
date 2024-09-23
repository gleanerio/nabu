package objects

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
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

		lastProcessed := false
		idList := make([]string, 0)
		for object := range objectCh {
			fo, err := mc.GetObject(context.Background(), bucket, object.Key, minio.GetObjectOptions{})
			if err != nil {
				fmt.Println(err)
				continue
			}
			var b bytes.Buffer
			bw := bufio.NewWriter(&b)
			_, err = io.Copy(bw, fo)
			if err != nil {
				log.Println(err)
				continue
			}
			s := string(b.Bytes())
			nq := ""
			if strings.HasSuffix(object.Key, ".nq") {
				nq = s
			} else {
				nq, err = graph.JSONLDToNQ(v1, s)
				if err != nil {
					log.Println(err)
					return
				}
			}
			var snq string
			if sf {
				snq = nq
			} else {
				snq, err = graph.Skolemization(nq, object.Key)
				if err != nil {
					return
				}
			}
			ctx, err := graph.MakeURN(v1, object.Key)
			if err != nil {
				return
			}
			csnq, err := graph.NtToNq(snq, ctx)
			if err != nil {
				return
			}
			_, err = pw.Write([]byte(csnq))
			if err != nil {
				return
			}
			idList = append(idList, ctx)
			lastProcessed = true
		}

		// Once we are done with the loop, put in the triples to associate all the graphURIs with the org.
		if lastProcessed {

			data := `_:b0 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://schema.org/DataCatalog> .
_:b0 <https://schema.org/dateCreated> "2024-09-20" .
_:b0 <https://schema.org/description> "This is an example data catalog containing various datasets from this organization" .
_:b0 <https://schema.org/provider> _:b1 .
_:b0 <https://schema.org/publisher> _:b2 .
_:b1 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://schema.org/Organization> .
_:b1 <https://schema.org/name> "Provider XYZ" .
_:b2 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://schema.org/Organization> .
_:b2 <https://schema.org/name> "DeCoder" .
`
			for _, item := range idList {
				data += `_:b0 <https://schema.org/dataset> <` + item + `> .` + "\n"
			}

			sdata, err := graph.Skolemization(data, "release graph prov for ORG")
			if err != nil {
				log.Println(err)
			}

			// Perform the final write to the pipe here
			// ilstr := strings.Join(idList, ",")
			_, err = pw.Write([]byte(sdata))
			if err != nil {
				log.Println(err)
			}
		}
	}()

	// go function to write to minio from pipe
	go func() {
		defer lwg.Done()
		_, err := mc.PutObject(context.Background(), bucket, fmt.Sprintf("%s/%s", destprefix, name), pr, -1, minio.PutObjectOptions{})
		//_, err := mc.PutObject(context.Background(), bucket, fmt.Sprintf("%s/%s", prefix, name), pr, -1, minio.PutObjectOptions{})
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
