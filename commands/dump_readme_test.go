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

func Test_Dump_Readme(t *testing.T) {
	ui := &cli.MockUi{}
	c := &DumpReadmeCommand{
		Ui: ui,
	}

	code := c.Run(nil)

	assert.Equal(t, 0, code)
}

func Test_Dump_Readme_Help(t *testing.T) {
	c := &DumpReadmeCommand{
		Ui: &cli.MockUi{},
	}

	assert.True(t, len(c.Help()) > 0)
	assert.True(t, len(c.Synopsis()) > 0)
}

func Test_Dump_Readme_InvalidRun(t *testing.T) {
	c := &DumpReadmeCommand{
		Ui: &cli.MockUi{},
	}

	assert.Equal(t, 1, c.Run([]string{"--foobar"}))
}
