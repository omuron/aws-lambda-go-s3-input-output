// main.go
package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// 出力バケット
const S3Output = "sample-output2"

// https://github.com/aws/aws-lambda-go/blob/main/events/README_S3.md を参考
func handler(ctx context.Context, s3Event events.S3Event) {
	// セッション獲得
	sess := session.Must(session.NewSession())
	// S3 clientを作成
	svc := s3.New(sess)

	for _, record := range s3Event.Records {
		s3rec := record.S3
		log.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3rec.Bucket.Name, s3rec.Object.Key)

		// オブジェクト取得
		obj, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(s3rec.Bucket.Name),
			Key:    aws.String(s3rec.Object.Key),
		})
		if err != nil {
			log.Fatal(err)
		}

		// オブジェクト読み込み
		rc := obj.Body
		defer rc.Close()
		content, err := ioutil.ReadAll(rc)
		if err != nil {
			log.Fatal(err)
		}

		// オブジェクト書き込み
		_, err = svc.PutObject(&s3.PutObjectInput{
			Body:   bytes.NewReader(content),
			Bucket: aws.String(S3Output), // バケットは適宜変更
			Key:    aws.String(s3rec.Object.Key),
		})
		if err != nil {
			log.Fatal(err)
		}
	}

}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
