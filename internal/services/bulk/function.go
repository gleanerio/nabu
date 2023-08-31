package bulk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gleanerio/nabu/pkg/config"

	"github.com/gleanerio/nabu/internal/graph"

	"github.com/gleanerio/nabu/internal/objects"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/minio/minio-go/v7"
)

// BulkLoad
// This functions could be used to load stored release graphs to the graph database
func BulkLoad(v1 *viper.Viper, mc *minio.Client, bucketName string, item string) (string, error) {
	//spql, err := config.GetSparqlConfig(v1)
	//if err != nil {
	//	return "", err
	//}
	epflag := v1.GetString("flags.endpoint")
	spql, err := config.GetEndpoint(v1, epflag, "bulk")
	if err != nil {
		log.Error(err)
	}
	ep := spql.URL
	md := spql.Method
	ct := spql.Accept

	// check for the required bulk endpoint, no need to move on from here
	if spql.URL == "" {
		return "", errors.New("The configuration file lacks an endpointBulk entry")
	}

	log.Printf("Object %s:%s for %s with method %s type %s", bucketName, item, ep, md, ct)

	b, _, err := objects.GetS3Bytes(mc, bucketName, item)
	if err != nil {
		return "", err
	}

	// NOTE:   commented out, but left.  Since we are loading quads, no need for a graph.
	// If (when) we add back in ntriples as a version, this could be used to build a graph for
	// All the triples in the bulk file to then load as triples + general context (graph)
	// Review if this graph g should b here since we are loading quads
	// I don't think it should b.   validate with all the tested triple stores
	//bn := strings.Replace(bucketName, ".", ":", -1) // convert to urn : values, buckets with . are not valid IRIs
	g, err := graph.MakeURN(v1, item)
	if err != nil {
		log.Error("gets3Bytes %v\n", err)
		return "", err // Assume return. since on this error things are not good?
	}
	url := fmt.Sprintf("%s?graph=%s", ep, g)

	// check if JSON-LD and convert to RDF
	if strings.Contains(item, ".jsonld") {
		nb, err := graph.JSONLDToNQ(v1, string(b))
		if err != nil {
			return "", err
		}
		b = []byte(nb)
	}

	req, err := http.NewRequest(md, url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", ct) // needs to be x-nquads for blaze, n-quads for jena and graphdb
	req.Header.Set("User-Agent", "EarthCube_DataBot/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	log.Println(resp)
	body, err := io.ReadAll(resp.Body) // return body if you want to debugg test with it
	if err != nil {
		log.Println(string(body))
		return string(body), err
	}

	// report
	log.Println(string(body))
	log.Printf("success: %s : %d  : %s\n", item, len(b), ep)

	return string(body), err
}
