// Copyright © 2016-present Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

func Test_Version(t *testing.T) {
	c := &VersionCommand{
		Ui: &cli.MockUi{},
	}

	code := c.Run(nil)

	assert.Equal(t, 0, code)
}

func Test_Version_Extended(t *testing.T) {
	ui := &cli.MockUi{}
	c := &VersionCommand{
		Ui:      ui,
		RefLog:  "sha1",
		Version: "1.0.0-TEST",
	}

	code := c.Run([]string{"-e"})

	assert.Equal(t, "1.0.0-TEST - sha1\n", ui.OutputWriter.String())
	assert.Equal(t, 0, code)
}

func Test_Version_Help(t *testing.T) {
	c := &VersionCommand{
		Ui: &cli.MockUi{},
	}

	assert.True(t, len(c.Help()) > 0)
	assert.True(t, len(c.Synopsis()) > 0)
}

func Test_Version_InvalidRun(t *testing.T) {
	c := &VersionCommand{
		Ui: &cli.MockUi{},
	}

	assert.Equal(t, 1, c.Run([]string{"--foobar"}))
}
