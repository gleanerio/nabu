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

// docfunc needs to be renamed to something like BulkUpate and made exported
// This functions could be used to load stored release graphs to the graph database
func docfunc(v1 *viper.Viper, mc *minio.Client, bucketName string, item string) (string, error) {
	spql, err := config.GetSparqlConfig(v1)
	if err != nil {
		return "", err
	}
	ep := spql.EndpointBulk
	md := spql.EndpointMethod
	ct := spql.ContentType

	// check for the required bulk endpoint, no need to move on from here
	if spql.EndpointBulk == "" {
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
	bn := strings.Replace(bucketName, ".", ":", -1) // convert to urn : values, buckets with . are not valid IRIs
	g, err := graph.MakeURN(item, bn)
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
