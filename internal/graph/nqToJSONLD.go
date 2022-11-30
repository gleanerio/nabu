package graph

import (
	"encoding/json"
	"github.com/piprate/json-gold/ld"
	log "github.com/sirupsen/logrus"
)

// NQToJSONLD takes nquads and converts to JSON-LD
func NQToJSONLD(triples string) ([]byte, error) {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	// add the processing mode explicitly if you need JSON-LD 1.1 features
	options.ProcessingMode = ld.JsonLd_1_1

	doc, err := proc.FromRDF(triples, options)
	if err != nil {
		log.Println("ERROR: converting from RDF/NQ to JSON-LD")
		log.Println(err)
	}

	// ld.PrintDocument("JSON-LD output", doc)
	b, err := json.MarshalIndent(doc, "", " ")

	return b, err
}
