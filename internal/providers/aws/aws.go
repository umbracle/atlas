package aws

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/umbracle/atlas/internal/framework"
	"github.com/umbracle/atlas/internal/proto"
	"github.com/umbracle/atlas/internal/schema"
	"github.com/umbracle/atlas/internal/userdata"

	"fmt"
)

var _ framework.Provider = &AwsProvider{}

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

func (a *AwsProvider) Config() interface{} {
	return &config{}
}

func (a *AwsProvider) Schema() *schema.Object {
	return &schema.Object{
		Fields: map[string]*schema.Field{
			"instance_type": {
				Type:    &schema.String,
				Default: "t2.small",
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

func (a *AwsProvider) Update(ctx context.Context, old, new interface{}, node *proto.Node) error {
	if node.Id == "" {
		return fmt.Errorf("id not set")
	}

	var instance *ec2.Instance

	if old == nil {
		// create the instance
		res, err := a.createInstance(ctx, new.(*config))
		if err != nil {
			return err
		}
		instance = res
	} else {
		var handleInput handle
		if err := json.Unmarshal([]byte(node.Handle.Handle), &handleInput); err != nil {
			return err
		}

		// update
		res, err := a.updateInstance(ctx, handleInput.InstanceID, old.(*config), new.(*config))
		if err != nil {
			return err
		}
		instance = res
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
	return nil
}

func (a *AwsProvider) updateInstance(ctx context.Context, instanceID string, old, new *config) (*ec2.Instance, error) {
	_, err := a.conn.StopInstances(&ec2.StopInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})
	if err != nil {
		return nil, err
	}

	// wait for the instance to stop
	if _, err := a.waitForState(ctx, instanceID, ec2.InstanceStateNameStopped); err != nil {
		return nil, err
	}

	// start again
	modifyInput := &ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(instanceID),
		InstanceType: &ec2.AttributeValue{
			Value: aws.String(new.Type),
		},
	}
	if _, err := a.conn.ModifyInstanceAttribute(modifyInput); err != nil {
		return nil, err
	}

	// start it again
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}
	if _, err := a.conn.StartInstances(input); err != nil {
		return nil, err
	}

	// wait for it to be running and update handle
	instance, err := a.waitForState(ctx, instanceID, ec2.InstanceStateNameRunning)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (a *AwsProvider) createInstance(ctx context.Context, config *config) (*ec2.Instance, error) {
	userDataInput, err := userdata.GetUserData("https://7693-88-9-192-173.ngrok.io/atlas")
	if err != nil {
		return nil, err
	}
	userData := base64.StdEncoding.EncodeToString([]byte(userDataInput))

	if config.Type == "" {
		config.Type = "t2.small"
	}
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
					VolumeSize: aws.Int64(100),
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
		return nil, err
	}

	instanceId := *runResult.Instances[0].InstanceId
	instance, err := a.waitForState(ctx, instanceId, ec2.InstanceStateNameRunning)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

type handle struct {
	InstanceID string
}
