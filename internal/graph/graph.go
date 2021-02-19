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

func NQToNTCtx(inquads string) (string, string, error) {
	// loop on tr and make a set of triples
	ntr := []rdf.Triple{}
	g := ""

	// HACK, delete thie..  was used to add missing . from tika index output (that is fixed now)
	// dec := rdf.NewQuadDecoder(strings.NewReader(fmt.Sprintf("%s .", inquads)), rdf.NQuads)
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
	ctx := tr[0].Ctx // Assume context of first triple is context of all triples
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

// IILTriple builds a IRI, IRI, Literal triple
func IILTriple(s, p, o, c string) (string, error) {
	buf := bytes.NewBufferString("")

	newctx, err := rdf.NewIRI(c) // this should be  c
	if err != nil {
		return buf.String(), err
	}
	ctx := rdf.Context(newctx)

	sub, err := rdf.NewIRI(s)
	if err != nil {
		log.Println("Error building subject IRI for tika triple")
		return buf.String(), err
	}
	pred, err := rdf.NewIRI(p)
	if err != nil {
		log.Println("Error building predicate IRI for tika triple")
		return buf.String(), err
	}
	obj, err := rdf.NewLiteral(o)
	if err != nil {
		log.Println("Error building object literal for tika triple")
		return buf.String(), err
	}

	t := rdf.Triple{Subj: sub, Pred: pred, Obj: obj}
	q := rdf.Quad{t, ctx}

	qs := q.Serialize(rdf.NQuads)
	if s != "" && p != "" && o != "" {
		fmt.Fprintf(buf, "%s", qs)
	}
	return buf.String(), err
}
