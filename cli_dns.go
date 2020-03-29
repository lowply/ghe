package main

import "github.com/urfave/cli/v2"

type dns struct {
	ip string
}

func (d *dns) dns() error {
	m, err := newDnsmasq()
	if err != nil {
		return err
	}

	err = m.Update(d.ip)
	if err != nil {
		return err
	}
	return nil
}

func cmd_dns(c *cli.Context) error {
	d := &dns{
		ip: c.Args().Get(0),
	}

	err := d.dns()
	if err != nil {
		return err
	}

	return nil
}
