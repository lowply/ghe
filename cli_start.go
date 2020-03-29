package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli/v2"
)

type start struct {
	ec2 *ec2svc
	ids []string
}

func (s *start) start() error {
	idPointers := []*string{}
	for _, v := range s.ids {
		idPointers = append(idPointers, aws.String(v))
	}

	input := &ec2.StartInstancesInput{
		InstanceIds: idPointers,
	}

	_, err := s.ec2.StartInstances(input)
	if err != nil {
		return err
	}

	return nil
}

func cmd_start(c *cli.Context) error {
	ec2, err := newEC2()
	if err != nil {
		return err
	}

	s := &start{
		ec2: ec2,
		ids: c.Args().Slice(),
	}

	err = s.start()
	if err != nil {
		return err
	}

	err = s.ec2.list()
	if err != nil {
		return err
	}

	return nil
}
