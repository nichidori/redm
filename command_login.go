package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/nichidori/redm/internal/config"
	"github.com/nichidori/redm/internal/redmineapi"
)

var CommandLogin = Command{
	Name:        "login",
	Description: "Set Redmine authentication",
	Setup: func(fs *flag.FlagSet, s *state) func() error {
		return func() error {
			if s.config != nil {
				return fmt.Errorf("you are already logged in")
			}

			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Redmine URL: ")
			rawURL, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("read URL: %w", err)
			}
			rawURL = strings.TrimSpace(rawURL)
			if rawURL == "" {
				return fmt.Errorf("URL must not be empty")
			}
			parsedURL, err := url.ParseRequestURI(rawURL)
			if err != nil {
				return fmt.Errorf("invalid URL: %w", err)
			}
			if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
				return fmt.Errorf("URL must start with http:// or https://")
			}
			rawURL = strings.TrimRight(rawURL, "/")

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

			testClient := redmineapi.NewClient(rawURL, apiKey)
			user, err := testClient.GetCurrentUser(context.Background())
			if err != nil {
				errStr := err.Error()
				if strings.Contains(errStr, "401") || strings.Contains(errStr, "403") {
					return fmt.Errorf("login failed: invalid URL or API key")
				}
				return fmt.Errorf("failed to connect: %w", err)
			}

			cfg := &config.Config{URL: rawURL, APIKey: apiKey}
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			s.config = cfg
			fmt.Printf("Logged in as %s %s (%s)\n", user.FirstName, user.LastName, user.Login)
			return nil
		}
	},
}
