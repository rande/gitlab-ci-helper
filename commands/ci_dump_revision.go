// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"flag"
	"github.com/mitchellh/cli"
	helper "github.com/rande/gitlab-ci-helper"
	"os"
	"strings"
)

type CiDumpRevisionCommand struct {
	Ui           cli.Ui
	Verbose      bool
	RevisionFile string
	Reference    string
}

func (c *CiDumpRevisionCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("ci:revision", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	cmdFlags.BoolVar(&c.Verbose, "verbose", false, "")
	cmdFlags.StringVar(&c.RevisionFile, "file", "REVISION", "The revision file")
	cmdFlags.StringVar(&c.Reference, "ref", helper.GetEnv("CI_COMMIT_SHA", os.Getenv("CI_BUILD_REF")), "The sha1 (default: env var CI_BUILD_REF)")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	fp, _ := os.Create(c.RevisionFile)
	defer fp.Close()

	fp.Write([]byte(c.Reference))

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
  -file               Target file (default: REVISION)
  -ref                The sha1 (default: env var CI_BUILD_REF)
  -verbose            Add verbose information to the output

Env Variables:

  CI_BUILD_REF        Get the revision from this variable
`
	return strings.TrimSpace(helpText)
}
