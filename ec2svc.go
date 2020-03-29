package main

import (
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
)

type ec2svc struct {
	*ec2.EC2
}

func newEC2() (*ec2svc, error) {
	cred := aws.Config{
		Credentials: credentials.NewSharedCredentials("", viper.GetString("aws.profile")),
		Region:      aws.String(viper.GetString("aws.region")),
	}

	sess, err := session.NewSession(&cred)
	if err != nil {
		return nil, err
	}

	e := &ec2svc{
		ec2.New(sess),
	}

	return e, nil
}

func (e *ec2svc) describeSortedInstances() ([]*ec2.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("key-name"),
				Values: []*string{
					aws.String(viper.GetString("aws.key_name")),
				},
			},
		},
	}

	result, err := e.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	instances := []*ec2.Instance{}

	for _, r := range result.Reservations {
	I:
		for _, i := range r.Instances {
			// Exclude "terminated" status
			if *i.State.Code == 48 {
				continue
			}
			// Exclude instances
			if len(viper.GetStringSlice("aws.excludes")) != 0 {
				for _, e := range viper.GetStringSlice("aws.excludes") {
					if e == *i.InstanceId {
						continue I
					}
				}
			}
			instances = append(instances, i)
		}
	}

	// Sort instances by the launch time
	sort.Slice(instances, func(i, j int) bool {
		vi := *instances[i].LaunchTime
		vj := *instances[j].LaunchTime
		return vi.Before(vj)

	})

	return instances, nil
}

func (e *ec2svc) list() error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Date",
		"Instance ID",
		"Name",
		"Public DNS",
		"Public IP",
		"State",
		"Version",
	})

	instances, err := e.describeSortedInstances()
	if err != nil {
		return err
	}

	for _, i := range instances {
		var name, version string
		for _, t := range i.Tags {
			if *t.Key == "Name" {
				name = *t.Value
			}
			if *t.Key == "Version" {
				version = *t.Value
			}
		}

		var publicdns, publicip string

		if i.PublicDnsName != nil {
			publicdns = *i.PublicDnsName
		}

		if i.PublicIpAddress != nil {
			publicip = *i.PublicIpAddress
		}

		table.Append([]string{
			i.LaunchTime.String(),
			*i.InstanceId,
			name,
			publicdns,
			publicip,
			*i.State.Name,
			version,
		})
	}
	table.Render()
	return nil
}
