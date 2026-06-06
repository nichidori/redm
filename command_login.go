package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/nichidori/redm/internal/config"
)

var CommandLogin = Command{
	Name:        "login",
	Description: "Save Redmine credentials",
	Setup: func(fs *flag.FlagSet, s *state) func() error {
		return func() error {
			if s.config != nil {
				return fmt.Errorf("you are already logged in")
			}

			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Redmine URL: ")
			url, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("read URL: %w", err)
			}
			url = strings.TrimSpace(url)
			if url == "" {
				return fmt.Errorf("URL must not be empty")
			}

			fmt.Print("API Key: ")
			rawKey, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				return fmt.Errorf("read API key: %w", err)
			}
			apiKey := strings.TrimSpace(string(rawKey))
			if apiKey == "" {
				return fmt.Errorf("API key must not be empty")
			}

			cfg := &config.Config{URL: url, APIKey: apiKey}
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			s.config = cfg
			fmt.Println("Login successful.")
			return nil
		}
	},
}
