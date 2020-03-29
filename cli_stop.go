package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli/v2"
)

type stop struct {
	ec2 *ec2svc
	ids []string
}

func (s *stop) stop() error {
	if s.ids[0] == "all" {
		instances, err := s.ec2.describeSortedInstances()
		if err != nil {
			return err
		}

		targets := []string{}

		for _, i := range instances {
			// State code 16 is "running"
			// ref: https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#InstanceState
			if *i.State.Code == 16 {
				targets = append(targets, *i.InstanceId)
			}
		}
		s.ids = targets
	}

	idPointers := []*string{}
	for _, i := range s.ids {
		idPointers = append(idPointers, aws.String(i))
	}

	input := &ec2.StopInstancesInput{
		InstanceIds: idPointers,
	}

	_, err := s.ec2.StopInstances(input)
	if err != nil {
		return err
	}

	return nil
}

func cmd_stop(c *cli.Context) error {
	ec2, err := newEC2()
	if err != nil {
		return err
	}

	s := &stop{
		ec2: ec2,
		ids: c.Args().Slice(),
	}

	err = s.stop()
	if err != nil {
		return err
	}

	err = s.ec2.list()
	if err != nil {
		return err
	}

	return nil
}
