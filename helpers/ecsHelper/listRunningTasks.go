package ecshelper

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	aws_conf "github.com/harsh082ip/Video-transcoder_Go/aws"
)

func ListRunningTask() (int, error) {

	ecsClient, err := aws_conf.GetECSClient()
	if err != nil {
		return 0, err
	}

	// List all tasks in the ECS cluster
	listTasksResp, err := ecsClient.ListTasks(context.TODO(), &ecs.ListTasksInput{
		Cluster: aws.String("Video-Transcoder"),
	})
	if err != nil {
		return 0, fmt.Errorf("error in listing task, %v", err.Error())
	}

	// Check if there are tasks to describe
	if len(listTasksResp.TaskArns) == 0 {
		return 0, nil
	}

	// Describe ECS tasks
	describeTasksResp, err := ecsClient.DescribeTasks(context.TODO(), &ecs.DescribeTasksInput{
		Cluster: aws.String("Video-Transcoder"),
		Tasks:   listTasksResp.TaskArns,
	})
	if err != nil {
		return 0, fmt.Errorf("error in describing tasks, %v", err.Error())
	}

	// Count RUNNING tasks
	runningTaskCount := 0
	for _, task := range describeTasksResp.Tasks {
		if *aws.String(*task.LastStatus) == "RUNNING" {
			runningTaskCount++
		}
	}

	return runningTaskCount, nil
}
