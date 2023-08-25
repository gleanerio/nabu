package graph

import (
	"fmt"
	"strings"
)

// URNToPrefix convert  urn string to a valid s3 prefix
// So
//
//	urn:gleaner.iow:iow:counties0:1bbde08c130150c48e7d424a60a00d574174bd28
//
// becomes
//
//	gleaner.iow/summoned/counties0/1bbde08c130150c48e7d424a60a00d574174bd28f.jsonld
func URNToPrefix(urn, pathelement, suffix string) (string, error) {
	parts := strings.Split(urn, ":")

	// Use the append function to create a new slice that includes the new element
	//pos := 3 // insert path element at this position
	//np := append(parts[:pos], append([]string{pathelement}, parts[pos:]...)...)

	parts[2] = "summoned"
	np := parts

	// modify last element (object) with suffix
	index := len(np) - 1
	np[index] = np[index] + suffix

	// drop first element (urn)
	np = np[2:]

	objectPrefix := fmt.Sprintf("/%s", strings.Join(np, "/"))

	//log.Print(objectPrefix)

	return objectPrefix, nil
}
