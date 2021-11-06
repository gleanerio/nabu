package objects

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gleanerio/nabu/run/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"log"
)

type Entry struct {
	Bucketname string
	Key        string
	Urlval     string
	Sha1val    string
	Jld        string
}

// MinioConnection Set up minio and initialize client
func MinioConnection(v1 *viper.Viper) (*minio.Client, error) {
	//mcfg := v1.GetStringMapString("minio")

	mcfg, err := config.GetMinioConfig(v1)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//endpoint := fmt.Sprintf("%s:%s", mcfg["address"], mcfg["port"])
	//accessKeyID := mcfg["accesskey"]
	//secretAccessKey := mcfg["secretkey"]
	//useSSL, err := strconv.ParseBool(fmt.Sprintf("%s", mcfg["useSSL"]))
	endpoint := fmt.Sprintf("%s:%d", mcfg.Address, mcfg.Port)
	accessKeyID := mcfg.Accesskey
	secretAccessKey := mcfg.Secretkey
	useSSL := mcfg.Ssl

	// minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, true)
	minioClient, err := minio.New(endpoint, &minio.Options{Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""), Secure: useSSL})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	// minioClient.SetCustomTransport(&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}})
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	return minioClient, err
}

// GetS3Bytes simply pulls the byes of an object from the store
func GetS3Bytes(mc *minio.Client, bucket, object string) ([]byte, string, error) {
	fo, err := mc.GetObject(context.Background(), bucket, object, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	oi, err := fo.Stat()
	if err != nil {
		log.Println("Issue with reading an object..  should I just fatal on this to make sure?")
	}

	// resuri := ""

	// urlval := ""
	// if len(oi.Metadata["X-Amz-Meta-Url"]) > 0 {
	// 	urlval = oi.Metadata["X-Amz-Meta-Url"][0] // also have  X-Amz-Meta-Sha1
	// }

	// shaval := ""
	// if len(oi.Metadata["X-Amz-Meta-Sha256"]) > 0 {
	// 	shaval = oi.Metadata["X-Amz-Meta-Sha256"][0]
	// }

	dgraph := ""
	if len(oi.Metadata["X-Amz-Meta-Dgraph"]) > 0 {
		dgraph = oi.Metadata["X-Amz-Meta-Dgraph"][0]
	}

	// fmt.Printf("%s %s %s \n", urlval, sha1val, resuri)

	sz := oi.Size        // what type is this...
	if sz > 1073741824 { // if bigger than 1 GB (which is very small) move on
		return nil, "", errors.New("file above tika processing size threshhold")
	}

	// TODO..   set an upper byte size  limit here and return error if the size is too big
	buf := new(bytes.Buffer)
	buf.ReadFrom(fo)
	bb := buf.Bytes() // Does a complete copy of the bytes in the buffer.

	return bb, dgraph, err
}

// GetMillObjects
func GetObjects(mc *minio.Client, bucketname string) []Entry {
	doneCh := make(chan struct{}) // Create a done channel to control 'ListObjectsV2' go routine.
	defer close(doneCh)           // Indicate to our routine to exit cleanly upon return.

	// isRecursive := true
	// objectCh := mc.ListObjects(context.Background(), bucketname, "", isRecursive, doneCh) // no v2 for swift

	objectCh := mc.ListObjects(context.Background(), bucketname, minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: true,
	}) // no v2 for swift

	var entries []Entry

	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return nil
		}

		// May need optional check here for .jsonld
		// or object metadata too...

		fo, err := mc.GetObject(context.Background(), bucketname, object.Key, minio.GetObjectOptions{})
		if err != nil {
			fmt.Println(err)
			return nil
		}

		oi, err := fo.Stat()
		if err != nil {
			log.Println("Issue with reading an object..  should I just fatal on this to make sure?")
		}
		urlval := ""
		sha1val := ""
		if len(oi.Metadata["X-Amz-Meta-Url"]) > 0 {
			urlval = oi.Metadata["X-Amz-Meta-Url"][0] // also have  X-Amz-Meta-Sha1
		}
		if len(oi.Metadata["X-Amz-Meta-Sha1"]) > 0 {
			sha1val = oi.Metadata["X-Amz-Meta-Sha1"][0]
		}

		// buf := new(bytes.Buffer)
		// buf.ReadFrom(fo)
		// jld := buf.String() // Does a complete copy of the bytes in the buffer.

		// Mock call for some validation (and a template for other millers)
		// Mock(bucketname, object.Key, urlval, sha1val, jld)
		entry := Entry{Bucketname: bucketname, Key: object.Key, Urlval: urlval, Sha1val: sha1val}
		entries = append(entries, entry)

	}

	fmt.Println(len(entries))
	// multiCall(entries)

	return entries
}

func PutS3Bytes(mc *minio.Client, bucketName, objectName, mimeType string, object []byte) (int, error) {
	usermeta := make(map[string]string) // what do I want to know?
	usermeta["url"] = "urlloc"
	usermeta["sha1"] = "bss"

	// Upload the file with FPutObject
	n, err := mc.PutObject(context.Background(), bucketName, objectName, bytes.NewReader(object), int64(len(object)), minio.PutObjectOptions{ContentType: mimeType, UserMetadata: usermeta})
	if err != nil {
		log.Printf("%s", objectName)
		log.Fatalln(err)
	}
	log.Printf("Uploaded Bucket:%s File:%s Size %d\n", bucketName, objectName, n)

	return 0, nil
}
