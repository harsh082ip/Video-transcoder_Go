package main

import (
	"log"
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {

	ecrRepositoryName := os.Getenv("ECR_REPOSITORY_NAME")
	if ecrRepositoryName == "" {
		log.Println("ECR_REPOSITORY_NAME cannot be empty")
		os.Exit(1)
	}

	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := ecr.NewRepository(ctx, ecrRepositoryName+"Repo", &ecr.RepositoryArgs{
			Name:               pulumi.String(ecrRepositoryName),
			ImageTagMutability: pulumi.String("MUTABLE"),
			EncryptionConfigurations: ecr.RepositoryEncryptionConfigurationArray{
				&ecr.RepositoryEncryptionConfigurationArgs{
					EncryptionType: pulumi.String("AES256"),
				},
			},
		})
		if err != nil {
			return err
		}
		return nil
	})
}

// error: deleting urn:pulumi:dev::setup-s3-temp::aws:s3/bucket:Bucket::test.harsh54323: 1 error occurred:
// * error deleting S3 Bucket (test.harsh54323): BucketNotEmpty: The bucket you tried to delete is not empty
// status code: 409, request id: CC9GQF7VYV0T3SE3, host id: U90fFbvxN/i0UIXBlewHvNTAdi54oNPTJdjw9YQFc/R6l4u8ufzqKva0IDqXMgp324Cx6VGtzxOEXbeEaT2L7Swy4+rmVryO
