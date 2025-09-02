package s3

import (
	"context"
	"fmt"
	"strings"

	cmnlogger "github.com/Justksenia/common/logger"
	"github.com/Justksenia/common/tracer"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/go-faster/errors"
	"go.uber.org/zap"
)

type FileStorage struct {
	client *s3.Client
	bucket string
}

func New(cfg *Config) (*FileStorage, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   cfg.PartitionID,
			URL:           cfg.Url,
			SigningRegion: cfg.Region,
		}, nil
	})

	credentialsProvider := aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     cfg.AccessKey,
			SecretAccessKey: cfg.SecretKey,
		}, nil
	})
	client := s3.NewFromConfig(aws.Config{
		Region:                      cfg.Region,
		EndpointResolverWithOptions: customResolver,
		Credentials:                 credentialsProvider,
	})
	return &FileStorage{client: client, bucket: cfg.Bucket}, nil
}

func (f *FileStorage) Upload(ctx context.Context, file *File) (string, error) {
	key := fmt.Sprintf("%s/%s", file.path, file.name)
	_, err := f.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(key),
		Body:   file.data,
	})
	if err != nil {
		return key, errors.Wrap(err, "failed to upload file")
	}
	return key, nil
}

func (f *FileStorage) GetAllPaths(ctx context.Context, pattern string) ([]string, error) {
	resp, err := f.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(f.bucket),
		Prefix: aws.String(pattern),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list objects")
	}
	paths := make([]string, 0, len(resp.Contents))
	for _, obj := range resp.Contents {
		paths = append(paths, *obj.Key)
	}
	return paths, nil
}

func (f *FileStorage) GetAll(ctx context.Context, path string) ([]File, error) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	logger := cmnlogger.FromContext(ctx).With(zap.String("Method", tracer.AutoFillName()))

	resp, err := f.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(f.bucket),
		Prefix: aws.String(path),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list objects")
	}

	files := make([]File, 0, len(resp.Contents))

	for _, obj := range resp.Contents {
		file, inErr := f.Get(ctx, *obj.Key)
		if inErr != nil {
			logger.Error("Failed to get file", zap.String("key", *obj.Key), zap.Error(inErr))
			continue
		}
		files = append(files, *file)
	}
	return files, nil
}

func (f *FileStorage) Get(ctx context.Context, path string) (*File, error) {
	resp, err := f.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get object")
	}
	fileName := path
	splitted := strings.Split(path, "/")
	if len(splitted) > 0 {
		fileName = splitted[len(splitted)-1]
	}

	return &File{
		name: fileName,
		path: path,
		data: resp.Body,
	}, nil
}

func (f *FileStorage) Delete(ctx context.Context, path string) error {
	_, err := f.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return errors.Wrap(err, "failed to delete object")
	}
	return nil
}

func (f *FileStorage) DeleteAll(ctx context.Context, path string) error {
	objects, err := f.GetAll(ctx, path)
	if err != nil {
		return errors.Wrap(err, "failed to get objects")
	}

	deleteObjectsId := make([]types.ObjectIdentifier, 0, len(objects))
	for _, obj := range objects {
		deleteObjectsId = append(deleteObjectsId, types.ObjectIdentifier{
			Key: &obj.path,
		})
	}

	_, err = f.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(f.bucket),
		Delete: &types.Delete{
			Objects: deleteObjectsId,
		},
	})

	if err != nil {
		return errors.Wrap(err, "failed to delete objects")
	}
	return nil
}
