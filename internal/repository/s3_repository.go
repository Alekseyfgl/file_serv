package repository

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Repository struct {
	Client     *s3.Client
	Uploader   *manager.Uploader
	BucketName string
}

func NewS3Repository(bucketName, accessKey, secretKey string) *S3Repository {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("ru-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if service == s3.ServiceID {
					return aws.Endpoint{
						URL:           "https://s3.timeweb.cloud",
						SigningRegion: "ru-1",
					}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			}),
		),
	)
	if err != nil {
		log.Fatalf("Ошибка загрузки AWS конфигурации: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	return &S3Repository{
		Client:     client,
		Uploader:   uploader,
		BucketName: bucketName,
	}
}

func (r *S3Repository) UploadFile(ctx context.Context, key, contentType string, body io.Reader) (string, error) {
	_, err := r.Uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.BucketName),
		Key:         aws.String(key),
		ACL:         "public-read",
		ContentType: aws.String(contentType),
		Body:        body,
	})
	if err != nil {
		return "", err
	}
	finalURL := fmt.Sprintf("https://%s.s3.timeweb.cloud/%s", r.BucketName, key)
	return finalURL, nil
}

func (r *S3Repository) ListFilesByPrefix(ctx context.Context, prefix string) ([]types.Object, error) {
	listResp, err := r.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(r.BucketName),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}
	return listResp.Contents, nil
}

func (r *S3Repository) DeleteFile(ctx context.Context, key string) error {
	_, err := r.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(key),
	})
	return err
}

// Удаление группы файлов (по одному в цикле).
func (r *S3Repository) DeleteFilesBatch(ctx context.Context, keys []string) error {
	for _, k := range keys {
		if err := r.DeleteFile(ctx, k); err != nil {
			return err
		}
	}
	return nil
}

// FolderExists проверяет, существует ли указанный префикс (папка) в S3.
// Например, folderName = "photos/".
func (r *S3Repository) FolderExists(ctx context.Context, folderName string) (bool, error) {
	// Можно сделать ListObjects или HeadObject с таким префиксом и проверить результат.
	// Простейший подход — ListObjects с префиксом "photos/" и MaxKeys=1.

	resp, err := r.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(r.BucketName),
		Prefix: aws.String(folderName),
		//MaxKeys: 1,
	})
	if err != nil {
		return false, err
	}

	// Если в ответе нет ни одного объекта, считаем что “папки” нет.
	return len(resp.Contents) > 0, nil
}
