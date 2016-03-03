// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"github.com/mitchellh/cli"
	"sort"
	"strings"
)

type DumpReadmeCommand struct {
	Ui       cli.Ui
	Verbose  bool
	Commands map[string]cli.CommandFactory
}

func (c *DumpReadmeCommand) Run(args []string) int {

	mk := make([]string, len(c.Commands))

	i := 0
	for k := range c.Commands {
		mk[i] = k
		i++
	}

	sort.Strings(mk)

	c.Ui.Output("## Commands")

	for _, name := range mk {

		if name == "dump:readme" {
			continue
		}

		cmd, _ := c.Commands[name]()

		c.Ui.Output(fmt.Sprintf("\n### %s\n", name))

		for _, l := range strings.Split(cmd.Help(), "\n") {
			c.Ui.Output(fmt.Sprintf("    %s", l))
		}
	}

	return 0
}

func (c *DumpReadmeCommand) Synopsis() string {
	return "Dump the command readme."
}

func (c *DumpReadmeCommand) Help() string {
	return strings.TrimSpace("")
}
