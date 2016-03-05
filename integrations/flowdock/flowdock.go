// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package flowdock

const (
	BUILD_PENDING   = 0
	BUILD_SUCCESS   = 1
	BUILD_RUNNING   = 2
	BUILD_FAILED    = 3
	BUILD_CANCELLED = 4
)

var flowdockConfiguration = `
Configuration:

  Please note, the command use the new Flowdock API.

  1. Go to https://www.flowdock.com/oauth/applications
  2. Create a new Application
  3. Enter name, description and make sure "Short application" is selected.
  4. Once validated, Go to "Tools for testing" and create a new source.
  5. Press "Generate Source" and store the generated token for later use
     as the FLOWDOCK_SOURCE_TOKEN
`

type FlowdockConfig struct {
	Organization string
	Flow         string
	Server       string
	Token        string
}

type FlowdockThreadStatus struct {
	Color string `json:"color"`
	Value string `json:"value"`
}

type FlowdockAuthor struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
}

type FlowdockField struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type FlowdockThread struct {
	Title       string                `json:"title"`
	Body        string                `json:"body"`
	ExternalUrl string                `json:"external_url"`
	Status      *FlowdockThreadStatus `json:"status"`
	Fields      []*FlowdockField      `json:"fields"`
}

type FlowdockMessage struct {
	Event      string          `json:"event"`
	Content    string          `json:"content,omitempty"`
	Message    string          `json:"message,omitempty"`
	Title      string          `json:"title,omitempty"`
	Username   string          `json:"external_user_name,omitempty"`
	ExternalId string          `json:"external_thread_id,omitempty"`
	Token      string          `json:"flow_token,omitempty"`
	Flow       string          `json:"flow,omitempty"`
	Parent     string          `json:"parent,omitempty"`
	Uuid       string          `json:"uuid,omitempty"`
	ThreadId   string          `json:"thread_id,omitempty"`
	Thread     *FlowdockThread `json:"thread,omitempty"`
	Author     *FlowdockAuthor `json:"author,omitempty"`
}
