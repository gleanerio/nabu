package objects

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"log"
)

// PutS3Bytes is used write an object
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
