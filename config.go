// Copyright © 2016-present Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gitlab_ci_helper

import "os"

type GitLabConfig struct {
	Host    string `json:"host"`
	Token   string `json:"token"`
	ApiPath string `json:"api_path"`
}

type MailerConfig struct {
	SubjectPrefix string   `json:"subject"`
	Sender        string   `json:"sender"`
	Dest          []string `json:"dest"`
	Host          string   `json:"host"`
	Username      string   `json:"username"`
	Password      string   `json:"password"`
}

type Config struct {
	Gitlab *GitLabConfig `json:"gitlab"`
}

func NewConfig() *Config {
	gitlab := &GitLabConfig{
		Host:    os.Getenv("GITLAB_HOST"),
		Token:   os.Getenv("GITLAB_TOKEN"),
		ApiPath: os.Getenv("GITLAB_API_PATH"),
	}

	if gitlab.ApiPath == "" {
		gitlab.ApiPath = "/api/v4"
	}

	return &Config{
		Gitlab: gitlab,
	}
}
