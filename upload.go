package awsm

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/roidaradal/fn/lang"
)

type UploadConfig struct {
	Profile    string
	Region     string
	Bucket     string
	FilePath   string
	BucketPath string
	ACL        types.ObjectCannedACL
}

// Upload file to S3 bucket
func UploadFile(cfg *UploadConfig) error {
	// Load AWS configuration
	profile := lang.Ternary(cfg.Profile == "", "default", cfg.Profile)
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(awsCfg)

	// Open local file
	file, err := os.Open(cfg.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", cfg.FilePath, err)
	}
	defer file.Close()

	// Create s3 uploader
	uploader := manager.NewUploader(client)

	// Upload file
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &cfg.Bucket,
		Key:    &cfg.BucketPath,
		Body:   file,
		ACL:    cfg.ACL,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file %q to S3: %w", cfg.FilePath, err)
	}

	return nil
}
