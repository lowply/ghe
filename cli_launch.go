package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

type launch struct {
	ec2     *ec2svc
	version string
	name    string
	replica bool
}

func (l *launch) launch() error {
	ami, err := newAmi(l.version)
	if err != nil {
		return err
	}

	input := &ec2.RunInstancesInput{
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/xvdf"),
				Ebs: &ec2.EbsBlockDevice{
					DeleteOnTermination: aws.Bool(true),
					Iops:                aws.Int64(1000),
					VolumeType:          aws.String("io1"),
					VolumeSize:          aws.Int64(100),
				},
			},
		},
		ImageId:      aws.String(ami.id),
		InstanceType: aws.String(viper.GetString("aws.type")),
		KeyName:      aws.String(viper.GetString("aws.key_name")),
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		SecurityGroupIds: []*string{
			aws.String(viper.GetString("aws.sec_group")),
		},
		SubnetId: aws.String(viper.GetString("aws.subnet_id")),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(l.name),
					},
					{
						Key:   aws.String("Version"),
						Value: aws.String(l.version),
					},
				},
			},
		},
	}

	reservation, err := l.ec2.RunInstances(input)
	if err != nil {
		return err
	}

	instance, err := l.wait(*reservation.Instances[0].InstanceId)
	if err != nil {
		return err
	}

	if !l.replica {
		fmt.Println("Updating dns...")

		d := &dns{
			ip: *instance.PublicIpAddress,
		}

		err := d.dns()
		if err != nil {
			return err
		}
	}

	err = l.ec2.list()
	if err != nil {
		return err
	}

	return nil
}

func (l *launch) wait(id string) (*ec2.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
	}
	for {
		fmt.Println("Waiting for " + id + " ...")
		result, err := l.ec2.DescribeInstances(input)
		if err != nil {
			return nil, err
		}
		if *result.Reservations[0].Instances[0].State.Code == 0 {
			time.Sleep(1 * time.Second)
			continue
		} else {
			fmt.Println("Launched: " + id)
			return result.Reservations[0].Instances[0], nil
		}
	}
}

func cmd_launch(c *cli.Context) error {
	if c.Args().Len() > 2 {
		return errors.New("Usage: launch <version> <replica>")
	}

	ec2, err := newEC2()
	if err != nil {
		return err
	}

	name := viper.GetString("basic.domain")
	replica := c.Args().Get(1) == "replica"
	if replica {
		name = name + " replica"
	}

	l := &launch{
		version: c.Args().Get(0),
		ec2:     ec2,
		name:    name,
		replica: replica,
	}

	err = l.launch()
	if err != nil {
		return err
	}

	return nil
}
