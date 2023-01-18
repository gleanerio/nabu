package prune

import (
	"fmt"
	"github.com/gleanerio/nabu/internal/graph"
	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/pkg/config"
	"github.com/minio/minio-go/v7"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Snip removes graphs in TS not in object store
func Snip(v1 *viper.Viper, mc *minio.Client) error {
	var pa []string
	//err := v1.UnmarshalKey("objects.prefix", &pa)
	objs, err := config.GetObjectsConfig(v1)
	bucketName, _ := config.GetBucketName(v1)
	if err != nil {
		log.Println(err)
	}
	pa = objs.Prefix

	fmt.Println(objs)

	for p := range pa {

		// do the object assembly
		oa, err := ObjectList(v1, mc, pa[p])
		if err != nil {
			log.Println(err)
			return err
		}

		// collect all the graphs from triple store
		// TODO resolve issue with Graph and graphList vs graphListStatements
		//ga, err := graphListStatements(v1, mc, pa[p])
		ga, err := graphList(v1, mc, pa[p])
		if err != nil {
			log.Println(err)
			return err
		}

		//objs := v1.GetStringMapString("objects") // from above
		// convert the object names to the URN pattern used in the graph
		oag := []string{} // object array graph mode
		for x := range oa {
			//s := strings.TrimSuffix(oa[x], ".rdf")
			//s2c := strings.Replace(s, "/", ":", -1)
			g, err := graph.MakeURN(oa[x], bucketName)
			if err != nil {
				log.Error("gets3Bytes %v\n", err)
				// should this just return. since on this error things are not good
			}
			oag = append(oag, g)
		}

		//compare lists, anything IN graph not in objects list should be removed
		d := difference(ga, oag) // return array of items in ga that are NOT in oa
		m := findMissingElements(oag, ga)

		fmt.Printf("Graph items: %d  Object items: %d  difference: %d\n", len(ga), len(oag), len(d))
		fmt.Printf("Missing item count: %d\n", len(m))

		// For each in d will delete that graph
		if len(d) > 0 {
			bar := progressbar.Default(int64(len(d)))
			for x := range d {
				log.Printf("Remove graph: %s\n", d[x])
				graph.Drop(v1, d[x])
				bar.Add(1)
			}
		}

		// load new ones..
		spql, err := config.GetSparqlConfig(v1)
		if err != nil {
			log.Error("prune -> config.GetSparqlConfig %v\n", err)
		}

		if len(m) > 0 {
			bar2 := progressbar.Default(int64(len(m)))
			log.Info("uploading missing %n objects", m)
			for x := range m {
				np, _ := graph.URNToPrefix(m[x], "summoned", ".jsonld")
				log.Tracef("Add graph: %s  %s \n", m[x], np)
				_, err := objects.PipeLoad(v1, mc, bucketName, np, spql.Endpoint)
				if err != nil {
					log.Error("prune -> pipeLoad %v\n", err)
				}
				bar2.Add(1)
			}
		}

	}

	return nil
}
