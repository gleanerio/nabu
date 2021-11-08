package graph

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/knakk/rdf"
	"github.com/piprate/json-gold/ld"
)

//NQtoNTCtx converts nquads to ntriples plus a context (graph) string
func NQToNTCtx(inquads string) (string, string, error) {
	// loop on tr and make a set of triples
	ntr := []rdf.Triple{}
	g := ""

	dec := rdf.NewQuadDecoder(strings.NewReader(inquads), rdf.NQuads)
	tr, err := dec.DecodeAll()
	if err != nil {
		log.Printf("Error decoding triples: %v\n", err)
		return "", g, err
	}

	// check we have triples
	if len(tr) < 1 {
		return "", g, errors.New("no triple")
	}

	for i := range tr {
		ntr = append(ntr, tr[i].Triple)
	}

	// Assume context of first triple is context of all triples  (again, a bit of a hack,
	// but likely valid as a single JSON-LD datagraph level).  This may be problematic for a "stitegraphs" where several
	// datagraph are represented in a single large JSON-LD via some collection concept.  There it is possible someone might
	// use the quad.  However, for most cases the quad is not important to us, it's local provenance, so we would still replace
	// it with our provenance (context)
	ctx := tr[0].Ctx
	g = ctx.String()

	// TODO output
	outtriples := ""
	buf := bytes.NewBufferString(outtriples)
	enc := rdf.NewTripleEncoder(buf, rdf.NTriples)
	err = enc.EncodeAll(ntr)
	if err != nil {
		log.Printf("Error encoding triples: %v\n", err)
	}
	enc.Close()

	tb := bytes.NewBuffer([]byte(""))
	for k := range ntr {
		tb.WriteString(ntr[k].Serialize(rdf.NTriples))
	}

	return tb.String(), g, err
}

// JSONLDToNQ takes JSON-LD and convets to nqquads (or ntriples if no graph?)
func JSONLDToNQ(jsonld string) (string, error) {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	// add the processing mode explicitly if you need JSON-LD 1.1 features
	options.ProcessingMode = ld.JsonLd_1_1
	options.Format = "application/n-quads"

	var myInterface interface{}
	err := json.Unmarshal([]byte(jsonld), &myInterface)
	if err != nil {
		log.Println("Error when transforming JSON-LD document to interface:", err)
		return "", err

	}

	triples, err := proc.ToRDF(myInterface, options) // returns triples but toss them, just validating
	if err != nil {
		log.Println("Error when transforming JSON-LD document to RDF:", err)
		return "", err

	}

	return fmt.Sprintf("%v", triples), err
}

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
