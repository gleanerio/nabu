package graph

import (
	"errors"
	"fmt"
	"strings"
)

func MakeURN(s, bucket string) (string, error) {

	var (
		g   string // build the URN for the graph context string we use
		err error
	)

	sr := strings.Replace(s, "/", ":", -1)
	s2c := getLastTwo(sr) // split the string and take last two segments

	if strings.Contains(s2c, ".rdf") {
		g = fmt.Sprintf("urn:%s:%s", bucket, strings.TrimSuffix(s2c, ".rdf"))
	} else if strings.Contains(s2c, ".jsonld") {
		g = fmt.Sprintf("urn:%s:%s", bucket, strings.TrimSuffix(s2c, ".jsonld"))
	} else if strings.Contains(s2c, ".nq") {
		g = fmt.Sprintf("urn:%s:%s", bucket, strings.TrimSuffix(s2c, ".nq"))
	} else {
		err = errors.New("unable to generate graph URI")
	}

	return g, err
}

// getLastTwo from chatGPT
func getLastTwo(s string) string {
	// Split the string on the ":" character.
	parts := strings.Split(s, ":")

	// Get the last two elements.
	lastTwo := parts[len(parts)-2:]

	// Join the last two elements and return the result.
	return strings.Join(lastTwo, ":")
}
