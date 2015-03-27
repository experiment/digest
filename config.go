package main

import (
	"github.com/ryanlower/setting"
)

type config struct {
	Repo struct {
		Owner string `env:"REPO_OWNER"`
		Name  string `env:"REPO_NAME"`
	}
	Github struct {
		AccessToken string `env:"GITHUB_ACCESS_TOKEN"`
	}
}

// Load config from environment
func (c *config) load() {
	setting.Load(c)
}
