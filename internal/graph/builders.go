package graph

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/knakk/rdf"
)

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
