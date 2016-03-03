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
)

type BuildDownloadArtifactsCommand struct {
	Ui        cli.Ui
	Verbose   bool
	JobName   string
	BuildId   string
	CommitId  string
	ProjectId string
}

func (c *BuildDownloadArtifactsCommand) Help() string {
	return `Serve gonode server (better be behing a http reverse proxy)`
}

func (c *BuildDownloadArtifactsCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("server", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	cmdFlags.BoolVar(&c.Verbose, "verbose", false, "")

	cmdFlags.StringVar(&c.JobName, "name", "package", "The job name related to the artifact")
	cmdFlags.StringVar(&c.BuildId, "build", "", "The build id")
	cmdFlags.StringVar(&c.CommitId, "commit", "", "The commit id")
	cmdFlags.StringVar(&c.ProjectId, "project", "", "The project id")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Job name: %s", c.JobName))
	c.Ui.Output(fmt.Sprintf("BuildId: %s", c.BuildId))
	c.Ui.Output(fmt.Sprintf("CommitId: %s", c.CommitId))
	c.Ui.Output(fmt.Sprintf("ProjectId: %s", c.ProjectId))

	config := helper.NewConfig()

	gitlab := gitlab.NewGitlab(config.Host, config.ApiPath, config.Token)

	c.Ui.Output("Trying to find project from options")

	gitlab.Projects()

	return 0
}

func (c *BuildDownloadArtifactsCommand) Synopsis() string {
	return "download artifact from a specific build"
}
