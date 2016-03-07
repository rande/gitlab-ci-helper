// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package flowdock

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	"net/http"
	"os"
	"strings"
)

type CiFlowdockMessageCommand struct {
	Ui      cli.Ui
	Verbose bool
}

func (c *CiFlowdockMessageCommand) Run(args []string) int {

	config := &FlowdockConfig{
		Server: "https://api.flowdock.com",
	}

	message := &FlowdockMessage{}

	cmdFlags := flag.NewFlagSet("ci:notification:flowdock", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	cmdFlags.BoolVar(&c.Verbose, "verbose", false, "")
	cmdFlags.StringVar(&config.Token, "token", os.Getenv("FLOWDOCK_SOURCE_TOKEN"), "The room's token (default: env var FLOWDOCK_SOURCE_TOKEN)")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	args = cmdFlags.Args()

	if len(args) != 3 {
		c.Ui.Error("Invalid number of arguments\n")
		cmdFlags.Usage()

		return 1
	}

	config.Organization = args[0]
	config.Flow = args[1]

	message.Event = "message"
	message.Content = args[2]
	message.Token = config.Token
	message.Flow = config.Flow

	buf := bytes.NewBuffer([]byte(""))
	e := json.NewEncoder(buf)

	e.Encode(message)

	client := &http.Client{}

	r, _ := http.NewRequest("POST", fmt.Sprintf("%s/flows/%s/%s/messages", config.Server, config.Organization, config.Flow), buf)
	r.Header.Add("Content-Type", "application/json")

	_, err := client.Do(r)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Err: %s", err.Error()))

		return 1
	}

	return 0
}

func (c *CiFlowdockMessageCommand) Synopsis() string {
	return "Send a message to one Flowdock's room."
}

func (c *CiFlowdockMessageCommand) Help() string {
	helpText := fmt.Sprintf(`
Usage: gitlab-ci-helper flowdock:message [options] organisation flow message

  Build a flowdock thread from the current build. Information are retrieved from
  environment variables.

Arguments:
  organisation        The organisation name
  flow                The flow reference
  message             The message to send

Options:
  -token              The flow's token (default: env var FLOWDOCK_SOURCE_TOKEN)
  -verbose            Add verbose information to the output

%s

`, flowdockConfiguration)

	return strings.TrimSpace(helpText)
}
