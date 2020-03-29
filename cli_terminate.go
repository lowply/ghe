package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli/v2"
)

type terminate struct {
	ec2 *ec2svc
	ids []string
}

func (t *terminate) terminate() error {
	if t.ids[0] == "all" {
		instances, err := t.ec2.describeSortedInstances()
		if err != nil {
			return err
		}

		targets := []string{}

		for _, i := range instances {
			targets = append(targets, *i.InstanceId)
		}
		t.ids = targets
	}
	idPointers := []*string{}
	for _, i := range t.ids {
		idPointers = append(idPointers, aws.String(i))
	}

	input := &ec2.TerminateInstancesInput{
		InstanceIds: idPointers,
	}

	_, err := t.ec2.TerminateInstances(input)
	if err != nil {
		return err
	}

	return nil
}

func cmd_terminate(c *cli.Context) error {
	ec2, err := newEC2()
	if err != nil {
		return err
	}

	t := &terminate{
		ec2: ec2,
		ids: c.Args().Slice(),
	}

	err = t.terminate()
	if err != nil {
		return err
	}

	err = t.ec2.list()
	if err != nil {
		return err
	}

	return nil
}
