package graph

import (
	"fmt"
	"strings"
)

// URNToPrefix convert  urn string to a valid s3 prefix
// So urn:gleaner.oih:edmo:00032788b3d1eecf4257bd8ffd42c5d56761a6bf
// becomes gleaner.oih/[summoned]/edmo/00032788b3d1eecf4257bd8ffd42c5d56761a6bf.jsonld
func URNToPrefix(urn, pathelement, suffix string) (string, error) {
	parts := strings.Split(urn, ":")

	pos := 2 // insert pathelement at this position

	// Use the append function to create a new slice that includes the new element
	np := append(parts[:pos], append([]string{pathelement}, parts[pos:]...)...)

	// modify last element (object) with suffix
	index := len(np) - 1
	np[index] = np[index] + suffix

	// drop first element (urn)
	np = np[2:]

	objectPrefix := fmt.Sprintf("/%s", strings.Join(np, "/"))

	return objectPrefix, nil
}
