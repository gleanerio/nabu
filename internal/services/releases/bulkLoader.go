package releases

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"path"
	"strings"
	"time"

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

	//ol, err := common.ObjectList(v1, mc, "graphs/latest")
	//if err != nil {
	//	return err
	//}

	// Let's move the current bulk graph to archive and clear the way for a new release graph
	//for o := range ol {
	//	for p := range pa {
	//		sp := strings.Split(pa[p], "/")
	//		spj := strings.Join(sp, "")
	//		if strings.Contains(ol[o], baseName(path.Base(spj))) {
	//			// move the match from graphs/latest to graphs/archive
	//			fmt.Println(ol[o])
	//			// copy it and change the prefix path from "latest" to "archive"
	//			err = objects.Copy(v1, mc, bucketName, ol[o], bucketName, strings.Replace(ol[o], "latest", "archive", 1))
	//			if err != nil {
	//				log.Println(err)
	//				return err
	//			}
	//			// remove it
	//			err = objects.Remove(v1, mc, bucketName, ol[o])
	//			if err != nil {
	//				log.Println(err)
	//				return err
	//			}
	//		}
	//	}
	//}

	for p := range pa {
		sp := strings.Split(pa[p], "/")
		sourcename := strings.Join(sp[1:], "__")
		spj := strings.Join(sp, "__")
		name_latest := fmt.Sprintf("%s_release.nq", baseName(path.Base(sourcename)))

		err = objects.PipeCopy(v1, mc, name_latest, bucketName, pa[p], "graphs/latest") // have this function return the object name and path, easy to load and remove then
		if err != nil {
			return err
		}
		// move a copy from latest to archive
		const layout = "2006-01-02-15-04-05"
		t := time.Now()
		name := fmt.Sprintf("%s/%s/%s_%s_release.nq", "graphs/archive", sourcename, baseName(path.Base(spj)), t.Format(layout))
		latest_fullpath := fmt.Sprintf("%s/%s", "graphs/latest", name_latest)
		err = objects.Copy(v1, mc, bucketName, latest_fullpath, bucketName, strings.Replace(name, "latest", "archive", 1))

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
