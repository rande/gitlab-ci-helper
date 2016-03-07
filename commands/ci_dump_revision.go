// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"flag"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

type CiDumpRevisionCommand struct {
	Ui      cli.Ui
	Verbose bool
}

func (c *CiDumpRevisionCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("ci:revision", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	cmdFlags.BoolVar(&c.Verbose, "verbose", false, "")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	fp, _ := os.Create("REVISION")
	defer fp.Close()

	fp.Write([]byte(os.Getenv("CI_BUILD_REF")))

	return 0
}

func (c *CiDumpRevisionCommand) Synopsis() string {
	return "Dump a REVISION with the current sha1."
}

func (c *CiDumpRevisionCommand) Help() string {
	helpText := `
Usage: gitlab-ci-helper ci:revision [options]

  Dump a REVISION file with the current sha1

Options:

  -verbose            Add verbose information to the output

Env Variables:

  CI_BUILD_REF        Get the revision from this variable
`
	return strings.TrimSpace(helpText)
}
