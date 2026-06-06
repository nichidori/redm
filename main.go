package main

import (
	"fmt"
	"os"

	"github.com/nichidori/redm/internal/config"
	"github.com/nichidori/redm/internal/redmineapi"
)

func main() {
	var baseURL, apiKey string
	var cfg *config.Config

	if c, err := config.Load(); err == nil {
		cfg = c
		baseURL = cfg.URL
		apiKey = cfg.APIKey
	}

	c := redmineapi.NewClient(baseURL, apiKey)

	cli := &CLI{
		Name: "redm",
		State: &state{
			client: c,
			config: cfg,
		},
	}

	cli.Register(CommandProject)
	cli.Register(CommandIssue)
	cli.Register(CommandLogin)
	cli.Register(CommandLogout)

	if err := cli.Execute(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
