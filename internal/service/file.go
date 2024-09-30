package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type FileService struct {
	Client           *s3.Client
	BucketName       string
	Region           string
	EndpointResolver string
}

func NewFileService(
	ctx context.Context, backetName, region, endpointResolverURL string,
) *FileService {
	client, err := initS3Client(ctx, region, endpointResolverURL)
	if err != nil {
		return nil
	}

	service := &FileService{
		Client:           client,
		BucketName:       backetName,
		Region:           region,
		EndpointResolver: endpointResolverURL,
	}
	return service
}

func initS3Client(ctx context.Context, regionInput, url string) (*s3.Client, error) {
	// Создаем кастомный обработчик эндпоинтов, который для сервиса S3 и региона ru-central1 выдаст корректный URL
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID && region == regionInput {
				return aws.Endpoint{
					PartitionID:   "yc",
					URL:           url,
					SigningRegion: regionInput,
				}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		},
	)

	// Подгружаем конфигрурацию из ~/.aws/*
	cfg, err := config.LoadDefaultConfig(ctx, config.WithEndpointResolverWithOptions(customResolver))

	if err != nil {
		log.Error("FileService - initS3Client - ", err)
	}

	// Создаем клиента для доступа к хранилищу S3
	client := s3.NewFromConfig(cfg)
	return client, nil
}

func (f *FileService) Upload(ctx context.Context, file graphql.Upload) (string, error) {
	result, err := f.Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	for _, bucket := range result.Buckets {
		log.Printf("Upload bucket=%s creation time=%s", aws.ToString(bucket.Name), bucket.CreationDate.Format("2006-01-02 15:04:05 Monday"))
	}

	name := uuid.New().String()
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	pathName := strings.Join([]string{name, ext}, "")
	_, err = f.Client.PutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Bucket: aws.String(f.BucketName), Key: aws.String(pathName), Body: file.File,
		},
	)
	if err != nil {
		return "", err
	}

	return pathName, nil
}

func (f *FileService) Delete(ctx context.Context, path string) (bool, error) {
	_, err := f.Client.DeleteObject(
		ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(f.BucketName),
			Key:    aws.String(path),
		},
	)

	if err != nil {
		log.Error("FileService - Delete - f.Client.DeleteObject - ", err)
		return false, err
	}
	return true, nil
}

func (f *FileService) BuildImageURL(pathName string) string {
	elems := []string{f.EndpointResolver, f.BucketName, pathName}
	path := strings.Join(elems, "/")
	return path
}
