package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
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
    return slice[i].Size < slice[j].Size;
}

func (slice results) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}


func main() {
  bucket := ""
  region := ""
  if len(os.Args) == 1  {
    fmt.Println("Usage: s3ls {bucketname} [region]\nError: Missing Bucket Name")
    os.Exit(1)
  } else {
    bucket = os.Args[1]
    if len(os.Args) == 3 {
      region = os.Args[2]
    } else {
      region = "us-east-1"
    }
  }
	svc := s3.New(&aws.Config{Region: aws.String(region)})
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
	for e := range resp.Contents {
		if *resp.Contents[e].Size != 0 {
      var temp s3items
      temp.Name, temp.Size = *resp.Contents[e].Key, *resp.Contents[e].Size
      list = append(list,temp)
			//fmt.Println(*resp.Contents[e].Size, *resp.Contents[e].Key)
		}
	}
  sort.Sort(list)
  var totalsize int64
  for e:= range list {
    totalsize = totalsize + list[e].Size
    //fmt.Println(list[e].Size,list[e].Name)
    fmt.Printf("%-15d%s\n", list[e].Size,list[e].Name)
  }
  fmt.Printf("Bucket has %d items.\n", len(list))
  fmt.Printf("%.2f Megabytes total space used.\n", float64(totalsize)/1000000)
}
