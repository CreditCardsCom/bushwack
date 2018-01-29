package main

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/CreditCardsCom/bushwack/bushwack"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var esUrl string

func init() {
	host := os.Getenv("ES_HOST")
	if host == "" {
		log.Fatalln("No ES_HOST defined in the environment!")
	}

	esUrl = fmt.Sprintf("%s/_bulk", host)
}

func main() {
	lambda.Start(EventHandler)
}

func EventHandler(event events.S3Event) {
	// TODO: Use S3Event region instead of hardcoded
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})
	if err != nil {
		log.Fatal(err)
	}

	dlManager := s3manager.NewDownloader(sess)
	var objects []s3manager.BatchDownloadObject
	var files []string

	// Build a batch download of all objects send in the event
	for _, r := range event.Records {
		e := r.S3

		// Verify we have what we thing we have...
		if !strings.HasSuffix(e.Object.Key, ".log.gz") {
			continue
		}

		tf, err := ioutil.TempFile("", "")
		if err != nil {
			log.Fatal(err)
		}
		defer tf.Close()
		defer os.Remove(tf.Name())

		o := s3manager.BatchDownloadObject{
			Writer: tf,
			Object: &s3.GetObjectInput{
				Bucket: &e.Bucket.Name,
				Key:    &e.Object.Key,
			},
		}

		files = append(files, tf.Name())
		objects = append(objects, o)
		log.Printf("Processing s3://%s/%s\n", e.Bucket.Name, e.Object.Key)
	}

	it := &s3manager.DownloadObjectsIterator{Objects: objects}
	ctx := aws.BackgroundContext()
	if err := dlManager.DownloadWithIterator(ctx, it); err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		processLog(f)
	}
}

func processLog(filename string) {
	contents, err := decompress(filename)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := bushwack.ParseLog(string(contents))
	if err != nil {
		log.Fatal(err)
	}

	if len(entries) == 0 {
		log.Println("No entries to index.")
		return
	}

	body, err := entries.SerializeBulkBody()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(esUrl, "application/x-ndjson", strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Printf("Received unexpected status code %d.\n", resp.StatusCode)

		reason, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Reason: %s", reason)
	}

	log.Printf("Sent off %d log entries.", len(entries))
}

func decompress(f string) ([]byte, error) {
	fd, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	r, err := gzip.NewReader(fd)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}
