package semsearch

import (
	"context"
	"encoding/json"
	"log"
	"path"
	"strings"
	"sync"

	"github.com/UFOKN/nabu/internal/jsonld"
	"github.com/UFOKN/nabu/internal/objects"

	"github.com/minio/minio-go/v7"
	"github.com/neuml/txtai.go"
	"github.com/spf13/viper"
)

type DataFrame struct {
	ID          string `json:"@id"`
	Description string `json:"description"`
}

// ObjectAssembly collects the objects from a bucket to load
func ObjectAssembly(v1 *viper.Viper, mc *minio.Client) error {
	objs := v1.GetStringMapString("objects")

	// My go func controller vars
	semaphoreChan := make(chan struct{}, 1) // a blocking channel to keep concurrency under control (1 == single thread)
	defer close(semaphoreChan)
	wg := sync.WaitGroup{} // a wait group enables the main process a wait for goroutines to finish

	// params for list objects calls
	doneCh := make(chan struct{}) // , N) Create a done channel to control 'ListObjectsV2' go routine.
	defer close(doneCh)           // Indicate to our routine to exit cleanly upon return.
	// isRecursive := true

	oa := []string{}
	// for object := range mc.ListObjects(objs["bucket"], objs["prefix"], isRecursive, doneCh) {
	for object := range mc.ListObjects(context.Background(), objs["bucket"],
		minio.ListObjectsOptions{Prefix: objs["prefix"], Recursive: true}) {
		wg.Add(1)
		go func(object minio.ObjectInfo) {
			// log.Println(object.Key)
			oa = append(oa, object.Key)
			wg.Done() // tell the wait group that we be done
			// log.Printf("Doc: %s error: %v ", name, err) // why print the status??
			<-semaphoreChan
		}(object)
	}
	wg.Wait()

	//parseloader(v1, mc, oa)
	frameloader(v1, mc, oa)
	return nil
}

func frameloader(v1 *viper.Viper, mc *minio.Client, oa []string) ([]byte, error) {
	objs := v1.GetStringMapString("objects")
	t := v1.GetStringMapString("txtaipkg")

	embeddings := txtai.Embeddings(t["endpoint"])

	for item := range oa {
		// s, err := loader(v1, mc, objs["bucket"], oa[item], spql["endpoint"])
		jo := strings.Replace(oa[item], "milled", "summoned", 1)
		jo2 := strings.Replace(jo, ".rdf", ".jsonld", 1)
		b, _, err := objects.GetS3Bytes(mc, objs["bucket"], jo2)
		if err != nil {
			log.Printf("%s : %s \n", objs["bucket"], jo2)
			log.Println(err)
			continue
		}

		desc, err := httpsDesc(b)
		if err != nil {
			return nil, err
		}

		// if desc nil try http version of the frame
		if desc == "" {
			desc, err = httpDesc(b)
			if err != nil {
				return nil, err
			}
		}

		// TESTING
		log.Printf("%s : %s \n", jo2, desc)

		if desc == "" {
			continue
		}

		// get the base part only of the object
		// osplt := strings.Split(oa[item], "/")
		// o := osplt[len(osplt)-1]
		//fmt.Printf("%s \n", strings.TrimSuffix(o, path.Ext(o)))

		// log.Println(df[0].Description)
		var documents []txtai.Document
		td := txtai.Document{Id: strings.TrimSuffix(oa[item], path.Ext(oa[item])), Text: desc}
		documents = append(documents, td)
		embeddings.Add(documents)
	}

	log.Println("Calling indexing, this will take some time. ")
	embeddings.Index()

	return nil, nil
}

func httpsDesc(b []byte) (string, error) {
	proc, options := jsonld.JLDProc()

	// "@type":       "Dataset",
	frame := map[string]interface{}{
		"@context":    map[string]interface{}{"@vocab": "https://schema.org/"},
		"@explicit":   true,
		"@type":       "Dataset",
		"description": map[string]interface{}{},
	}

	var myInterface interface{}
	err := json.Unmarshal(b, &myInterface)
	if err != nil {
		log.Println("Error when transforming JSON-LD document to interface:", err)
	}

	framedDoc, err := proc.Frame(myInterface, frame, options) // do I need the options set in order to avoid the large context that seems to be generated?
	if err != nil {
		log.Println("Error when trying to frame document", err)
	}

	graph := framedDoc["@graph"]
	// log.Printf("%s : %s \n", objs["bucket"], jo2)
	// ld.PrintDocument("JSON-LD frame succeeded", framedDoc)
	// ld.PrintDocument("JSON-LD graph section", graph) // debug print....

	jsonm, err := json.MarshalIndent(graph, "", " ")
	if err != nil {
		log.Println("Error trying to marshal data", err)
	}

	df := []DataFrame{}
	json.Unmarshal(jsonm, &df)

	if len(df) > 0 {
		return df[0].Description, nil
	}

	return "", nil
}

func httpDesc(b []byte) (string, error) {
	proc, options := jsonld.JLDProc()

	// "@type":       "Dataset",
	frame := map[string]interface{}{
		"@context":    map[string]interface{}{"@vocab": "http://schema.org/"},
		"@explicit":   true,
		"@type":       "Dataset",
		"description": map[string]interface{}{},
	}

	var myInterface interface{}
	err := json.Unmarshal(b, &myInterface)
	if err != nil {
		log.Println("Error when transforming JSON-LD document to interface:", err)
	}

	framedDoc, err := proc.Frame(myInterface, frame, options) // do I need the options set in order to avoid the large context that seems to be generated?
	if err != nil {
		log.Println("Error when trying to frame document", err)
	}

	graph := framedDoc["@graph"]
	// log.Printf("%s : %s \n", objs["bucket"], jo2)
	// ld.PrintDocument("JSON-LD frame succeeded", framedDoc)
	// ld.PrintDocument("JSON-LD graph section", graph) // debug print....

	jsonm, err := json.MarshalIndent(graph, "", " ")
	if err != nil {
		log.Println("Error trying to marshal data", err)
	}

	df := []DataFrame{}
	json.Unmarshal(jsonm, &df)

	if len(df) > 0 {
		return df[0].Description, nil
	}

	return "", nil
}
