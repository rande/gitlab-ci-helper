// Copyright © 2016-present Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package flowdock

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
	gitlab "gopkg.in/plouc/go-gitlab-client.v1"
	helper "github.com/rande/gitlab-ci-helper"
)

type CiFlowdockStatusCommand struct {
	Ui           cli.Ui
	Verbose      bool
	Last         bool
	BuildRef     string
	BuildName    string
	BuildRefName string
	Project      string
}

func (c *CiFlowdockStatusCommand) Run(args []string) int {

	config := &FlowdockConfig{
		Server: "https://api.flowdock.com",
	}

	levels := map[string]int{
		"pending":  BUILD_PENDING,
		"running":  BUILD_RUNNING,
		"success":  BUILD_SUCCESS,
		"failed":   BUILD_FAILED,
		"canceled": BUILD_CANCELLED,
	}

	// flowdock colors: red, green, yellow, cyan, orange, grey, black, lime, purple, blue
	colors := map[int]string{
		BUILD_PENDING:   "black",
		BUILD_RUNNING:   "orange",
		BUILD_SUCCESS:   "green",
		BUILD_FAILED:    "red",
		BUILD_CANCELLED: "grey",
	}

	cmdFlags := flag.NewFlagSet("flowdock:thread", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	cmdFlags.BoolVar(&c.Verbose, "verbose", false, "")
	cmdFlags.BoolVar(&c.Last, "last", false, "")
	cmdFlags.StringVar(&config.Token, "token", os.Getenv("FLOWDOCK_SOURCE_TOKEN"), "The room's token (default: env var FLOWDOCK_SOURCE_TOKEN)")
	cmdFlags.StringVar(&c.BuildRef, "ref", helper.GetEnv("CI_COMMIT_SHA", os.Getenv("CI_BUILD_REF")), "The commit related to the build (default: env var CI_COMMIT_SHA/CI_BUILD_REF)")
	cmdFlags.StringVar(&c.Project, "project", os.Getenv("CI_PROJECT_ID"), "The project related to the build (default: env var CI_PROJECT_ID)")
	cmdFlags.StringVar(&c.BuildName, "name", helper.GetEnv("CI_JOB_NAME", os.Getenv("CI_BUILD_NAME")), "The build's name (default: env var CI_BUILD_NAME)")
	cmdFlags.StringVar(&c.BuildRefName, "ref-name", helper.GetEnv("CI_COMMIT_REF_NAME", os.Getenv("CI_BUILD_REF_NAME")), "The reference name (default: env var CI_COMMIT_REF_NAME/CI_BUILD_REF_NAME)")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	args = cmdFlags.Args()

	if len(args) != 2 {
		c.Ui.Error("Invalid number of arguments\n")
		cmdFlags.Usage()

		return 1
	}

	config.Organization = args[0]
	config.Flow = args[1]

	gitlabConfig := helper.NewConfig()
	gitLabClient := gitlab.NewGitlab(gitlabConfig.Gitlab.Host, gitlabConfig.Gitlab.ApiPath, gitlabConfig.Gitlab.Token)

	project, err := helper.GetProject(c.Project, gitLabClient)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to fetch the project: %s", err.Error()))

		return 1
	}

	c.Ui.Output(fmt.Sprintf("Found project: %s/%s (id: %d)", project.Namespace.Name, project.Name, project.Id))

	builds, err := gitLabClient.ProjectCommitBuilds(strconv.FormatInt(int64(project.Id), 10), c.BuildRef)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error: %s", err.Error()))

		return 1
	}

	var level = 0
	var status = "-"

	activeBuilds := make(map[string]*gitlab.Build, 0)

	// filter builds to only get active ones
	for _, b := range builds {
		if ib, ok := activeBuilds[b.Name]; !ok {
			activeBuilds[b.Name] = b
		} else {
			if ib.Id > b.Id {
				activeBuilds[b.Name] = ib
			}
		}
	}

	// found the valid status build
	for _, b := range activeBuilds {
		if levels[b.Status] > level {
			level = levels[b.Status]
			status = b.Status
		}
	}

	// if it the last pass with no previous error, we force the final job to be green
	// the command should be the last one (reporting ...)
	if c.Last && level == BUILD_RUNNING {
		status = "success"
		level = BUILD_SUCCESS

		activeBuilds[c.BuildName].Status = "success"
	}

	tpl := `<h3>{{ .Project.Namespace.Name }} / {{ .Project.Name }}</h3>

    <ul>
        <li><b>Project:</b> <a href="{{ .Project.WebUrl }}">{{ .Project.WebUrl }}</a></li>
    </ul>

    <hr />
    <ul>
        <li><b>Author:</b> {{ .Build.Commit.Author_Name }}</li>
        <li><b>Title:</b> {{ .Build.Commit.Title }}</li>
        <li><b>Builds commit:</b> <a href="{{ .Project.WebUrl }}/commit/{{ .Build.Commit.Short_Id }}/pipelines">{{ .Project.WebUrl }}/commit/{{ .Build.Commit.Id }}/pipelines</a></li>
    </ul>

    <table class="build-status" style="width: 100%">
        <thead>
            <tr>
                <th style="padding: 4px; text-align: left;">Status</th>
                <th style="padding: 4px; text-align: left;">Name</th>
                <th style="padding: 4px; text-align: left;">Stage</th>
            </tr>
        </thead>
        <tbody>
            {{ range .Builds }}
                <tr>
                    <td style="padding: 4px;">
                        {{ $k := index $.Levels .Status }}
                        {{ $color := index $.Colors $k }}

                        <a href="{{ $.Project.WebUrl }}/builds/{{ .Id }}" style="color: {{ $color }}">&#10026;</style> #{{ .Id }} - {{ .Status }}</a>

                        {{ if gt .ArtifactsFile.Size 0 }}
                             <a href="{{ $.Project.WebUrl }}/builds/{{ .Id }}/artifacts/browse">&#128230;</a>
                        {{ end }}
                    </td>

                    <td style="padding: 4px;">{{ .Name }}</td>
                    <td style="padding: 4px;">{{ .Stage }}</td>
                </tr>
            {{ end }}
        </tbody>
    </table>
`

	t, err := template.New("foo").Parse(tpl)

	body := bytes.NewBuffer([]byte(""))

	err = t.Execute(body, map[string]interface{}{
		"Builds":  activeBuilds,
		"Build":   builds[0],
		"Project": project,
		"Levels":  levels,
		"Colors":  colors,
	})

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error: %s", err.Error()))

		return 1
	}

	message := &FlowdockMessage{
		Event:      "activity",
		ExternalId: fmt.Sprintf("gitlab:%s:%s", c.BuildRef, c.BuildRefName),
		Title:      fmt.Sprintf("Update status, job:%s", c.BuildName),
		Thread: &FlowdockThread{
			Status: &FlowdockThreadStatus{
				Color: colors[level],
				Value: status,
			},
			Title:       fmt.Sprintf("Jobs for %s - %s", project.Name, c.BuildRefName),
			Body:        body.String(),
			ExternalUrl: fmt.Sprintf("%s/commit/%s/pipelines", project.WebUrl, c.BuildRef),
		},
		Author: &FlowdockAuthor{
			Name:   "GitlabCi",
			Avatar: "https://about.gitlab.com/images/gitlab_logo.png",
		},
	}

	message.Token = config.Token
	message.Flow = config.Flow

	buf := bytes.NewBuffer([]byte(""))
	e := json.NewEncoder(buf)

	e.Encode(message)

	client := &http.Client{}

	r, _ := http.NewRequest("POST", fmt.Sprintf("%s/flows/%s/%s/messages", config.Server, config.Organization, config.Flow), buf)
	r.Header.Add("Content-Type", "application/json")

	_, err = client.Do(r)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Err: %s", err.Error()))

		return 1
	}

	return 0
}

func (c *CiFlowdockStatusCommand) Synopsis() string {
	return "Send a message to one Flowdock's room."
}

func (c *CiFlowdockStatusCommand) Help() string {
	helpText := fmt.Sprintf(`
Usage: gitlab-ci-helper flowdock:message [options] organisation flow

  Build a flowdock thread from the current build. Information are retrieved from
  environment variables.

  The external thread id is: gitlab:sha1

  You can use the -last option to indicate that the current job is the last one.

Arguments:
  organisation        The organisation name
  flow                The flow reference

Options:
  -ref                The commit related to the build (default:
                        9.x: CI_COMMIT_SHA or 8.x: CI_BUILD_REF)
  -project            The project related to the build (default: env var CI_PROJECT_ID)
  -name               The build's name (default:
                        9.x: CI_JOB_NAME or 8.x: CI_BUILD_NAME)
  -ref-name           The reference name (default:
                        9.x: CI_COMMIT_REF_NAME or 8.x: CI_BUILD_REF_NAME)
  -last               Indicate if the current build is the last one
  -token              The flow's token (default: env var FLOWDOCK_SOURCE_TOKEN)
  -verbose            Add verbose information to the output

%s

Gitlab's credentials are retrieved from environment:

  GITLAB_HOST         The gitlab host
  GITLAB_TOKEN        The user's token
  GITLAB_API_PATH     (optional) the api path, default to: "/api/v4"

`, flowdockConfiguration)

	return strings.TrimSpace(helpText)
}
