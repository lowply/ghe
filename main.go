package main

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const version = "0.0.2"

var commands = []*cli.Command{
	&cli.Command{
		Name:        "ami",
		Usage:       "Get the Amazon AMI ID of the given GitHub Enterprise Server version, based on your AWS region",
		Description: "ghe ami 2.20.0",
		Action:      cmd_ami,
	},
	&cli.Command{
		Name:        "launch",
		Usage:       "Launch a GitHub Enterprise Server instance of the given version on EC2",
		Description: "",
		Action:      cmd_launch,
	},
	&cli.Command{
		Name:        "list",
		Usage:       "List all GitHub Enterprise Server instance with detailed information",
		Description: "",
		Action:      cmd_list,
	},
	&cli.Command{
		Name:        "start",
		Usage:       "Start GitHub Enterprise Server instances. Multiple argument can be passed",
		Description: "",
		Action:      cmd_start,
	},
	&cli.Command{
		Name:        "stop",
		Usage:       "Stop GitHub Enterprise Server instances. Multiple argument, or 'all' can be passed",
		Description: "",
		Action:      cmd_stop,
	},
	&cli.Command{
		Name:        "terminate",
		Usage:       "Terminate GitHub Enterprise Server instances. Multiple argument, or 'all' can be passed",
		Description: "",
		Action:      cmd_terminate,
	},
	&cli.Command{
		Name:        "dns",
		Usage:       "Update your dnsmasq config file to make your domain point to specific IP address",
		Description: "",
		Action:      cmd_dns,
	},
	&cli.Command{
		Name:        "configure",
		Usage:       "Initiate GitHub Enterprise Server configuration using the configuration file",
		Description: "",
		Action:      cmd_configure,
	},
}

func main() {
	viper.SetConfigName("ghe")
	viper.AddConfigPath("$HOME/.config/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	r := newResolver()
	err = r.check()
	if err != nil {
		log.Fatal(err)
	}

	app := cli.NewApp()
	app.Name = "ghe"
	app.Usage = "Launch GitHub Enterprise instances on AWS"
	app.Version = version
	app.Authors = []*cli.Author{{
		Name:  "lowply",
		Email: "lowply@github.com",
	}}
	app.Commands = commands

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
