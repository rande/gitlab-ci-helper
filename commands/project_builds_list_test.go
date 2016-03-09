// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"github.com/mitchellh/cli"
	helper "github.com/rande/gitlab-ci-helper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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
