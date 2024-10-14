package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get environment variables
		keyPairName := os.Getenv("KEY_PAIR_NAME")
		if keyPairName == "" {
			log.Fatal("KEY_PAIR_NAME cannot be empty")
			os.Exit(1)
		}

		// Fetch the latest Ubuntu 24.04 LTS AMI ID
		ami, err := ec2.LookupAmi(ctx, &ec2.LookupAmiArgs{
			Filters: []ec2.GetAmiFilter{
				{
					Name:   "name",
					Values: []string{"ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"},
				},
				{
					Name:   "state",
					Values: []string{"available"},
				},
			},
			Owners:     []string{"099720109477"}, // Canonical's AWS ID
			MostRecent: pulumi.BoolRef(true),
		})
		if err != nil {
			log.Fatalf("Error fetching AMI: %v", err)
		}

		// Lookup the existing key pair
		keyPair, err := ec2.LookupKeyPair(ctx, &ec2.LookupKeyPairArgs{
			KeyName: &keyPairName,
		}, nil)
		if err != nil {
			return fmt.Errorf("failed to find key pair: %w", err)
		}

		// Create a security group to allow SSH traffic
		secGroup, err := ec2.NewSecurityGroup(ctx, "videoTranscoderWorkerSecGroup", &ec2.SecurityGroupArgs{
			Ingress: ec2.SecurityGroupIngressArray{
				&ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(22),
					ToPort:     pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				&ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"), // Allow all traffic
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		})
		if err != nil {
			return err
		}

		// Define user data script
		userData := `#!/bin/bash
		sudo apt update -y && sudo apt upgrade -y
		sudo apt install git -y
		git clone https://github.com/harsh082ip/Video-transcoder_Go`

		// Create the EC2 instance
		instance, err := ec2.NewInstance(ctx, "videoTranscoderWorker", &ec2.InstanceArgs{
			InstanceType:   pulumi.String("t2.micro"),
			Ami:            pulumi.String(ami.Id), // Use ami.Id directly
			KeyName:        pulumi.String(*keyPair.KeyName),
			SecurityGroups: pulumi.StringArray{secGroup.Name},
			UserData:       pulumi.String(userData),
			RootBlockDevice: &ec2.InstanceRootBlockDeviceArgs{
				VolumeSize: pulumi.Int(8),
				VolumeType: pulumi.String("gp3"),
			},
			Tags: pulumi.StringMap{
				"Name": pulumi.String("video-transcoder-worker"),
			},
		})
		if err != nil {
			return err
		}

		// Export useful instance details
		ctx.Export("instancePublicIP", instance.PublicIp)
		ctx.Export("instancePublicDns", instance.PublicDns)

		return nil
	})
}
