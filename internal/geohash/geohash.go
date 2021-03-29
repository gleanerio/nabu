package geohash

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/UFOKN/nabu/internal/objects"
	"github.com/schollz/progressbar/v3"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

// ObjectAssembly collects the objects from a bucket to load
func ObjectAssembly(v1 *viper.Viper, mc *minio.Client) error {
	objs := v1.GetStringMapString("objects")
	spql := v1.GetStringMapString("sparql")

	// My go func controller vars
	semaphoreChan := make(chan struct{}, 1) // a blocking channel to keep concurrency under control (1 == single thread)
	defer close(semaphoreChan)
	// wg := sync.WaitGroup{} // a wait group enables the main process a wait for goroutines to finish

	// params for list objects calls
	doneCh := make(chan struct{}) // , N) Create a done channel to control 'ListObjectsV2' go routine.
	defer close(doneCh)           // Indicate to our routine to exit cleanly upon return.

	oa := []string{}

	// NEW
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	objectCh := mc.ListObjects(ctx, objs["bucket"],
		minio.ListObjectsOptions{Prefix: objs["prefix"], Recursive: true})

	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return object.Err
		}
		// fmt.Println(object)
		oa = append(oa, object.Key)
	}

	log.Printf("%s:%s object count: %d\n", objs["bucket"], objs["prefix"], len(oa))
	bar := progressbar.Default(int64(len(oa)))
	for item := range oa {
		_, err := PipeLoad(v1, mc, objs["bucket"], oa[item], spql["endpoint"])
		if err != nil {
			log.Println(err)
		}
		bar.Add(1)
		// log.Println(string(s)) // get "s" on pipeload and send to a log file
	}

	return nil
}

// LoadHash reads from an object and loads directly into a triplestore
func LoadHash(v1 *viper.Viper, mc *minio.Client, bucket, object, spql string) ([]byte, error) {
	// build our quad/graph from the object path

	s2c := strings.Replace(object, "/", ":", -1)
	g := fmt.Sprintf("urn:%s:%s", bucket, strings.TrimSuffix(s2c, ".rdf"))

	// Turn checking off while testing other parts of Nabu
	c, err := gexists(spql, g)
	if err != nil {
		log.Println(err)
	}
	if c {
		return nil, nil // our graph is loaded already..
	}

	b, _, err := objects.GetS3Bytes(mc, bucket, object)
	if err != nil {
		fmt.Printf("gets3Bytes %v\n", err)
	}

	fmt.Println(len(string(b)))

	// need to parse B and get out the lat long
	// call a client to load the lat long point into tile38

	return []byte("resp.Status"), err
}
