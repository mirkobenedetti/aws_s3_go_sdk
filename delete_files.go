package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the name of your Bucket: ")
	bucket, _ := reader.ReadString('\n')

	fmt.Print("Enter the name of the file to delete: ")
	filename, _ := reader.ReadString('\n')

	for _, delimiter := range []string{"\n", "\r"} {
		bucket = strings.TrimSuffix(bucket, delimiter)
		filename = strings.TrimSuffix(filename, delimiter)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	svc := s3.New(sess)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(filename)})
	if err != nil {
		exitErrorf("Unable to delete object %q from bucket %q, %v", filename, bucket, err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		exitErrorf("Error occurred while waiting for object %q to be deleted, %v", filename)
	}

	fmt.Printf("Object %q successfully deleted\n", filename)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
