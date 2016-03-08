// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	"strings"
)

type VersionCommand struct {
	Ui       cli.Ui
	Extended bool
	Version  string
	RefLog   string
}

func (c *VersionCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("version", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	cmdFlags.BoolVar(&c.Extended, "e", false, "")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if c.Extended {
		c.Ui.Output(fmt.Sprintf("%s - %s", c.Version, c.RefLog))
	} else {
		c.Ui.Output(c.Version)
	}

	return 0
}

func (c *VersionCommand) Synopsis() string {
	return "Display the application version."
}

func (c *VersionCommand) Help() string {
	helpText := `
Usage: gitlab-ci-helper version

  Display the version number

Options:

  -e                  Extended version with sha1

`
	return strings.TrimSpace(helpText)
}
