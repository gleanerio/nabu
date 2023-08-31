package graph

import (
	"errors"
	"fmt"
	"github.com/gleanerio/nabu/pkg/config"
	"github.com/spf13/viper"
	"strings"
)

// MakeURN formats a URN following the ADR 0001-URN-decision.md  which at the
// time of this coding resulted in     urn:{program}:{organization}:{provider}:{sha}
func MakeURN(v1 *viper.Viper, s string) (string, error) {
	gcfg, _ := config.GetImplNetworkConfig(v1)

	var (
		g   string // build the URN for the graph context string we use
		err error
	)

	sr := strings.Replace(s, "/", ":", -1)
	s2c := getLastTwo(sr) // split the string and take last two segments

	if strings.Contains(s2c, ".rdf") {
		g = fmt.Sprintf("urn:%s:%s", gcfg.Orgname, strings.TrimSuffix(s2c, ".rdf"))
	} else if strings.Contains(s2c, ".jsonld") {
		g = fmt.Sprintf("urn:%s:%s", gcfg.Orgname, strings.TrimSuffix(s2c, ".jsonld"))
	} else if strings.Contains(s2c, ".nq") {
		g = fmt.Sprintf("urn:%s:%s", gcfg.Orgname, strings.TrimSuffix(s2c, ".nq"))
	} else {
		err = errors.New("unable to generate graph URI")
	}

	return g, err
}

// MakeURNPrefix formats a URN following the ADR 0001-URN-decision.md  which at the
// time of this coding resulted in     urn:{program}:{organization}:{provider}:{sha}
// the "prefix" version only returns the prefix part of the urn, for use in the prune
// command
func MakeURNPrefix(v1 *viper.Viper, prefix string) (string, error) {
	gcfg, _ := config.GetImplNetworkConfig(v1)

	var (
		g   string // build the URN for the graph context string we use
		err error
	)

	ps := strings.Split(prefix, "/")

	g = fmt.Sprintf("urn:%s:%s", gcfg.Orgname, ps[len(ps)-1])

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
