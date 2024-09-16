package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

func NewAWSS3Storage() (*s3.Client, error) {
	conf, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			viper.GetString("s3.access_key_id"),
			viper.GetString("s3.secret_access_key"),
			viper.GetString("s3.session_token"),
		)),
	)

	s3client := s3.NewFromConfig(conf)
	return s3client, err
}
