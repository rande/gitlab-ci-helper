// Copyright © 2016-present Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package flowdock

import (
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

func Test_CiFlowdockStatusCommand_Help(t *testing.T) {
	c := &CiFlowdockStatusCommand{
		Ui: &cli.MockUi{},
	}

	assert.True(t, len(c.Help()) > 0)
	assert.True(t, len(c.Synopsis()) > 0)
}

func Test_CiFlowdockStatusCommand_InvalidRun(t *testing.T) {
	c := &CiFlowdockStatusCommand{
		Ui: &cli.MockUi{},
	}

	assert.Equal(t, 1, c.Run([]string{"--foobar"}))
}
