package main

import (
	"github.com/aws/aws-sdk-go/aws"
	//  "github.com/aws/aws-sdk-go/aws/awsutil"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"sort"
)

type s3items struct {
	Name string
	Size int64
}
type results []s3items

func (slice results) Len() int {
	return len(slice)
}

func (slice results) Less(i, j int) bool {
	return slice[i].Size < slice[j].Size
}

func (slice results) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func main() {
	bucket := ""
	if len(os.Args) < 2 {
		fmt.Println("Usage: s3ls {bucketname}\nError: Missing Bucket Name")
		os.Exit(1)
	} else {
		bucket = os.Args[1]
	}

	svc := s3.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket), // Required
	}
	resp, err := svc.ListObjects(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Generic AWS error with Code, Message, and original error (if any)
			fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				// A service error occurred
				fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
			}
		} else {
			// This case should never be hit, the SDK should always return an
			// error which satisfies the awserr.Error interface.
			fmt.Println(err.Error())
		}
	}

	// Pretty-print the response data.
	var list results
	for e := 0; e < len(resp.Contents); e++ {
		if *resp.Contents[e].Size != 0 {
			var temp s3items
			temp.Name, temp.Size = *resp.Contents[e].Key, *resp.Contents[e].Size
			list = append(list, temp)
			//fmt.Println(*resp.Contents[e].Size, *resp.Contents[e].Key)
		}
	}
	sort.Sort(list)
	for e := 0; e < len(list); e++ {
		fmt.Println(list[e].Size, list[e].Name)
	}
	fmt.Printf("Bucket has %d items\n", len(list))
}
