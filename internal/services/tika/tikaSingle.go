package tika

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/internal/prune"
	"github.com/schollz/progressbar/v3"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

// SingleBuild is a test function to call a single item
func SingleBuild(v1 *viper.Viper, mc *minio.Client) error {
	objs := v1.GetStringMapString("objects")
	tkcfg := v1.GetStringMapString("tika")
	tb := tkcfg["outbucket"]
	tp := tkcfg["outprefix"]

	var pa []string
	err := v1.UnmarshalKey("objects.prefix", &pa)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(pa)

	for p := range pa {

		// collect the objects...
		oa, err := prune.ObjectList(v1, mc, pa[p])
		if err != nil {
			log.Println(err)
			return err
		}

		fmt.Printf("Object items: %d   \n", len(oa))

		// read the bytes and feed into tika
		bar := progressbar.Default(int64(len(oa)))
		for n := range oa {
			bar.Add(1)

			// check it's not a .jsonld file, we don't tika index those
			if strings.HasSuffix(oa[n], ".jsonld") {
				continue
			}

			// Get the names for the target bucket and prefix we will use, so we can check for them
			// to see if they exist already.
			na := strings.Split(oa[n], "/")
			nl := na[len(na)-1]
			nb := strings.TrimSuffix(nl, filepath.Ext(nl))
			to := fmt.Sprintf("%s/%s.rdf", tp, nb)

			// TODO need a way to skip this check in case we want to force rebuild all files (cfg file option)
			// check we haven't already don't this file
			if ObjectExists(mc, tb, to) == nil {
				continue
			}

			// Get the bytes and generate the triples
			b, _, err := objects.GetS3Bytes(mc, objs["bucket"], oa[n])
			if err != nil {
				fmt.Printf("gets3Bytes %v\n", err)
			}
			s, err := EngineTika(v1, b)
			t, err := fullTextTrpls(s, oa[n])

			bs := bytes.NewBufferString(t)
			_, err = mc.PutObject(context.Background(), tb, to, bs, int64(bs.Len()), minio.PutObjectOptions{})
			if err != nil {
				fmt.Printf("putObject error:  %v\n", err)
			}
		}
	}

	return nil
}

// ObjectExists returns true if the object is found
func ObjectExists(mc *minio.Client, bucket, object string) error {

	_, err := mc.StatObject(context.Background(), bucket, object, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func processObject(v1 *viper.Viper, mc *minio.Client, bucket, prefix string, message minio.ObjectInfo) (string, error) {
	fo, err := mc.GetObject(context.Background(), bucket, message.Key, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("get object %s", err)
		// return "", err
	}

	var b bytes.Buffer
	bw := bufio.NewWriter(&b)

	_, err = io.Copy(bw, fo)
	if err != nil {
		log.Printf("iocopy %s", err)
	}

	s, err := EngineTika(v1, b.Bytes())
	t, err := fullTextTrpls(s, message.Key)

	if err != nil {
		log.Println(err)
	}

	return t, err
}

func fullTextTrpls(s, obj string) (string, error) {
	t := fmt.Sprintf("<https://opencoredata.org/id/%s>  <https://schema.org/text>  \"%s\"  .", obj, s)
	return t, nil
}

// EngineTika sends a byte array to tika for processing into text
func EngineTika(v1 *viper.Viper, b []byte) (string, error) {
	tkcfg := v1.GetStringMapString("tika")
	tikaurl := tkcfg["tikaurl"]

	req, err := http.NewRequest("PUT", tikaurl, bytes.NewReader(b))
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("User-Agent", "EarthCube_DataBot/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	// fmt.Println("Tika Response Status:", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	// Do stop words (make a function)
	//sw := stopwords.CleanString(string(body), "en", true) // remove stop words..   no reason for them in the search
	sw := string(body)

	// remove anything not text and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Println(err)
	}
	processedString := reg.ReplaceAllString(sw, " ")

	// TODO remove duplicate words.. DONE..   but needs review
	//return dedup(processedString), err

	return processedString, err
}

//func dedup(input string) string {
//unique := []string{}

//words := strings.Split(input, " ")
//for _, word := range words {
//// If we alredy have this word, skip.
//if contains(unique, word) {
//continue
//}
//unique = append(unique, word)
//}

//return strings.Join(unique, " ")
//}

//func contains(strs []string, str string) bool {
//for _, s := range strs {
//if s == str {
//return true
//}
//}
//return false
//}
