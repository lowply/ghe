package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type dnsmasq struct {
	path   string
	domain string
	ip     net.IP
}

func newDnsmasq() (*dnsmasq, error) {
	cmd := exec.Command("dnsmasq", "-v")
	err := cmd.Run()
	if err != nil {
		return nil, errors.New("dnsmasq not found")
	}

	path := "/dnsmasq.d/ghe.conf"

	if runtime.GOARCH == "arm64" {
		path = "/opt/homebrew/etc" + path
	} else {
		path = "/usr/local/etc" + path
	}

	d := &dnsmasq{
		path:   path,
		domain: viper.GetString("basic.domain"),
	}

	_, err = os.Stat(d.path)
	if err != nil {
		fmt.Println("Creating " + d.path + " ...")
		_, err := os.Create(d.path)
		if err != nil {
			return nil, err
		}
	}

	file, err := os.Open(d.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), viper.GetString("basic.domain")) {
			ar := strings.Split(scanner.Text(), "/")
			if net.ParseIP(ar[2]) == nil {
				return nil, errors.New("Invalid IP Address: " + ar[2])
			}
			d.ip = net.ParseIP(ar[2])
		}
	}
	return d, nil
}

func (d *dnsmasq) Update(ipstring string) error {
	ip := net.ParseIP(ipstring)

	if ip == nil {
		return errors.New("Invalid argument: " + ipstring + ". Needs an IP address")
	}

	if d.ip.Equal(ip) {
		return errors.New("No change")
	}

	data := []byte("address=/" + d.domain + "/" + ip.String() + "\n")
	err := ioutil.WriteFile(d.path, data, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Restarting dnsmasq...")
	err = exec.Command("brew", "services", "restart", "dnsmasq").Run()
	if err != nil {
		return err
	}

	return nil
}
