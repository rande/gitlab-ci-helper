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
	"io"
	"os"
	"strconv"
	"strings"
)

type ProjectBuildArtifactCommand struct {
	Ui            cli.Ui
	Verbose       bool
	ExtractPath   string
	ArtifactsFile string
	BuildId       string
	Stage         string
	Ref           string
}

func (c *ProjectBuildArtifactCommand) Run(args []string) int {

	flags := flag.NewFlagSet("project:builds:artifacts", flag.ContinueOnError)
	flags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	flags.BoolVar(&c.Verbose, "verbose", false, "")
	flags.StringVar(&c.ArtifactsFile, "file", "artifacts.zip", "Artifacts file name")
	flags.StringVar(&c.ExtractPath, "path", "", "The path to extract the artifacts")
	flags.StringVar(&c.BuildId, "build", "", "The build number to get the artifacts")

	flags.StringVar(&c.Stage, "stage", "package", "The stage to search the artifacts")
	flags.StringVar(&c.Ref, "ref", "", "The reference (sha1) to search the artifacts")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()

	if len(args) != 1 {
		flags.Usage()

		c.Ui.Error(fmt.Sprintf("Error: %s", "Invalid number of arguments"))

		return 1
	}

	config := helper.NewConfig()
	client := gitlab.NewGitlab(config.Gitlab.Host, config.Gitlab.ApiPath, config.Gitlab.Token)

	project, err := helper.GetProject(args[0], client)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error: %s", err.Error()))

		return 1
	}

	c.Ui.Output(fmt.Sprintf("Found project: %s/%s (id: %d)", project.Namespace.Name, project.Name, project.Id))

	var build *gitlab.Build

	if len(c.BuildId) > 0 {
		build, err = helper.GetBuild(project, c.BuildId, client)

		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error: %s", err.Error()))

			return 1
		}

	} else if len(c.Stage) > 0 && len(c.Ref) > 0 {
		builds, err := client.ProjectBuilds(strconv.FormatInt(int64(project.Id), 10))

		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error: %s", err.Error()))

			return 1
		}

		for _, b := range builds {
			if b.Stage == c.Stage && b.Commit.Id == c.Ref {
				build = b
				break
			}
		}
	}

	if build == nil {
		c.Ui.Error(fmt.Sprintf("Error: %s", "Unable to found the build"))

		return 1
	}

	c.Ui.Output(fmt.Sprintf("Found build - stage:%s status:%s id:%d", build.Stage, build.Status, build.Id))

	fp, err := os.Create(c.ArtifactsFile)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error: %s", err.Error()))

		return 1
	}

	r, err := client.ProjectBuildArtifacts(strconv.FormatInt(int64(project.Id), 10), strconv.FormatInt(int64(build.Id), 10))

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error: %s", err.Error()))

		return 1
	}

	c.Ui.Output(fmt.Sprintf("Downloading artifacts... (%s)", fp.Name()))
	_, err = io.Copy(fp, r)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error: %s", err.Error()))

		return 1
	}

	if len(c.ExtractPath) > 0 {
		c.Ui.Output(fmt.Sprintf("Extracting package... (%s)", c.ExtractPath))

		err = helper.Unzip(c.ArtifactsFile, "package")

		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error: %s", err.Error()))

			return 1
		}
	}

	c.Ui.Output(fmt.Sprintf("Done!"))

	return 0
}

func (c *ProjectBuildArtifactCommand) Synopsis() string {
	return "Download artifact from a job."
}

func (c *ProjectBuildArtifactCommand) Help() string {
	helpText := `
Usage: gitlab-helper project:builds:artifacts [options] project build

  Download an artifacts and extract it if the 'path' option is provided

Options:

  -build=XX           The build number used to retrieved the related artifact
  -stage=XX           The stage to search the build (must be used with -ref, default: package)
  -ref=XX             The sha1 linked to the build (must be used with -stage)
  -file=artifacts.zip The path to the artifact file (default: artifacts.zip)
  -path=./package     The path to extract the command. If not set, the artifact will not
                      be extracted.
  -verbose            Add verbose information to the output
`

	return strings.TrimSpace(helpText)
}
