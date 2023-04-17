package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const bucketName = "aws-demo-test-bucket-pungrumpy"
const regionName = "ap-southeast-1"

func main() {
	var (
		s3Client *s3.Client
		err      error
		out      []byte
	)

	ctx := context.Background()
	if s3Client, err = initS3Client(ctx, regionName); err != nil {
		fmt.Printf("unable to initialize S3 client: %s\n", err)
		os.Exit(1)
	}
	if err = createS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("unable to create S3 bucket: %s\n", err)
		os.Exit(1)
	}
	if err = uploadToS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("unable to upload to S3 bucket: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully uploaded to S3 bucket %s\n", bucketName)
	if out, err = downloadFromS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("unable to download from S3 bucket: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully downloaded from S3 bucket: %s\n", string(out))
}

func initS3Client(ctx context.Context, region string) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %s", err)
	}

	return s3.NewFromConfig(cfg), nil
}

func createS3Bucket(ctx context.Context, s3Client *s3.Client) error {
	allBuckets, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("unable to list S3 buckets: %s", err)
	}

	found := false
	for _, bucket := range allBuckets.Buckets {
		if *bucket.Name == bucketName {
			found = true
			fmt.Printf("S3 bucket %s already exists\n", bucketName)
		}
	}

	if !found {
		_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: regionName,
			},
		})
		if err != nil {
			return fmt.Errorf("unable to create S3 bucket: %s", err)
		}
	}

	return nil
}

func uploadToS3Bucket(ctx context.Context, s3Client *s3.Client) error {
	testFile, err := ioutil.ReadFile("test.txt")
	if err != nil {
		return fmt.Errorf("unable to read test file: %s", err)
	}
	uploader := manager.NewUploader(s3Client)
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("test.txt"),
		Body:   bytes.NewReader(testFile),
	})
	if err != nil {
		return fmt.Errorf("unable to upload to S3 bucket: %s", err)
	}

	return nil
}

func downloadFromS3Bucket(ctx context.Context, s3Client *s3.Client) ([]byte, error) {
	downloader := manager.NewDownloader(s3Client)
	buffer := manager.NewWriteAtBuffer([]byte{})

	numBytes, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("test.txt"),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to download from S3 bucket: %s", err)
	}

	if numBytesReceived := len(buffer.Bytes()); numBytes != int64(numBytesReceived) {
		return nil, fmt.Errorf("number of bytes downloaded (%d) does not match number of bytes received (%d)", numBytes, numBytesReceived)
	}

	return buffer.Bytes(), nil
}
