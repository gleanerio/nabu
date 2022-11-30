package graph

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"

	"github.com/piprate/json-gold/ld"
)

// ContextMapping holds the JSON-LD mappings for cached context
type ContextMapping struct {
	Prefix string
	File   string
}

// TODO   we create this all the time..  stupidly..  Generate these pointers
// and pass them around, don't keep making it over and over
// Ref:  https://schema.org/docs/howwework.html and https://schema.org/docs/jsonldcontext.json

// JLDProc builds the JSON-LD processor and sets the options object
// for use in framing, processing and all JSON-LD actions
func JLDProc() (*ld.JsonLdProcessor, *ld.JsonLdOptions) { // TODO make a booklean
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	client := &http.Client{}
	nl := ld.NewDefaultDocumentLoader(client)

	m := make(map[string]string)

	f := "./web/jsonldcontext.json"
	if fileExists(f) {
		m["http://schema.org/"] = f
	} else {
		log.Printf("Could not find: %s", f)
	}

	f = "./web/jsonldcontext.json"
	if fileExists(f) {
		m["https://schema.org/"] = f
	} else {
		log.Printf("Could not find: %s", f)
	}

	// Read mapping from config file
	cdl := ld.NewCachingDocumentLoader(nl)
	cdl.PreloadWithMapping(m)
	options.DocumentLoader = cdl

	// TODO let this be set later via config
	// Set to a default format..
	options.Format = "application/nquads"

	return proc, options
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
