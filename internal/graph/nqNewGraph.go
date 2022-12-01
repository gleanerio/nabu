package graph

import (
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"
	"strings"

	"github.com/knakk/rdf"
)

//NQNewGraph converts nquads to nquads with a new context graph
func NQNewGraph(inquads, newctx string) (string, string, error) {

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

	//ntr = append(ntr, )

	// Assume context of first triple is context of all triples  (again, a bit of a hack,
	// but likely valid as a single JSON-LD datagraph level).  This may be problematic for a "stitegraphs" where several
	// datagraph are represented in a single large JSON-LD via some collection concept.  There it is possible someone might
	// use the quad.  However, for most cases the quad is not important to us, it's local provenance, so we would still replace
	// it with our provenance (context)
	ctx := tr[0].Ctx
	g = ctx.String()

	// TODO make a new Ctx from newCtx and repate the existing one

	// TODO update the following to output quads

	// TODO output
	outtriples := ""
	buf := bytes.NewBufferString(outtriples)
	//enc := rdf.NewQuadEncoder(buf, rdf.NQuads)

	enc := rdf.NewTripleEncoder(buf, rdf.NTriples)

	err = enc.EncodeAll(ntr)
	if err != nil {
		log.Printf("Error encoding triples: %v\n", err)
	}
	enc.Close()

	tb := bytes.NewBuffer([]byte(""))
	for k := range ntr {
		tb.WriteString(ntr[k].Serialize(rdf.NQuads))
	}

	return tb.String(), g, err
}
