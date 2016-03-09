// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/cli"
	helper "github.com/rande/gitlab-ci-helper"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"os"
	"testing"
)

func Test_Ci_Dump_Meta(t *testing.T) {
	ui := &cli.MockUi{}
	c := &CiDumpMetaCommand{
		Ui: ui,
	}

	path := fmt.Sprintf("%s/gitlab_helper.ci.json", os.TempDir())
	code := c.Run([]string{"-file", path})

	assert.Equal(t, 0, code)

	_, err := os.Stat(path)

	assert.NoError(t, err)
	assert.False(t, os.IsNotExist(err))

	os.Remove(path)
}

func Test_Ci_Dump_Meta_With_Values(t *testing.T) {
	path := fmt.Sprintf("%s/gitlab_helper.ci.json", os.TempDir())

	reqs := []*helper.FakeRequest{}
	envs := map[string]string{
		"CI_BUILD_ID":        "CI_BUILD_ID",
		"CI_BUILD_REF":       "CI_BUILD_REF",
		"CI_BUILD_REF_NAME":  "CI_BUILD_REF_NAME",
		"CI_BUILD_TAG":       "CI_BUILD_TAG",
		"CI_BUILD_STAGE":     "CI_BUILD_STAGE",
		"CI_BUILD_NAME":      "CI_BUILD_NAME",
		"CI_PROJECT_ID":      "CI_PROJECT_ID",
		"CI_PROJECT_DIR":     "CI_PROJECT_DIR",
		"CI_SERVER_NAME":     "CI_SERVER_NAME",
		"CI_SERVER_REVISION": "CI_SERVER_REVISION",
		"CI_SERVER_VERSION":  "CI_SERVER_VERSION",
	}

	meta := &Meta{
		Build: &MetaBuild{
			Id:      "CI_BUILD_ID",
			Ref:     "CI_BUILD_REF",
			RefName: "CI_BUILD_REF_NAME",
			Tag:     "CI_BUILD_TAG",
			Stage:   "CI_BUILD_STAGE",
			JobName: "CI_BUILD_NAME",
		},
		Project: &MetaProject{
			Id:  "CI_PROJECT_ID",
			Dir: "CI_PROJECT_DIR",
		},
		Server: &MetaServer{
			Name:     "CI_SERVER_NAME",
			Revision: "CI_SERVER_REVISION",
			Version:  "CI_SERVER_VERSION",
		},
	}

	helper.WrapperTestCommand(reqs, envs, t, func(ts *httptest.Server) {
		ui := &cli.MockUi{}
		c := &CiDumpMetaCommand{
			Ui: ui,
		}

		code := c.Run([]string{"-file", path})

		assert.Equal(t, 0, code)

		r, err := os.Open(path)

		assert.NoError(t, err)
		defer r.Close()

		m := &Meta{}
		err = json.NewDecoder(r).Decode(m)

		assert.NoError(t, err)

		assert.Equal(t, meta, m)
	})

	os.Remove(path)
}
