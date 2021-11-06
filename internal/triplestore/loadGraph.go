package triplestore

import (
	"fmt"
	"strings"

	"github.com/gleanerio/nabu/internal/graph"
	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gosuri/uiprogress"
	"github.com/minio/minio-go/v7"
)

// TODO  review and ensure this is not used and then remove.  This functionality is better done in something like Nabu or a groups own approach
// LoadGraph attempts to load data from srcbkt into graph.  It assumes the data is stored in JSON-LD
func DEPRECATEDLoadGraph(mc *minio.Client, srcbckt string, sue string) error {
	var err error

	// Get the objects...   may need to add filter on .jsonld ?
	ea := objects.GetObjects(mc, srcbckt)

	count := len(ea)
	// log.Println(count)
	// uiprogress.Start() // start rendering
	bar2 := uiprogress.AddBar(count).AppendCompleted().PrependElapsed()
	bar2.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("RDF load: %s (%d/%d)", srcbckt, b.Current(), count)
	})

	for _, v := range ea {
		bar2.Incr()

		// fmt.Printf("%s %s \n", v.Bucketname, v.Key)

		// NOTE
		// For JRSO we moved to DO and DO metadatain the same
		// bucket with a .jsonld extension on metadata.  So
		// break here if don't have that
		// if !strings.Contains(v.Key, ".jsonld") {
		// 	continue
		// }

		b, _, err := objects.GetS3Bytes(mc, v.Bucketname, v.Key)
		if err != nil {
			fmt.Printf("gets3Bytes %v\n", err)
		}
		// fmt.Println(string(b))

		nq, err := graph.JSONLDToNQ(string(b))
		if err != nil {
			fmt.Printf("Error in JSONLDToNQ %v \n", err)
		}
		// fmt.Println(nq)

		// BUG FIX HACK
		// I typo-ed schema.org/AdditonalType as schema.org/AdditionType
		// here I could do a simple find and replace to resolve that.   Then remove
		// later when I hide my shame (it's been fixed in VaultWalker FYI)
		nqhacked := strings.ReplaceAll(nq, "additionType", "additionalType")

		_, err = BlazeUpdateNQ([]byte(nqhacked), sue) // TODO  add the above graph string to target a database on jena

		// _, err = triplestore.JenaUpdateNQ([]byte(nqhacked), sue) // TODO  add the above graph string to target a database on jena

		if err != nil {
			fmt.Printf("Error in update call: %v\n", err)
		}
		// fmt.Printf("graphUploader: %s \n", string(r))

		// TODO  review /home/fils/src/Projects/LDN/GoLDeN/internal/graph and incorporate here...??
	}

	return err // why do I return OK  :)
}
