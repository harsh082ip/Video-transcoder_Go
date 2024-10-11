package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
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

		// Enable function URL
		_, err = lambda.NewFunctionUrl(ctx, "goTestLambdaUrl", &lambda.FunctionUrlArgs{
			FunctionName:      function.Name,
			AuthorizationType: pulumi.String("NONE"),
		})
		if err != nil {
			return err
		}

		// Export the function URL (assuming no qualifier is used):
		ctx.Export("functionUrl", pulumi.Sprintf("https://%s.lambda-url.<region>.on.aws", function.Arn))

		return nil
	})
}
