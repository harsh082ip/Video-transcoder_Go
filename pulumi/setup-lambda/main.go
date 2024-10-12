package main

import (
	"log"
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Define the temperory bucket
		tempBucketName := os.Getenv("TEMP_BUCKET_NAME")
		if tempBucketName == "" {
			log.Println("TEMP_BUCKET_NAME cannot be empty")
			os.Exit(1)
		}

		// Get the existing s3 bucket details
		bucket, err := s3.LookupBucket(ctx, &s3.LookupBucketArgs{
			Bucket: tempBucketName,
		}, nil)
		if err != nil {
			return err
		}

		// Create an IAM role for the lambda function
		role, err := iam.NewRole(ctx, "lambdaRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`
	{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {
				"Service": "lambda.amazonaws.com"
			},
			"Action": "sts:AssumeRole"
		}]
	}`),
		})
		if err != nil {
			return err
		}

		// Attach AWSLambdaBasicExecutionRole policy to the role
		_, err = iam.NewRolePolicyAttachment(ctx, "lambdaRoleAttachment", &iam.RolePolicyAttachmentArgs{
			Role:      role.Name,
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
		})
		if err != nil {
			return err
		}

		// Create the Lambda function
		function, err := lambda.NewFunction(ctx, "goTestLambda", &lambda.FunctionArgs{
			Code:    pulumi.NewFileArchive("../../build/TriggerS3Upload/deployment.zip"),
			Runtime: pulumi.String("provided.al2"), // Amazon Linux 2023 uses provided.al2 runtime
			Architectures: pulumi.StringArray{
				pulumi.String("arm64"),
			},
			Handler: pulumi.String("hello.handler"),
			Role:    role.Arn,
			Name:    pulumi.String("go-test"),
		})
		if err != nil {
			return err
		}

		// Grant invoke permissions for S3 to invoke the Lambda function
		_, err = lambda.NewPermission(ctx, "s3ToLambdaPermission", &lambda.PermissionArgs{
			Action:    pulumi.String("lambda:InvokeFunction"),
			Function:  function.Name,
			Principal: pulumi.String("s3.amazonaws.com"),
			SourceArn: pulumi.String(bucket.Arn),
		})
		if err != nil {
			return err
		}

		// Configure bucket notification to trigger Lambda function on object create events
		_, err = s3.NewBucketNotification(ctx, "bucketNotification", &s3.BucketNotificationArgs{
			Bucket: pulumi.String(tempBucketName),
			LambdaFunctions: s3.BucketNotificationLambdaFunctionArray{
				&s3.BucketNotificationLambdaFunctionArgs{
					Events: pulumi.StringArray{
						pulumi.String("s3:ObjectCreated:*"),
					},
					LambdaFunctionArn: function.Arn,
				},
			},
		})
		if err != nil {
			return err
		}

		// Export the function URL (assuming no qualifier is used):
		ctx.Export("functionArn", pulumi.Sprintf("%s", function.Arn))

		return nil
	})
}
