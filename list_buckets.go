package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	if len(result.Buckets) > 0 {
		fmt.Println("Here are your Buckets:")
		showBuckets(svc, result)
	} else {
		fmt.Println("You have no Buckets yet.")
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func showBuckets(svc *s3.S3, result *s3.ListBucketsOutput) {
	for _, b := range result.Buckets {

		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))

		resp, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(*b.Name)})
		if err != nil {
			exitErrorf("Unable to list items in bucket %q, %v", *b.Name, err)
		}

		for _, item := range resp.Contents {
			fmt.Println("Name:         ", *item.Key)
			fmt.Println("Last modified:", *item.LastModified)
			fmt.Println("Size:         ", *item.Size)
			fmt.Println("Storage class:", *item.StorageClass)
			fmt.Println("")
		}

		fmt.Println("Found", len(resp.Contents), "items in bucket", *b.Name)
		fmt.Println("")

	}
}
