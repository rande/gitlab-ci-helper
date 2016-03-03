package main

import (
	"github.com/mitchellh/cli"
	"github.com/rande/gitlab-ci-helper/commands"
	"log"
	"os"
)

func main() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	c := cli.NewCLI("gitlab-helper", "0.0.1-DEV")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"build:artifacts": func() (cli.Command, error) {
			return &commands.BuildDownloadArtifactsCommand{
				Ui: ui,
			}, nil
		},
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
	}

	exitStatus, err := c.Run()

	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
