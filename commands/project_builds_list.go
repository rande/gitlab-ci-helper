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
	"strconv"
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

func (c *ProjectBuildsListCommand) Help() string {
	return `Return builds available for the provided project.`
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

		fmt.Printf("Error: %s", "Invalid number of arguments")

		return 1
	}

	config := helper.NewConfig()
	client := gitlab.NewGitlab(config.Host, config.ApiPath, config.Token)

	project, err := helper.GetProject(args[0], client)

	if err != nil {
		fmt.Println(err)

		return 1
	}

	c.Ui.Output(fmt.Sprintf("Project: %s/%s (id: %d)", project.Namespace.Name, project.Name, project.Id))

	builds, err := client.ProjectBuilds(strconv.FormatInt(int64(project.Id), 10))

	if err != nil {
		fmt.Printf("Error: %s", err.Error())

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

		c.Ui.Output(fmt.Sprintf(" > %s  %s % 4d - %-15s ref: %-25s short id: %s", status, artifacts, b.ID, b.Name, b.Ref, b.Commit.Short_Id))
	}

	return 0
}

func (c *ProjectBuildsListCommand) Synopsis() string {
	return "Return builds available for the provided project."
}
