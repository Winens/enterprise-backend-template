package service

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

func NewS3Storage() (*minio.Client, error) {
	opts := &minio.Options{
		Creds: credentials.NewStaticV4(viper.GetString("s3.access_key_id"), viper.GetString("s3.secret_access_key"),
			viper.GetString("s3.session_token")),
		Secure: viper.GetBool("s3.use_ssl"),
	}

	return minio.New(viper.GetString("s3.endpoint"), opts)
}

func MigrateS3Buckets(s3 *minio.Client) error {
	ctx := context.Background()
	if err := s3.MakeBucket(ctx, "attachments", minio.MakeBucketOptions{}); err != nil {
		return err
	}

	return nil
}
