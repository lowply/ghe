package main

import (
	"github.com/spf13/viper"
)

type release interface {
	getId() string
}

type releaseV1 struct {
	Amis []struct {
		ID     string `yaml:"id"`
		Region string `yaml:"region"`
	}
}

type releaseV2 struct {
	Version  string `yaml:"version"`
	Packages struct {
		Ami struct {
			Releases []struct {
				ID     string `yaml:"id"`
				Region string `yaml:"region"`
			} `yaml:"releases"`
		} `yaml:"ami"`
	}
}

func (r *releaseV1) getId() string {
	for _, v := range r.Amis {
		if v.Region == viper.GetString("aws.region") {
			return v.ID
		}
	}
	return ""
}

func (r *releaseV2) getId() string {
	for _, v := range r.Packages.Ami.Releases {
		if v.Region == viper.GetString("aws.region") {
			return v.ID
		}
	}
	return ""
}
