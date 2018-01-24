package main

import (
	"compress/gzip"
	"fmt"
	//"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/s3"
	"github.com/CreditCardsCom/bushwack/bushwack"
	"github.com/Jeffail/tunny"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const urlFormat = "<Elastic search endpoint>"

func main() {
	//sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})
	//if err != nil {
	//	log.Fatal(err)
	//}

	//service := s3.New(sess)
	//s3Resp, err := service.ListObjects(&s3.ListObjectsInput{
	//	Bucket: aws.String("cccom-elb-logs"),
	//	Prefix: aws.String("transnode/AWSLogs/799335648850/elasticloadbalancing/us-west-2/2017/01/23/"),
	//})

	//if err != nil {
	//	log.Fatal(err)
	//}

	//for _, item := range s3Resp.Contents {
	//	fmt.Println("Name:         ", *item.Key)
	//	fmt.Println("Last modified:", *item.LastModified)
	//	fmt.Println("Size:         ", *item.Size)
	//	fmt.Println("Storage class:", *item.StorageClass)
	//	fmt.Println("")
	//}

	pool := tunny.NewFunc(20, func(p interface{}) interface{} {
		processLog(p.(string))

		return nil
	})
	defer pool.Close()

	files, err := ioutil.ReadDir("logs")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			p := filepath.Join("logs", f.Name())

			pool.Process(p)
		}
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

	body, err := entries.SerializeBulkBody()
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf(urlFormat, "2018.01.23")
	resp, err := http.Post(url, "application/x-ndjson", strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Println("Yikes, looks like we had a request return a non-200 status")
	}

	log.Printf("Sent off %d log entires.", len(entries))
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
