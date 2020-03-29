package main

import (
	"errors"
	"fmt"

	ver "github.com/hashicorp/go-version"
	"github.com/urfave/cli/v2"
)

const minimumVersion = "1.9.9"

type ami struct {
	version *ver.Version
	id      string
}

func newAmi(vs string) (*ami, error) {
	v, err := ver.NewSemver(vs)
	if err != nil {
		return nil, err
	}

	min, err := ver.NewSemver(minimumVersion)
	if err != nil {
		return nil, err
	}

	if v.LessThan(min) {
		return nil, errors.New("GitHub Enterprise version should be greater than 2.0.0")
	}

	y, err := newReleaseYaml(v)
	if err != nil {
		return nil, err
	}

	oldScheme, err := ver.NewSemver("2.3.0")
	if err != nil {
		return nil, err
	}

	id := ""

	if y.version.GreaterThan(oldScheme) {
		r := &releaseV2{}
		id, err = y.getId(r)
	} else {
		r := &releaseV1{}
		id, err = y.getId(r)
	}

	if err != nil {
		return nil, err
	}

	ami := &ami{
		version: v,
		id:      id,
	}

	return ami, nil
}

func cmd_ami(c *cli.Context) error {
	if c.Args().Len() != 1 {
		return errors.New("Wrong argument")
	}

	ami, err := newAmi(c.Args().Get(0))
	if err != nil {
		return err
	}

	fmt.Println(ami.id)

	return nil
}
