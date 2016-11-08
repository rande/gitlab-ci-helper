// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/mitchellh/cli"
	"github.com/rande/gitlab-ci-helper/commands"
	"github.com/rande/gitlab-ci-helper/integrations/flowdock"
	"github.com/rande/gitlab-ci-helper/integrations/hipchat"
	"os"
)

var (
	Version = "0.0.1-Dev"
	RefLog  = "master"
)

func main() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	c := cli.NewCLI("gitlab-ci-helper", "0.0.1-DEV")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"project:list": func() (cli.Command, error) {
			return &commands.ProjectsListCommand{
				Ui: ui,
			}, nil
		},
		"project:builds": func() (cli.Command, error) {
			return &commands.ProjectBuildsListCommand{
				Ui: ui,
			}, nil
		},
		"project:builds:artifacts": func() (cli.Command, error) {
			return &commands.ProjectBuildArtifactCommand{
				Ui: ui,
			}, nil
		},
		"ci:meta": func() (cli.Command, error) {
			return &commands.CiDumpMetaCommand{
				Ui: ui,
			}, nil
		},
		"ci:revision": func() (cli.Command, error) {
			return &commands.CiDumpRevisionCommand{
				Ui: ui,
			}, nil
		},
		"hipchat:message": func() (cli.Command, error) {
			return &hipchat.CiNotificationHipchatCommand{
				Ui: ui,
			}, nil
		},
		"jenkins:start": func() (cli.Command, error) {
			return &commands.JenkinsStartJobCommand{
				Ui: ui,
			}, nil
		},
		"flowdock:message": func() (cli.Command, error) {
			return &flowdock.CiFlowdockMessageCommand{
				Ui: ui,
			}, nil
		},
		"flowdock:status": func() (cli.Command, error) {
			return &flowdock.CiFlowdockStatusCommand{
				Ui: ui,
			}, nil
		},
		"dump:readme": func() (cli.Command, error) {
			return &commands.DumpReadmeCommand{
				Ui:       ui,
				Commands: c.Commands,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &commands.VersionCommand{
				Ui:      ui,
				Version: Version,
				RefLog:  RefLog,
			}, nil
		},
		"s3:archive": func() (cli.Command, error) {
			return &commands.S3ArchiveCommand{
				Ui: ui,
			}, nil
		},
		"s3:extract": func() (cli.Command, error) {
			return &commands.S3ExtractCommand{
				Ui: ui,
			}, nil
		},
	}

	exitStatus, _ := c.Run()

	os.Exit(exitStatus)
}
