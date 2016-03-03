// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gitlab_ci_helper

import "os"

type Config struct {
	Host    string `json:"host"`
	Token   string `json:"token"`
	ApiPath string `json:"api_path"`
}

func NewConfig() *Config {
	c := &Config{
		Host:    os.Getenv("GITLAB_HOST"),
		Token:   os.Getenv("GITLAB_TOKEN"),
		ApiPath: os.Getenv("GITLAB_API_PATH"),
	}

	if c.ApiPath == "" {
		c.ApiPath = "/api/v3"
	}

	return c
}
