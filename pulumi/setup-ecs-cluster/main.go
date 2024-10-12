package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an ECS cluster
		cluster, err := ecs.NewCluster(ctx, "goCluster", &ecs.ClusterArgs{
			Name: pulumi.String("go-Cluster"),
		})
		if err != nil {
			return err
		}

		// Export the cluster ARN
		ctx.Export("clusterArn", cluster.Arn)

		return nil
	})
}
