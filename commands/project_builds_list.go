// Copyright © 2016-present Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
	gitlab "github.com/plouc/go-gitlab-client"
	helper "github.com/rande/gitlab-ci-helper"
)

var (
	icon_green     = "🍏"
	icon_red       = "🍅"
	icon_pending   = "🍊"
	icon_artifacts = "🍞"
)

type ProjectBuildsListCommand struct {
	Ui      cli.Ui
	Verbose bool
}

func (c *ProjectBuildsListCommand) Run(args []string) int {

	flags := flag.NewFlagSet("server", flag.ContinueOnError)
	flags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	flags.BoolVar(&c.Verbose, "verbose", false, "")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		flags.Usage()

		c.Ui.Error(fmt.Sprintf("\nError: %s", "Invalid number of arguments"))

		return 1
	}

	config := helper.NewConfig()
	client := gitlab.NewGitlab(config.Gitlab.Host, config.Gitlab.ApiPath, config.Gitlab.Token)

	project, err := helper.GetProject(args[0], client)

	if err != nil {
		fmt.Println(err)

		return 1
	}

	c.Ui.Output(fmt.Sprintf("Project: %s/%s (id: %d)", project.Namespace.Name, project.Name, project.Id))

	builds, err := client.ProjectBuilds(strconv.FormatInt(int64(project.Id), 10))

	if err != nil {
		flags.Usage()

		c.Ui.Error(fmt.Sprintf("\nError: %s", err.Error()))

		return 1
	}

	for _, b := range builds {
		artifacts := " "

		if b.ArtifactsFile.Size > 0 {
			artifacts = icon_artifacts
		}

		status := icon_pending
		switch b.Status {
		case "success":
			status = icon_green
		case "failed":
			status = icon_red
		}

		c.Ui.Output(fmt.Sprintf(" > %s  %s % 4d - %-15s ref: %-25s short id: %s", status, artifacts, b.Id, b.Name, b.Ref, b.Commit.Short_Id))
	}

	return 0
}

func (c *ProjectBuildsListCommand) Synopsis() string {
	return "Return builds available for the provided project."
}

func (c *ProjectBuildsListCommand) Help() string {
	helpText := `
Usage: gitlab-ci-helper project:builds:list [options] project

  List all builds available for the provide project

Arguments:

  project             Can be an id or a string: namespace/name

Options:

  -verbose            Add verbose information to the output

Credentials are retrieved from environment:

  GITLAB_HOST         The gitlab host
  GITLAB_TOKEN        The user's token
  GITLAB_API_PATH     (optional) the api path, default to: "/api/v4"

`
	return strings.TrimSpace(helpText)
}
