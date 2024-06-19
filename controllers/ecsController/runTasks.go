package ecscontroller

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	aws_conf "github.com/harsh082ip/Video-transcoder_Go/aws"
)

func RunECSTask(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, SOURCE_IMAGE, DESTINATION_1080, DESTINATION_720, DESTINATION_360 string) (bool, error) {

	ecsClient, err := aws_conf.GetECSClient()
	if err != nil {
		return false, fmt.Errorf("error in getting a ecs client, %v", err.Error())
	}

	// Define environment variables
	environment := []ecsTypes.KeyValuePair{
		{
			Name:  aws.String("AWS_ACCESS_KEY_ID"),
			Value: aws.String(AWS_ACCESS_KEY_ID),
		},
		{
			Name:  aws.String("AWS_SECRET_ACCESS_KEY"),
			Value: aws.String(AWS_SECRET_ACCESS_KEY),
		},
		{
			Name:  aws.String("AWS_DEFAULT_REGION"),
			Value: aws.String("ap-south-1"),
		},
		{
			Name:  aws.String("SOURCE_IMAGE"),
			Value: aws.String(SOURCE_IMAGE),
		},
		{
			Name:  aws.String("DESTINATION_1080"),
			Value: aws.String(DESTINATION_1080),
		},
		{
			Name:  aws.String("DESTINATION_720"),
			Value: aws.String(DESTINATION_720),
		},
		{
			Name:  aws.String("DESTINATION_360"),
			Value: aws.String(DESTINATION_360),
		},
	}

	// Define container overrides
	containerOverride := ecsTypes.ContainerOverride{
		Name:        aws.String("ffmpeg-container"),
		Environment: environment,
	}

	// Define task overrides
	taskOverride := ecsTypes.TaskOverride{
		ContainerOverrides: []ecsTypes.ContainerOverride{containerOverride},
	}

	// Run task
	runTaskInput := &ecs.RunTaskInput{
		Cluster:        aws.String("go-Cluster"),   // Replace with your cluster name
		TaskDefinition: aws.String("go-task-v1"),   // Replace with your task definition
		LaunchType:     ecsTypes.LaunchTypeFargate, // Or ecsTypes.LaunchTypeEc2 if using EC2 launch type
		NetworkConfiguration: &ecsTypes.NetworkConfiguration{ // Required for FARGATE
			AwsvpcConfiguration: &ecsTypes.AwsVpcConfiguration{
				Subnets:        []string{"subnet-0dd67375bfa2bc37d", "subnet-0bd233120f2dc36b0", "subnet-01234dc71b3c22354"}, // Replace with your subnet
				SecurityGroups: []string{"sg-03ff3853d671a1697"},                                                             // Replace with your security group
				AssignPublicIp: ecsTypes.AssignPublicIpEnabled,
			},
		},
		Overrides: &taskOverride,
	}

	result, err := ecsClient.RunTask(context.TODO(), runTaskInput)
	if err != nil {
		return false, fmt.Errorf("failed to run task, %v", err)
	}

	// Print result
	for _, task := range result.Tasks {
		fmt.Printf("Started task: %s\n", *task.TaskArn)
		return true, nil
	}

	return false, fmt.Errorf("unexpected Error in Starting a task")
}

// ffmpeg-container
