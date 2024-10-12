package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		bucketName := os.Getenv("TEMP_BUCKET_NAME")
		if bucketName == "" {
			log.Println("TEMP_BUCKET_NAME cannot be empty")
			os.Exit(1)
		}
		log.Println("bucketname:", bucketName)

		// Create an S3 bucket with the specified name
		bucket, err := s3.NewBucket(ctx, bucketName, &s3.BucketArgs{
			Acl: pulumi.String("private"),
			ServerSideEncryptionConfiguration: &s3.BucketServerSideEncryptionConfigurationArgs{
				Rule: &s3.BucketServerSideEncryptionConfigurationRuleArgs{
					ApplyServerSideEncryptionByDefault: &s3.BucketServerSideEncryptionConfigurationRuleApplyServerSideEncryptionByDefaultArgs{
						SseAlgorithm: pulumi.String("AES256"), // Amazon S3 managed keys (SSE-S3)
					},
					BucketKeyEnabled: pulumi.Bool(true),
				},
			},
			Bucket: pulumi.String(bucketName),
		})
		if err != nil {
			return err
		}

		// Disable block all public access
		publicAccessBlock, err := s3.NewBucketPublicAccessBlock(ctx, bucketName+"PublicAccessBlock", &s3.BucketPublicAccessBlockArgs{
			Bucket:                bucket.ID(),
			BlockPublicAcls:       pulumi.Bool(false),
			IgnorePublicAcls:      pulumi.Bool(false),
			BlockPublicPolicy:     pulumi.Bool(false),
			RestrictPublicBuckets: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// Use bucket.ID().ApplyT to dynamically access the bucket name
		bucket.ID().ApplyT(func(bucketID string) (string, error) {
			fmt.Println("bucketID: ", bucketID)

			// Define the bucket policy
			bucketPolicy := fmt.Sprintf(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "Stmt14055921390d0",
						"Effect": "Allow",
						"Principal": "*",
						"Action": "s3:*",
						"Resource": [
							"arn:aws:s3:::%s",
							"arn:aws:s3:::%s/*"
						]
					}
				]
			}`, bucketID, bucketID)

			/*
					I added a dependency between the BucketPolicy and the BucketPublicAccessBlock.
					This ensures that the public access block is fully disabled before attempting
					to apply the bucket policy. This is important because public policies are
				 	blocked if the BlockPublicPolicy setting is enabled.
			*/

			// Attach the bucket policy with a dependency on the public access block setting
			_, err = s3.NewBucketPolicy(ctx, bucketName+"Policy", &s3.BucketPolicyArgs{
				Bucket: pulumi.String(bucketID), // Use the bucketID string
				Policy: pulumi.String(bucketPolicy),
			}, pulumi.DependsOn([]pulumi.Resource{publicAccessBlock}))
			if err != nil {
				return "", err
			}
			return bucketID, nil
		})

		// Export the bucket name
		ctx.Export("bucketName", bucket.Bucket)
		return nil
	})
}
