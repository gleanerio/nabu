package releases

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"path"
	"strings"
	"time"

	"github.com/gleanerio/nabu/internal/prune"

	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/pkg/config"
	"github.com/spf13/viper"

	"github.com/minio/minio-go/v7"
)

// BulkRelease collects the objects from a bucket to load
func BulkRelease(v1 *viper.Viper, mc *minio.Client) error {
	log.Println("Release:BulkAssembly")
	bucketName, _ := config.GetBucketName(v1)
	objCfg, _ := config.GetObjectsConfig(v1)
	pa := objCfg.Prefix

	var err error

	ol, err := prune.ObjectList(v1, mc, "graphs/latest")
	if err != nil {
		return err
	}

	for o := range ol {
		for p := range pa {
			sp := strings.Split(pa[p], "/")
			spj := strings.Join(sp, "")
			if strings.Contains(ol[o], baseName(path.Base(spj))) {
				// move the match from graphs/latest to graphs/archive
				fmt.Println(ol[o])
				// copy it and change the prefix path from "latest" to "archive"
				err = objects.Copy(v1, mc, bucketName, ol[o], bucketName, strings.Replace(ol[o], "latest", "archive", 1))
				if err != nil {
					log.Println(err)
					return err
				}
				// remove it
				err = objects.Remove(v1, mc, bucketName, ol[o])
				if err != nil {
					log.Println(err)
					return err
				}
			}
		}
	}

	// Set and use a "single file flag" to bypass skolimaization since if it is a single file
	// the JSON-LD to RDF will correctly map blank nodes.
	sf := true
	if len(pa) > 1 {
		sf = false
	}

	for p := range pa {
		sp := strings.Split(pa[p], "/")
		spj := strings.Join(sp, "")
		const layout = "2006-01-02-15-04-05"
		t := time.Now()
		name := fmt.Sprintf("%s_%s_release.nq", baseName(path.Base(spj)), t.Format(layout))

		err = objects.PipeCopy(v1, mc, name, bucketName, pa[p], "graphs/latest", sf) // have this function return the object name and path, easy to load and remove then
		if err != nil {
			return err
		}
	}

	return err
}

func baseName(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n == -1 {
		return s
	}
	return s[:n]
}
