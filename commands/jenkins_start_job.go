// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"flag"
	"fmt"
	jenkins "github.com/bndr/gojenkins"
	"github.com/mitchellh/cli"
	helper "github.com/rande/gitlab-ci-helper"
	"os"
	"strings"
)

type JenkinsStartJobCommand struct {
	Ui            cli.Ui
	Verbose       bool
	Job           string
	JobToken      string
	JobParameters helper.Parameters
}

func (c *JenkinsStartJobCommand) Run(args []string) int {
	flags := flag.NewFlagSet("jenkins:start", flag.ContinueOnError)
	flags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	flags.BoolVar(&c.Verbose, "verbose", false, "")
	flags.StringVar(&c.Job, "job", "", "Job to launch")

	flags.StringVar(&c.JobToken, "job-token", os.Getenv("JENKINS_JOB_TOKEN"), "The Token to launch the job")

	c.JobParameters = make(helper.Parameters, 0)

	flags.Var(&c.JobParameters, "parameter", "-parameter name:value")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()

	var jobArgs = make(map[string]string)

	for _, unparsedParam := range c.JobParameters {
		param := strings.Split(unparsedParam, ":")
		jobArgs[param[0]] = param[1]
	}

	// Command parameters have been parsed, now to processing...

	config := helper.NewConfig()
	var job *jenkins.Job

	if client, err := jenkins.CreateJenkins(config.Jenkins.Host, config.Jenkins.User, config.Jenkins.ApiToken).Init(); err != nil {
		c.Ui.Error(fmt.Sprintf("Could not find Jenkins %s", c.Job))

		return 1
	} else if job, err = client.GetJob(c.Job); err != nil {
		c.Ui.Error(fmt.Sprintf("Could not find job %s", c.Job))

		return 1
	}

	_, err := job.Invoke(nil, false, jobArgs, "", c.JobToken)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Could not launch job %s", job.GetName()))
		c.Ui.Error(err.Error())

		return 1
	}

	c.Ui.Output(fmt.Sprintf("Done!"))

	return 0
}

func (c *JenkinsStartJobCommand) Synopsis() string {
	return "Start a Jenkins job."
}

func (c *JenkinsStartJobCommand) Help() string {
	helpText := `
Usage: gitlab-ci-helper jenkins:start [options]

  Download an artifacts and extract it if the 'path' option is provided

Options:

  -job=XX               The Jenkins job to start
  -job-token=XX         The job token to use to call the API (default: JENKINS_JOB_TOKEN)
  -parameter=name:value Add a parameter to the job
  -verbose              Add verbose information to the output

Credentials are retrieved from environment:

  JENKINS_HOST      The jenkins host
  JENKINS_USER      The username to login with
  JENKINS_API_TOKEN The API Token associated with this user
`

	return strings.TrimSpace(helpText)
}
