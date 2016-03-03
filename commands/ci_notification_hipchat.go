// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	helper "github.com/rande/gitlab-ci-helper"
	"io"
	"net/http"
	"os"
	"strings"
)

type HipChatConfig struct {
	Room   string
	Server string
	Token  string
}

type HipChatDescription struct {
	Value  int    `json:"value"`
	Format string `json:"format"`
}

type HipChatCard struct {
	Style       string              `json:"style,omitempty"`
	Title       string              `json:"title,omitempty"`
	Description *HipChatDescription `json:"description,omitempty"`
	Id          string              `json:"id,omitempty"`
	Url         string              `json:"url,omitempty"`
}

type HipChatMessage struct {
	Format  string       `json:"message_format,omitempty"`
	Color   string       `json:"color,omitempty"`
	Notify  bool         `json:"notify,omitempty"`
	Message string       `json:"message"`
	Card    *HipChatCard `json:"card,omitempty"`
}

type CiNotificationHipchatCommand struct {
	Ui      cli.Ui
	Verbose bool
}

func (c *CiNotificationHipchatCommand) Run(args []string) int {

	config := &HipChatConfig{}
	message := &HipChatMessage{
	//Card: &HipChatCard{
	//	Description: &HipChatDescription {
	//		Value: 1,
	//		Format: "html",
	//	},
	//},
	}

	cmdFlags := flag.NewFlagSet("ci:notification:hipchat", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	cmdFlags.BoolVar(&c.Verbose, "verbose", false, "")
	cmdFlags.StringVar(&config.Token, "token", os.Getenv("HIPCHAT_TOKEN"), "The room's token (default: env var HIPCHAT_TOKEN)")
	cmdFlags.StringVar(&config.Server, "server", helper.GetEnv("HIPCHAT_SERVER", "https://api.hipchat.com"), "The hipchat server, default to env var HIPCHAT_SERVER, then https://api.hipchat.com")

	cmdFlags.StringVar(&message.Color, "color", "gray", "The message color (default: gray, values: yellow, green, red, purple, gray, random)")
	cmdFlags.BoolVar(&message.Notify, "notify", false, "Whether this message should trigger a user notification (default: false)")

	// @todo: support the card option, seems to be a quite complex to get it right
	//cmdFlags.StringVar(&message.Card.Style, "style", "application", "Type of the card (default: application, values: file, image, application, link, media)")
	//cmdFlags.StringVar(&message.Card.Title, "title", "Gitlab Notification", "The title of the card")
	//cmdFlags.StringVar(&message.Card.Id, "id", "", "An id that will help HipChat recognise the same card when it is sent multiple times")
	//cmdFlags.StringVar(&message.Card.Url, "url", "", "An URL")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	args = cmdFlags.Args()

	if len(args) != 2 {
		cmdFlags.Usage()

		return 1
	}

	config.Room = args[0]
	message.Message = args[1]

	buf := bytes.NewBuffer([]byte(""))
	e := json.NewEncoder(buf)

	e.Encode(message)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", fmt.Sprintf("%s/v2/room/%s/notification", config.Server, config.Room), buf)
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.Token))
	r.Header.Add("Content-Type", "application/json")

	resp, _ := client.Do(r)

	fmt.Println(resp.Status)

	io.Copy(os.Stdout, resp.Body)

	resp.Body.Close()

	return 1
}

func (c *CiNotificationHipchatCommand) Synopsis() string {
	return "Send a message to one HipChat's room."
}

func (c *CiNotificationHipchatCommand) Help() string {
	helpText := `
Usage: gitlab-helper ci:notify:hipchat [options] room message

  Send a message to one HipChat's room

Arguments:
  room                The room reference
  message             The message to send

Options:
  -token              The room's token (default: env var HIPCHAT_TOKEN)
  -server             The hipchat server, default to env var HIPCHAT_SERVER, then https://api.hipchat.com
  -verbose            Add verbose information to the output
`
	return strings.TrimSpace(helpText)
}
