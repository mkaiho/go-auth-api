package infrastructure

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kelseyhightower/envconfig"
	"github.com/mkaiho/go-auth-api/adapter/storage"
	"github.com/mkaiho/go-auth-api/util"
)

var _ storage.Client = (*S3Client)(nil)

type S3Config struct {
	JWKBucket string `envconfig:"JWK_BUCKET" required:"true"`
}

func LoadS3Config() (*S3Config, error) {
	var c S3Config
	if err := envconfig.Process("S3", &c); err != nil {
		return nil, err
	}
	return &c, nil
}

type S3Client struct {
	bucket string
	client *s3.Client
}

func NewS3Client(bucket string, conf aws.Config) *S3Client {
	return &S3Client{
		bucket: bucket,
		client: s3.NewFromConfig(conf),
	}
}

func (c S3Client) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	output, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    &path,
	})
	if err != nil {
		return nil, err
	}

	return output.Body, nil
}

func (c S3Client) Save(ctx context.Context, path string, mime storage.MimeType, body []byte) error {
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &c.bucket,
		Key:         &path,
		ContentType: util.ToPointer(mime.String()),
		Body:        bytes.NewReader(body),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c S3Client) Remove(ctx context.Context, path string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &c.bucket,
		Key:    &path,
	})
	if err != nil {
		return err
	}

	return nil
}
