package textindex

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"

	"github.com/blevesearch/bleve"
	"github.com/gleanerio/gleaner/internal/common"
	"github.com/gleanerio/gleaner/internal/millers/millerutils"
	minio "github.com/minio/minio-go"
)

// GetObjects test a concurrent version of calling mock
func GetObjects(mc *minio.Client, bucketname string) {
	indexname := fmt.Sprintf("%s", bucketname)
	fp := millerutils.NewinitBleve(indexname) //  initBleve(indexname)
	entries := common.GetMillObjects(mc, bucketname)
	multiCall(entries, fp)

}

func multiCall(e []common.Entry, indexname string) {
	// TODO..   open the bleve index here once and pass by reference to text
	index, berr := bleve.Open(indexname)
	if berr != nil {
		// should panic here?..  no index..  no reason to keep living  :(
		log.Printf("Bleve error making index %v \n", berr)
	}

	// Set up the the semaphore and conccurancey
	semaphoreChan := make(chan struct{}, 1) //For direct write like this must be SINGLE THREADED!!!!!!
	defer close(semaphoreChan)
	wg := sync.WaitGroup{} // a wait group enables the main process a wait for goroutines to finish

	for k := range e {
		wg.Add(1)
		log.Printf("About to run #%d in a goroutine\n", k)
		go func(k int) {
			semaphoreChan <- struct{}{}

			status := textIndexer(e[k].Urlval, e[k].Jld, index)

			wg.Done() // tell the wait group that we be done
			log.Printf("#%d done with %s", k, status)
			<-semaphoreChan
		}(k)
	}
	wg.Wait()

	index.Close()
}

// index some jsonld with an ID
func textIndexer(ID string, jsonld string, index bleve.Index) string {
	berr := index.Index(ID, jsonld)
	log.Printf("Bleve Indexed item with ID %s\n", ID)
	if berr != nil {
		log.Printf("Bleve error indexing %v \n", berr)
	}

	return "done"
}
