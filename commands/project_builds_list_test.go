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

func Test_Project_BuildsList_From_Args(t *testing.T) {
	fpBuild, err := os.Open("../fixtures/builds.json")
	assert.NoError(t, err)

	fpProject, err := os.Open("../fixtures/project.json")
	assert.NoError(t, err)

	reqs := []*helper.FakeRequest{
		{
			Path:   "/api/v3/projects/3",
			Method: "GET",
			Response: &http.Response{
				Body: fpProject,
			},
		},
		{
			Path:   "/api/v3/projects/3/builds",
			Method: "GET",
			Response: &http.Response{
				Body: fpBuild,
			},
		},
	}

	envs := map[string]string{}

	helper.WrapperTestCommand(reqs, envs, t, func(ts *httptest.Server) {
		ui := &cli.MockUi{}
		c := &ProjectBuildsListCommand{
			Ui: ui,
		}

		code := c.Run([]string{"3"})

		assert.Equal(t, 0, code)
	})
}

func Test_Project_BuildsList_Help(t *testing.T) {
	c := &ProjectBuildsListCommand{
		Ui: &cli.MockUi{},
	}

	assert.True(t, len(c.Help()) > 0)
	assert.True(t, len(c.Synopsis()) > 0)
}

func Test_Project_BuildsList_InvalidRun(t *testing.T) {
	c := &ProjectBuildsListCommand{
		Ui: &cli.MockUi{},
	}

	assert.Equal(t, 1, c.Run([]string{"--foobar"}))
}
