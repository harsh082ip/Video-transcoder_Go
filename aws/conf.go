package aws_conf

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// GetAwsConf loads the AWS configuration
// It returns an aws.Config object and an error if any occurred
func GetAwsConf() (aws.Config, error) {
	// Load the default AWS configuration with the specified region
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		// Return an empty config and a formatted error message
		return aws.Config{}, fmt.Errorf("unable to load sdk config: %v", err.Error())
	}
	// Return the loaded configuration and no error
	return cfg, nil
}

// GetS3Client creates an S3 client using the loaded AWS configuration
// It returns a pointer to an s3.Client object and an error if any occurred
func GetS3Client() (*s3.Client, error) {
	// Get the AWS configuration
	cfg, err := GetAwsConf()
	if err != nil {
		// Return an empty S3 client and a formatted error message
		return nil, fmt.Errorf("error in getting S3 client: %v", err.Error())
	}
	// Create an S3 client from the configuration
	s3Client := s3.NewFromConfig(cfg)
	// Return the S3 client and no error
	return s3Client, nil
}
