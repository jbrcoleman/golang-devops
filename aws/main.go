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
	if instanceId, err = createEC2(ctx, "us-east-1"); err != nil {
		fmt.Printf("createEC2 error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("Instance id: %s\n", instanceId)
}

func createEC2(ctx context.Context, region string) (string, error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile("test-go"),
		config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("Unable to load SDK config, %s", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)

	keypairs, err := ec2Client.DescribeKeyPairs(ctx, &ec2.DescribeKeyPairsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("key-name"),
				Values: []string{"go-aws-demo"},
			},
		},
		//KeyNames: []string{"go-aws-demo"},
	})
	if err != nil {
		return "", fmt.Errorf("DescribeKeyPair error: %s", err)
	}

	if len(keypairs.KeyPairs) == 0 {
		keyPair, err := ec2Client.CreateKeyPair(ctx, &ec2.CreateKeyPairInput{

			KeyName: aws.String("go-aws-demo"),
		})
		if err != nil {
			return "", fmt.Errorf("CreateKeyPair error: %s", err)
		}
		//keyPair.KeyMaterial
		err = os.WriteFile("go-aws-ec2.pem", []byte(*keyPair.KeyMaterial), 0600)
		if err != nil {
			return "", fmt.Errorf("WriteFile error: %s", err)
		}
	}

	imageOutput, err := ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("virtualization-type"),
				Values: []string{"hvm"},
			},
			{},
		},
		Owners: []string{"099720109477"},
	})
	if err != nil {
		return "", fmt.Errorf("DescribeImages error: %s", err)
	}
	if len(imageOutput.Images) == 0 {
		return "", fmt.Errorf("imageOutput.Images is of 0 length")
	}
	runInstance, err := ec2Client.RunInstances(ctx, &ec2.RunInstancesInput{
		ImageId:      imageOutput.Images[0].ImageId,
		KeyName:      aws.String("go-aws-demo"),
		InstanceType: types.InstanceTypeT3Micro,
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
	})
	//imageOutput.Images[0].ImageId

	if len(runInstance.Instances) == 0 {
		return "", fmt.Errorf("instance.Instances is of 0 length")
	}

	return *runInstance.Instances[0].InstanceId, nil
}
