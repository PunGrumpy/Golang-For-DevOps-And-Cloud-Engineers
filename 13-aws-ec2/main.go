package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {
	var (
		instanceId string
		err        error
	)
	ctx := context.Background()
	if instanceId, err = createEC2(ctx, "ap-southeast-1"); err != nil {
		fmt.Printf("Create EC2 failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created EC2 instance: %s\n", instanceId)
}

func createEC2(ctx context.Context, region string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %s\n", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)

	keypairs, err := ec2Client.DescribeKeyPairs(ctx, &ec2.DescribeKeyPairsInput{
		// KeyNames: []string{"go-aws-demo"},

		// If used Filters: ... instead of KeyNames: ... wouldn't get an error so might opt for that approach (but not all describe endpoints have this filter option)
		Filters: []types.Filter{
			{
				Name:   aws.String("key-name"),
				Values: []string{"go-aws-demo"},
			},
		},
	})
	// If used KeyNames
	// if err != nil && !strings.Contains(err.Error(), "InvalidKeyPair.NotFound") {
	// If used Filters
	if err != nil {
		return "", fmt.Errorf("Describe key pairs failed: %s\n", err)
	}

	// If used KeyNames
	// if keypairs == nil || len(keypairs.KeyPairs) == 0 {
	// If used Filters
	if len(keypairs.KeyPairs) == 0 {
		keypair, err := ec2Client.CreateKeyPair(ctx, &ec2.CreateKeyPairInput{
			KeyName: aws.String("go-aws-demo"),
		})
		if err != nil {
			return "", fmt.Errorf("Create key pair failed: %s\n", err)
		}
		err = os.WriteFile("go-aws-ec2.pem", []byte(*keypair.KeyMaterial), 0600) // 0600: read/write by owner only
		if err != nil {
			return "", fmt.Errorf("Write key pair to file failed: %s\n", err)
		}
	}

	imageOutput, err := ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("name"),
				Values: []string{"ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"},
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: []string{"hvm"},
			},
		},
		Owners: []string{"099720109477"},
	})
	if err != nil {
		return "", fmt.Errorf("Describe images failed: %s\n", err)
	}

	if len(imageOutput.Images) == 0 {
		return "", fmt.Errorf("ImageOutput.Images is empty\n")
	}

	instance, err := ec2Client.RunInstances(ctx, &ec2.RunInstancesInput{
		ImageId:      imageOutput.Images[0].ImageId,
		KeyName:      aws.String("go-aws-demo"),
		InstanceType: types.InstanceTypeT3Micro,
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
	})
	if err != nil {
		return "", fmt.Errorf("Run instances failed: %s\n", err)
	}

	if len(instance.Instances) == 0 {
		return "", fmt.Errorf("Instance.Instances is empty\n")
	}

	return *instance.Instances[0].InstanceId, nil
}
