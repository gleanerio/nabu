package graph

import (
	"github.com/piprate/json-gold/ld"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

// ContextMapping holds the JSON-LD mappings for cached context
type ContextMapping struct {
	Prefix string
	File   string
}

// JLDProc builds the JSON-LD processor and sets the options object
// for use in framing, processing and all JSON-LD actions
func JLDProc(v1 *viper.Viper) (*ld.JsonLdProcessor, *ld.JsonLdOptions) {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	//spql, _ := config.GetSparqlConfig(v1)
	// cntxmap, _ := config.GetContextMapConfig(v1)
	// TODO, modled after above, need a "contextmap:" in the config file
	// with several KV pairs like
	//	contextmap:
	//		- key: https://schema.org/
	//		  value: ./assets/schemaorg-current-https.jsonld
	//		- key: http://schema.org/
	//		  value: ./assets/schemaorg-current-http.jsonld

	client := &http.Client{}
	nl := ld.NewDefaultDocumentLoader(client)

	m := make(map[string]string)

	// remove the hardcoded location (see TODO above)
	f := "./assets/schemaorg-current-http.jsonld"
	if fileExists(f) {
		m["http://schema.org/"] = f
	} else {
		log.Printf("Could not find: %s", f)
	}

	f = "./assets/schemaorg-current-https.jsonld"
	if fileExists(f) {
		m["https://schema.org/"] = f
	} else {
		log.Printf("Could not find: %s", f)
	}

	// Read mapping from config file
	cdl := ld.NewCachingDocumentLoader(nl)
	cdl.PreloadWithMapping(m)
	options.DocumentLoader = cdl

	options.ProcessingMode = ld.JsonLd_1_1 // add mode explicitly if you need JSON-LD 1.1 features
	options.Format = "application/nquads"  // Set to a default format. (make an option?)

	return proc, options
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
