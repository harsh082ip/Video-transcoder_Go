package aws_conf

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// GetAwsConf loads the AWS configuration
// It returns an aws.Config object and an error if any occurred

func GetAwsConf() (aws.Config, error) {
	// Hardcoded credentials (not recommended for prod)
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")

	// Load the AWS configuration with custom credentials and region
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-south-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)

	if err != nil {
		return aws.Config{}, fmt.Errorf("unable to load SDK config: %v", err)
	}

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

// ecs client
func GetECSClient() (*ecs.Client, error) {
	// Get the AWS configuration
	cfg, err := GetAwsConf()
	if err != nil {
		// Return an empty ecs client and a formatted error message
		return nil, fmt.Errorf("error in getting ecs client: %v", err.Error())
	}
	// Create an ecs client from the configuration
	ecsClient := ecs.NewFromConfig(cfg)
	// Return the ecs client and no error
	return ecsClient, nil
}
