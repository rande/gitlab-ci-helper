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

func Test_Project_List(t *testing.T) {
	fp, err := os.Open("../fixtures/projects.json")

	assert.NoError(t, err)

	reqs := []*helper.FakeRequest{
		{
			Path:   "/api/v4/projects",
			Method: "GET",
			Response: &http.Response{
				Body: fp,
			},
		},
	}

	envs := map[string]string{
		"GITLAB_TOKEN": "THE_SECRET_GITLAB_TOKEN",
	}

	helper.WrapperTestCommand(reqs, envs, t, func(ts *httptest.Server) {
		ui := &cli.MockUi{}
		c := &ProjectsListCommand{
			Ui: ui,
		}

		code := c.Run(nil)

		assert.Equal(t, 0, code)

		expected := "Trying to find project from options\n >    3 - Diaspora Client - Diaspora\n >    6 - Puppet - Brightbox\n"
		assert.Equal(t, expected, ui.OutputWriter.String())
	})
}

func Test_Project_List_Help(t *testing.T) {
	c := &ProjectsListCommand{
		Ui: &cli.MockUi{},
	}

	assert.True(t, len(c.Help()) > 0)
	assert.True(t, len(c.Synopsis()) > 0)
}

func Test_Project_List_InvalidRun(t *testing.T) {
	c := &ProjectsListCommand{
		Ui: &cli.MockUi{},
	}

	assert.Equal(t, 1, c.Run([]string{"--foobar"}))
}
