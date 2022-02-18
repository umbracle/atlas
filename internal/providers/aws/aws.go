package aws

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/umbracle/atlas/internal/proto"
	"github.com/umbracle/atlas/internal/schema"
	"github.com/umbracle/atlas/internal/userdata"

	"fmt"
)

type config struct {
	Type string `json:"type"`
}

type AwsProvider struct {
	conn *ec2.EC2
}

func (a *AwsProvider) Init() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		panic(err)
	}
	a.conn = ec2.New(sess)
}

func (a *AwsProvider) Config() (interface{}, error) {
	panic("X")
}

func (a *AwsProvider) Schema() *schema.Object {
	return &schema.Object{
		Fields: map[string]*schema.Field{
			"instance_type": {
				Type: &schema.String,
			},
		},
	}
}

// t2.small, t2.medium

func (a *AwsProvider) resourceInstanceFind(id string) (*ec2.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice([]string{id}),
	}
	resp, err := a.conn.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	if len(resp.Reservations) == 0 {
		return nil, nil
	}

	instances := resp.Reservations[0].Instances
	if len(instances) == 0 {
		return nil, nil
	}

	return instances[0], nil
}

func (a *AwsProvider) waitForState(ctx context.Context, id string, expectedState string) (*ec2.Instance, error) {
	for {
		instance, err := a.resourceInstanceFind(id)
		if err != nil {
			return nil, err
		}
		if instance == nil || instance.State == nil {
			continue
		}

		state := *instance.State.Name

		fmt.Println("--state -")
		fmt.Println(state, expectedState)

		if state == expectedState {
			return instance, nil
		}

		// sleep
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context done")
		case <-time.After(2 * time.Second):
		}
	}
}

func (a *AwsProvider) Update(ctx context.Context, node *proto.Node) error {
	if node.Id == "" {
		return fmt.Errorf("id not set")
	}
	if node.ExpectedConfig == "" {
		node.ExpectedConfig = "{}"
	}

	var config *config
	if err := json.Unmarshal([]byte(node.ExpectedConfig), &config); err != nil {
		return err
	}
	if config.Type == "" {
		config.Type = "t2.medium"
	}

	fmt.Println("-- dd")
	fmt.Println(node.ExpectedConfig)
	fmt.Println(node.CurrentConfig)

	if node.Handle != nil {

		// do only stuff if config is different
		if node.ExpectedConfig == node.CurrentConfig {
			panic("this should not happen, we always call this if we have somethign to update")
		}

		var handleInput handle
		if err := json.Unmarshal([]byte(node.Handle.Handle), &handleInput); err != nil {
			return err
		}

		_, err := a.conn.StopInstances(&ec2.StopInstancesInput{
			InstanceIds: []*string{aws.String(handleInput.InstanceID)},
		})
		if err != nil {
			return err
		}

		// wait for the instance to stop
		if _, err := a.waitForState(ctx, handleInput.InstanceID, ec2.InstanceStateNameStopped); err != nil {
			return err
		}

		// start again
		modifyInput := &ec2.ModifyInstanceAttributeInput{
			InstanceId: aws.String(handleInput.InstanceID),
			InstanceType: &ec2.AttributeValue{
				Value: aws.String(config.Type),
			},
		}
		if _, err := a.conn.ModifyInstanceAttribute(modifyInput); err != nil {
			return err
		}

		// start it again
		input := &ec2.StartInstancesInput{
			InstanceIds: []*string{aws.String(handleInput.InstanceID)},
		}
		if _, err := a.conn.StartInstances(input); err != nil {
			return err
		}

		// wait for it to be running and update handle
		instance, err := a.waitForState(ctx, handleInput.InstanceID, ec2.InstanceStateNameRunning)
		if err != nil {
			return err
		}

		awsHandle := &handle{
			InstanceID: *instance.InstanceId,
		}
		handleRaw, err := json.Marshal(awsHandle)
		if err != nil {
			return err
		}

		node.Handle = &proto.Node_Handle{
			Handle: string(handleRaw),
			Ip:     *instance.PublicIpAddress,
		}
		node.CurrentConfig = node.ExpectedConfig
		return nil
	}

	userData := base64.StdEncoding.EncodeToString([]byte(userdata.GetUserData()))

	fmt.Println(userData)

	instanceInput := &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-0341aeea105412b57"),
		InstanceType: aws.String(config.Type),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sdh"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize: aws.Int64(1),
				},
			},
		},
		KeyName:  aws.String("atlas"),
		UserData: aws.String(userData),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("MyFirstInstance"),
					},
				},
			},
		},
	}
	runResult, err := a.conn.RunInstances(instanceInput)
	if err != nil {
		return err
	}

	instanceId := *runResult.Instances[0].InstanceId
	instance, err := a.waitForState(ctx, instanceId, ec2.InstanceStateNameRunning)
	if err != nil {
		return err
	}

	awsHandle := &handle{
		InstanceID: *instance.InstanceId,
	}
	handleRaw, err := json.Marshal(awsHandle)
	if err != nil {
		return err
	}

	handle := &proto.Node_Handle{
		Handle: string(handleRaw),
		Ip:     *instance.PublicIpAddress,
	}

	node.Handle = handle
	node.CurrentConfig = node.ExpectedConfig
	return nil
}

type handle struct {
	InstanceID string
}
