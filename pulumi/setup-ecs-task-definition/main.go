package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecs"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Retrieve the image URI from environment variables
		imageURI := os.Getenv("IMAGE_URI")
		if imageURI == "" {
			log.Println("IMAGE_URI cannot be empty")
			os.Exit(1)
		}

		// Lookup ECS Task Execution Role
		executionRole, err := iam.LookupRole(ctx, &iam.LookupRoleArgs{
			Name: "ecsTaskExecutionRole",
		}, nil)
		if err != nil {
			return fmt.Errorf("failed to get execution role: %w", err)
		}

		// Lookup ECS Task Role (can be same or different from execution role)
		taskRole, err := iam.LookupRole(ctx, &iam.LookupRoleArgs{
			Name: "ecsTaskExecutionRole",
		}, nil)
		if err != nil {
			return fmt.Errorf("failed to get task role: %w", err)
		}

		// Define container configurations
		containerDef := fmt.Sprintf(`[
			{
				"name": "ffmpeg-container",
				"image": "%s",
				"cpu": 1024,
				"memory": 3072,
				"memoryReservation": 1024,
				"essential": true,
				"portMappings": [
					{
						"containerPort": 80,
						"hostPort": 80,
						"protocol": "tcp",
						"appProtocol": "http"
					}
				],
				"logConfiguration": {
					"logDriver": "awslogs",
					"options": {
						"awslogs-group": "/ecs/go-task-v1",
						"awslogs-create-group": "true",
						"awslogs-region": "ap-south-1",
						"awslogs-stream-prefix": "ecs"
					}
				}
			}
		]`, imageURI)

		// Create ECS Task Definition
		taskDefinition, err := ecs.NewTaskDefinition(ctx, "ecsTaskDefinition", &ecs.TaskDefinitionArgs{
			Family:                  pulumi.String("go-task-v1"),
			Cpu:                     pulumi.String("1024"), // 1 vCPU
			Memory:                  pulumi.String("3072"), // 3 GB
			NetworkMode:             pulumi.String("awsvpc"),
			RequiresCompatibilities: pulumi.StringArray{pulumi.String("FARGATE")},
			RuntimePlatform: &ecs.TaskDefinitionRuntimePlatformArgs{
				CpuArchitecture:       pulumi.String("X86_64"),
				OperatingSystemFamily: pulumi.String("LINUX"),
			},
			TaskRoleArn:          pulumi.String(taskRole.Arn),
			ExecutionRoleArn:     pulumi.String(executionRole.Arn),
			ContainerDefinitions: pulumi.String(containerDef),
			EphemeralStorage: &ecs.TaskDefinitionEphemeralStorageArgs{
				SizeInGib: pulumi.Int(21), // Configurable ephemeral storage
			},
		})
		if err != nil {
			return fmt.Errorf("error creating ECS task definition: %w", err)
		}

		// Export the ARN of the task definition
		ctx.Export("taskDefinitionArn", taskDefinition.Arn)
		return nil
	})
}
