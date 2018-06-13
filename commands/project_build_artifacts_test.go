// Copyright © 2016-present Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mitchellh/cli"
	helper "github.com/rande/gitlab-ci-helper"
	"github.com/stretchr/testify/assert"
)

func Test_Project_Builds_Artifacts(t *testing.T) {
	fpBuilds, err := os.Open("../fixtures/builds.json")
	assert.NoError(t, err)

	fpProjects, err := os.Open("../fixtures/projects.json")
	assert.NoError(t, err)

	fpProject, err := os.Open("../fixtures/project.json")
	assert.NoError(t, err)

	fpCommits, err := os.Open("../fixtures/commits.json")
	assert.NoError(t, err)

	fpArchive, err := os.Open("../fixtures/artifacts.zip")
	assert.NoError(t, err)

	headers := http.Header{
		"Content-Type": []string{"application/json"},
	}

	reqs := []*helper.FakeRequest{
		{
			Path:   "/api/v3/projects/3",
			Method: "GET",
			Response: &http.Response{
				Body:   fpProject,
				Header: headers,
			},
		},
		{
			Path:   "/api/v3/projects",
			Method: "GET",
			Response: &http.Response{
				Body:   fpProjects,
				Header: headers,
			},
		},
		{
			Path:   "/api/v3/projects/3/builds",
			Method: "GET",
			Response: &http.Response{
				Body:   fpBuilds,
				Header: headers,
			},
		},
		{
			Path:   "/api/v3/projects/3/repository/commits/889935cf4d3e7558ae6c0d4dd62e20ea600f5a57/builds",
			Method: "GET",
			Response: &http.Response{
				Body:   fpCommits,
				Header: headers,
			},
		},
		{
			Path:   "/api/v3/projects/3/builds/69/artifacts",
			Method: "GET",
			Response: &http.Response{
				Body: fpArchive,
				Header: http.Header{
					"Content-Type": []string{"application/zip"},
				},
			},
		},
	}

	envs := map[string]string{}

	helper.WrapperTestCommand(reqs, envs, t, func(ts *httptest.Server) {
		ui := &cli.MockUi{}
		c := &ProjectBuildArtifactCommand{
			Ui: ui,
		}

		code := c.Run([]string{"-project", "3", "-ref", "889935cf4d3e7558ae6c0d4dd62e20ea600f5a57", "-job", "rubocop"})

		assert.Equal(t, 0, code)

		expected := "Found project: Diaspora/Diaspora Project Site (id: 3)\nFound build - stage:test status:canceled id:69\nDownloading artifacts... (artifacts.zip)\nDone!\n"
		assert.Equal(t, expected, ui.OutputWriter.String())
		assert.Equal(t, "", ui.ErrorWriter.String())

		os.Remove(c.ArtifactsFile)
	})

}

func Test_Project_Builds_Artifacts_Help(t *testing.T) {
	c := &ProjectBuildArtifactCommand{
		Ui: &cli.MockUi{},
	}

	assert.True(t, len(c.Help()) > 0)
	assert.True(t, len(c.Synopsis()) > 0)
}

func Test_Project_Builds_Artifacts_InvalidRun(t *testing.T) {
	c := &ProjectBuildArtifactCommand{
		Ui: &cli.MockUi{},
	}

	assert.Equal(t, 1, c.Run([]string{"--foobar"}))
}
