// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/mitchellh/cli"
	"github.com/rande/gitlab-ci-helper/commands"
	"os"
)

func main() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	c := cli.NewCLI("gitlab-helper", "0.0.1-DEV")
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
			return &commands.CiDumpInfoCommand{
				Ui: ui,
			}, nil
		},
		"ci:revision": func() (cli.Command, error) {
			return &commands.CiDumpRevisionCommand{
				Ui: ui,
			}, nil
		},
		"ci:notify:hipchat": func() (cli.Command, error) {
			return &commands.CiNotificationHipchatCommand{
				Ui: ui,
			}, nil
		},
		"dump:readme": func() (cli.Command, error) {
			return &commands.DumpReadmeCommand{
				Ui:       ui,
				Commands: c.Commands,
			}, nil
		},
	}

	exitStatus, _ := c.Run()

	os.Exit(exitStatus)
}
