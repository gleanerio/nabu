package context

import (
	"log"
	"net/http"
	"os"

	"github.com/piprate/json-gold/ld"
)

// ContextMapping holds the JSON-LD mappings for cached context
type ContextMapping struct {
	Prefix string
	File   string
}

// JLDProc build the JSON-LD processer and sets the options object
// to use in framing, processing and all JSON-LD actions
// TODO   we create this all the time..  stupidly..  Generate these pointers
// and pass them around, don't keep making it over and over
// Ref:  https://schema.org/docs/howwework.html and https://schema.org/docs/jsonldcontext.json
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

	// Set a default format..  let this be set later...
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
