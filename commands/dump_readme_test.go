// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Dump_Readme(t *testing.T) {
	ui := &cli.MockUi{}
	c := &DumpReadmeCommand{
		Ui: ui,
	}

	code := c.Run(nil)

	assert.Equal(t, 0, code)
}
