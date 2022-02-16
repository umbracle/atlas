package aws

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/umbracle/atlas/internal/proto"

	"fmt"
)

type config struct {
	Type string `json:"type"`
}

type AwsProvider struct {
}

func (a *AwsProvider) Config() (interface{}, error) {
	panic("X")
}

var userDataRaw = `#!/bin/bash

yum update -y
amazon-linux-extras install docker
service docker start
systemctl enable docker
usermod -a -G docker ec2-user
docker info

# install yum
yum install -y tmux

# download atlas
echo "atlas"
curl -o /usr/bin/atlas https://4e88-88-9-192-173.ngrok.io/atlas && chmod +x /usr/bin/atlas
echo "atlas done"

# start the agent session
tmux new-session -d -s atlas '/usr/bin/atlas agent'
`

// t2.small, t2.medium

func (a *AwsProvider) readInstanceById(instanceId string) (*ec2.Instance, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		return nil, err
	}

	// Create EC2 service client
	svc := ec2.New(sess)

	output, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: []*string{aws.String(instanceId)}})
	if err != nil {
		return nil, err
	}
	return output.Reservations[0].Instances[0], nil
}

func (a *AwsProvider) Update(ctx context.Context, node *proto.Node) error {

	var config *config
	if err := json.Unmarshal([]byte(node.ProviderConfig), &config); err != nil {
		return err
	}

	// config type
	if config.Type == "" {
		config.Type = "t2.medium"
	}

	fmt.Println("-- aws config --")
	fmt.Println(config.Type)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		return err
	}

	// Create EC2 service client
	svc := ec2.New(sess)

	if node.Handle != nil {
		// this is an update
		var xxx handle
		if err := json.Unmarshal([]byte(node.Handle.Handle), &xxx); err != nil {
			panic(err)
		}

		log.Printf("[INFO] Stop instance")

		_, err := svc.StopInstances(&ec2.StopInstancesInput{
			InstanceIds: []*string{aws.String(xxx.InstanceID)},
		})
		if err != nil {
			panic(err)
		}

		time.Sleep(10 * time.Second)

		_, err = svc.ModifyInstanceAttribute(&ec2.ModifyInstanceAttributeInput{
			InstanceId: aws.String(xxx.InstanceID),
			InstanceType: &ec2.AttributeValue{
				Value: aws.String(config.Type),
			},
		})
		if err != nil {
			panic(err)
		}

		log.Printf("[INFO] Starting Instance %q after instance_type change", xxx.InstanceID)

		time.Sleep(10 * time.Second)

		input := &ec2.StartInstancesInput{
			InstanceIds: []*string{aws.String(xxx.InstanceID)},
		}
		if _, err := svc.StartInstances(input); err != nil {
			panic(err)
		}
		return nil
	}

	userData := base64.StdEncoding.EncodeToString([]byte(userDataRaw))

	fmt.Println(userData)

	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		ImageId:      aws.String("ami-0341aeea105412b57"),
		InstanceType: aws.String(config.Type),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sdh"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize: aws.Int64(1024),
				},
			},
		},
		KeyName:  aws.String("atlas"),
		UserData: aws.String(userData),
	})
	if err != nil {
		return err
	}

	instance := runResult.Instances[0]
	fmt.Println("Created instance", *instance.InstanceId)

	instanceId := instance.InstanceId

	time.Sleep(2 * time.Second)

	// loop until we get the ip address
	if err := svc.WaitUntilInstanceExists(&ec2.DescribeInstancesInput{InstanceIds: []*string{instanceId}}); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	output, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: []*string{instanceId}})
	if err != nil {
		return err
	}

	fmt.Println("-- otuput --")
	fmt.Println(output)

	ipAddress := *output.Reservations[0].Instances[0].PublicIpAddress
	fmt.Println(ipAddress)

	awsHandle := &handle{
		InstanceID: *instance.InstanceId,
	}
	handleRaw, err := json.Marshal(awsHandle)
	if err != nil {
		return err
	}

	handle := &proto.Node_Handle{
		Handle: string(handleRaw),
		Ip:     ipAddress,
	}

	// Add tags to the created instance
	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{instance.InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("MyFirstInstance"),
			},
		},
	})
	if errtag != nil {
		return errtag
	}

	fmt.Println("Successfully tagged instance")
	node.Handle = handle
	return nil
}

type handle struct {
	InstanceID string
}
