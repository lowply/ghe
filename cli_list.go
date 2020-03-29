package main

import (
	"github.com/urfave/cli/v2"
)

func cmd_list(c *cli.Context) error {
	ec2, err := newEC2()
	if err != nil {
		return err
	}

	ec2.list()

	return nil
}
