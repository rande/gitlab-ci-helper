// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	gitlab "github.com/plouc/go-gitlab-client"
	helper "github.com/rande/gitlab-ci-helper"
	"strings"
)

type ProjectsListCommand struct {
	Ui      cli.Ui
	Verbose bool
}

func (c *ProjectsListCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("project:list", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	cmdFlags.BoolVar(&c.Verbose, "verbose", false, "")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	config := helper.NewConfig()

	client := gitlab.NewGitlab(config.Gitlab.Host, config.Gitlab.ApiPath, config.Gitlab.Token)

	c.Ui.Output("Trying to find project from options")

	projects, err := client.Projects()

	if err != nil {
		c.Ui.Error(err.Error())

		return 1
	}

	for _, p := range projects {
		c.Ui.Output(fmt.Sprintf(" > % 4d - %s - %s", p.Id, p.Name, p.Namespace.Name))
	}

	return 0
}

func (c *ProjectsListCommand) Synopsis() string {
	return "Return the list of projects available."
}

func (c *ProjectsListCommand) Help() string {
	helpText := `
Usage: gitlab-ci-helper project:list [options] project

  List all projects available

Options:

  -verbose            Add verbose information to the output

Credentials are retrieved from environment:

  GITLAB_HOST         The gitlab host
  GITLAB_TOKEN        The user's token
  GITLAB_API_PATH     (optional) the api path, default to: "/api/v3"

`
	return strings.TrimSpace(helpText)
}
