package services

import (
	"io"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// WriteToS3 using a reader, write to S3 bucket.
func WriteToS3(s3Svc *s3manager.Uploader, r io.Reader, s3Bucket, fileName, s3Key string) error {
	log.Printf("Starting to write response to file : %v to S3 bucket path: %v/%v", fileName, s3Bucket, s3Key)
	result, err := s3Svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Key + fileName),
		Body:   r,
	})

	if err != nil {
		log.Printf("failed to write response to file : %v to S3 bucket path: %v/%v", fileName, s3Bucket, s3Key)
		return err
	}
	log.Printf("filed uploaded to, %s\n", result.Location)
	return nil
}

// GenerateFileName generates a file path for S3.
func GenerateFileName(file, fileType string) string {
	if !strings.HasPrefix(fileType, ".") {
		fileType = "." + fileType
	}

	currentTime := time.Now()
	formattedTime := strings.ReplaceAll(currentTime.Format("2006-01-02-15:04:05"), ":", "")
	fileName := formattedTime + "-" + file + fileType
	return fileName
}
