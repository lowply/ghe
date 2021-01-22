package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	ver "github.com/hashicorp/go-version"
	yaml "gopkg.in/yaml.v2"
)

type releaseYaml struct {
	version *ver.Version
	content []byte
}

func newReleaseYaml(v *ver.Version) (*releaseYaml, error) {
	y := &releaseYaml{
		version: v,
	}

	nv := y.version.String()
	if v.Prerelease() != "" {
		nv = strings.Replace(nv, "-", ".", -1)
		nv = strings.Replace(nv, "rc.", "rc", -1)
		fmt.Println(nv)
	}

	baseURL := "https://github-enterprise.s3.amazonaws.com"
	url := baseURL + "/release/release-" + nv + ".yml"
	fmt.Println("Getting " + url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Error: failed to load yaml, " + resp.Status)
	}

	defer resp.Body.Close()

	y.content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return y, nil
}

func (y *releaseYaml) getId(r release) (string, error) {
	err := yaml.Unmarshal(y.content, r)
	if err != nil {
		return "", err
	}

	version := r.getId()

	if version == "" {
		return "", errors.New("AMI not found")
	}
	return version, nil
}
