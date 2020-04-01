package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

type configure struct {
	domain     string
	port       string
	baseurl    string
	configPath string
	data       Config
}

func newConfigure() *configure {
	c := &configure{
		domain: viper.GetString("basic.domain"),
		port:   "8443",
	}
	c.baseurl = "https://api_key:" + viper.GetString("basic.password") + "@" + c.dp() + "/setup/api"
	c.configPath = os.Getenv("HOME") + "/.config/ghe-config.json"
	return c
}

func (c *configure) dp() string {
	return net.JoinHostPort(c.domain, c.port)
}

func (c *configure) oneliner(path string) (string, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	n := strings.NewReplacer(
		"\r\n", "\n",
		"\r", "\n",
		"\n", "\n",
	).Replace(string(file))

	return n, nil
}

func (c *configure) generate_settings() error {
	ssl_cert, err := c.oneliner(viper.GetString("basic.ssl_cert"))
	if err != nil {
		return err
	}

	ssl_key, err := c.oneliner(viper.GetString("basic.ssl_key"))
	if err != nil {
		return err
	}

	syslog_cert, err := c.oneliner(viper.GetString("syslog.cert"))
	if err != nil {
		return err
	}

	// Basic
	c.data.Enterprise.PrivateMode = true
	c.data.Enterprise.SubdomainIsolation = true
	c.data.Enterprise.GithubHostname = c.domain
	c.data.Enterprise.AuthMode = "ldap"
	c.data.Enterprise.BuiltinAuthFallback = true
	c.data.Enterprise.AdminPassword = viper.GetString("basic.password")
	c.data.Enterprise.GithubSsl.Enabled = true
	c.data.Enterprise.GithubSsl.Cert = ssl_cert
	c.data.Enterprise.GithubSsl.Key = ssl_key

	// LDAP
	c.data.Enterprise.Ldap.Host = viper.GetString("ldap.host")
	c.data.Enterprise.Ldap.Port = viper.GetInt("ldap.port")
	c.data.Enterprise.Ldap.Base = []string{viper.GetString("ldap.base")}
	c.data.Enterprise.Ldap.UID = "uid"
	c.data.Enterprise.Ldap.BindDn = viper.GetString("ldap.dn")
	c.data.Enterprise.Ldap.Password = viper.GetString("ldap.pass")
	c.data.Enterprise.Ldap.Method = viper.GetString("ldap.method")
	c.data.Enterprise.Ldap.AdminGroup = viper.GetString("ldap.admin")
	c.data.Enterprise.Ldap.PosixSupport = true
	c.data.Enterprise.Ldap.UserSyncEmails = true
	c.data.Enterprise.Ldap.UserSyncKeys = true
	c.data.Enterprise.Ldap.UserSyncInterval = 1
	c.data.Enterprise.Ldap.TeamSyncInterval = 1
	c.data.Enterprise.Ldap.SyncEnabled = true
	c.data.Enterprise.Ldap.ExternalAuthTokenRequired = true
	c.data.Enterprise.Ldap.VerifyCertificate = true
	c.data.Enterprise.Ldap.Profile.UID = "uid"
	c.data.Enterprise.Ldap.Profile.Mail = "mail"
	c.data.Enterprise.Ldap.Profile.Key = "sshPublicKey"

	// SAML
	c.data.Enterprise.Saml.SsoURL = viper.GetString("saml.url")
	c.data.Enterprise.Saml.Certificate = viper.GetString("saml.cert")
	c.data.Enterprise.Saml.Issuer = viper.GetString("saml.issuer")
	c.data.Enterprise.Saml.NameIDFormat = viper.GetString("saml.format")

	// SMTP
	c.data.Enterprise.SMTP.Enabled = true
	c.data.Enterprise.SMTP.Address = viper.GetString("smtp.host")
	c.data.Enterprise.SMTP.Port = 587
	c.data.Enterprise.SMTP.SupportAddress = viper.GetString("smtp.support_email")
	c.data.Enterprise.SMTP.SupportAddressType = "email"
	c.data.Enterprise.SMTP.NoreplyAddress = viper.GetString("smtp.noreply_email")

	// Syslog
	c.data.Enterprise.Syslog.Server = viper.GetString("syslog.host")
	c.data.Enterprise.Syslog.Cert = syslog_cert
	c.data.Enterprise.Syslog.ProtocolName = "tcp"
	c.data.Enterprise.Syslog.TLSEnabled = true

	// Pages
	c.data.Enterprise.Pages.Enabled = true

	j, err := json.Marshal(c.data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.configPath, j, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *configure) initconfig() error {
	fmt.Printf("Initializing the instance...\n")
	target := "https://" + c.dp() + "/setup/api/start"

	file, err := os.Open(viper.GetString("basic.license"))
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("license", file.Name())
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	password, err := writer.CreateFormField("password")
	if err != nil {
		return err
	}

	_, err = password.Write([]byte(viper.GetString("basic.password")))
	if err != nil {
		return err
	}

	settings, err := writer.CreateFormField("settings")
	if err != nil {
		return err
	}

	configfile, err := ioutil.ReadFile(c.configPath)
	if err != nil {
		return err
	}

	_, err = settings.Write(configfile)
	if err != nil {
		return err
	}

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, target, body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if len(content) > 0 {
		// em := &ErrorMessage{}
		em := &struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}{}
		json.Unmarshal(content, em)
		if em.Error == "password-error" {
			fmt.Printf("Message in initconfig: %s\n", em.Message)
			return nil
		} else {
			return errors.New(fmt.Sprintf("Error in initconfig: %s - %s\n", em.Error, em.Message))
		}
	}

	return nil
}

func (c *configure) sendconfig() error {
	fmt.Printf("Sending configuration...\n")

	target := c.baseurl + "/settings"

	file, err := ioutil.ReadFile(c.configPath)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("settings", string(file))

	buf := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest(http.MethodPut, target, buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if len(content) > 0 {
		return errors.New(fmt.Sprintf("Error in sendconfig: %s\n", content))
	}

	return nil
}

func (c *configure) applyconfig() error {
	target := c.baseurl + "/configure"
	body := &bytes.Buffer{}

	req, err := http.NewRequest(http.MethodPost, target, body)
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if len(content) > 0 {
		return errors.New(fmt.Sprintf("Error in applyconfig: %s\n", content))
	}

	return nil
}

func (c *configure) addkey() error {
	target := c.baseurl + "/settings/authorized-keys"

	file, err := os.Open(viper.GetString("basic.pubkey"))
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("authorized_key", file.Name())
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, target, body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if len(content) > 0 {
		fmt.Println(string(content))
		kp := []struct {
			Key         string `json:"key"`
			PrettyPrint string `json:"pretty-print"`
			Comment     string `json:"comment"`
		}{}
		json.Unmarshal(content, kp)
		fmt.Println("Keys:")
		for _, v := range kp {
			fmt.Printf("%s - %s - %s\n", v.Key, v.PrettyPrint, v.Comment)
		}
	}

	return nil
}

func (c *configure) check_port() error {
	for {
		conn, err := net.DialTimeout("tcp", c.dp(), 3*time.Second)
		if err != nil {
			fmt.Println("Waiting for port " + c.port + " ...")
			time.Sleep(1 * time.Second)
			continue
		}
		defer conn.Close()
		break
	}
	fmt.Println("Port " + c.port + " is ready")
	return nil
}

func (c *configure) checkprogress() error {
	target := c.baseurl + "/configcheck"
	body := &bytes.Buffer{}

	req, err := http.NewRequest(http.MethodGet, target, body)
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	s := &struct {
		Status   string `json:"status"`
		Progress []struct {
			Key    string `json:"key"`
			Status string `json:"status"`
		} `json:"progress"`
	}{}
	fmt.Println("Configuration in progress...")

	for {
		res, err := client.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(content, s)
		if err != nil {
			return err
		}

		if s.Status != "running" {
			fmt.Println("Done")
			break
		}

		table := tablewriter.NewWriter(os.Stdout)
		statusheader := []string{"Status"}
		statusline := []string{s.Status}
		for _, v := range s.Progress {
			statusheader = append(statusheader, v.Key)
			statusline = append(statusline, v.Status)
		}
		table.SetHeader(statusheader)
		table.Append(statusline)
		table.Render()
		time.Sleep(6 * time.Second)
	}

	if s.Status == "success" {
		fmt.Println(c.domain + " is now ready!")
	} else {
		fmt.Println("Failed to configure " + c.domain + " :(")
		fmt.Println("Status: " + s.Status)
	}

	return nil
}

func (c *configure) updateconfig() error {
	err := c.sendconfig()
	if err != nil {
		return err
	}

	err = c.applyconfig()
	if err != nil {
		return err
	}

	return nil
}

func cmd_configure(c *cli.Context) error {
	co := newConfigure()

	err := co.generate_settings()
	if err != nil {
		return err
	}

	err = co.check_port()
	if err != nil {
		return err
	}

	err = co.initconfig()
	if err != nil {
		return err
	}

	err = co.applyconfig()
	if err != nil {
		return err
	}

	err = co.checkprogress()
	if err != nil {
		return err
	}

	err = co.addkey()
	if err != nil {
		return err
	}

	ec2, err := newEC2()
	if err != nil {
		return err
	}

	err = ec2.list()
	if err != nil {
		return err
	}

	return nil
}
