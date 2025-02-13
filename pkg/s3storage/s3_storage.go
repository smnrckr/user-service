package s3storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	AccessKey       string
	SecretAccessKey string
	Region          string
}

type S3Service struct {
	Client *s3.Client
}

func NewS3Service(cfg *S3Config) (*S3Service, error) {
	credentials := credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretAccessKey, "")
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsConfig)
	return &S3Service{Client: client}, nil
}

func (s *S3Service) UploadFile(BucketName string, Key string, Body []byte) (string, error) {
	uploader := manager.NewUploader(s.Client)

	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(Key),
		Body:   bytes.NewReader(Body),
	})
	if err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", BucketName, Key)
	return fileURL, nil
}
