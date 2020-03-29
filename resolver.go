package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type resolver struct {
	nameserver string
	dnsport    string
	path       string
}

func newResolver() *resolver {
	r := &resolver{
		nameserver: "127.0.0.1",
		dnsport:    "65353",
		path:       "/etc/resolver/" + viper.GetString("basic.domain"),
	}
	return r
}

func (r *resolver) check() error {
	_, err := os.Stat(r.path)

	if err != nil {
		fmt.Printf("Creating " + r.path + "\n")
		data := []byte("nameserver 127.0.0.1\nport 65353\n")
		err := ioutil.WriteFile(r.path, data, 0644)
		if err != nil {
			return err
		}
	} else {
		file, err := os.Open(r.path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			f := strings.Fields(scanner.Text())
			if f[0] == "nameserver" {
				if f[1] != r.nameserver {
					return errors.New("'nameserver' in " + r.path + " should be " + r.nameserver)
				}
			}
			if f[0] == "port" {
				if f[1] != r.dnsport {
					return errors.New("'port' in " + r.path + " should be " + r.dnsport)
				}
			}
		}
	}

	return nil
}
