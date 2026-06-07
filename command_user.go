package main

import (
	"context"
	"flag"
	"fmt"
)

var CommandUser = Command{
	Name:        "user",
	Description: "Show current logged in user",
	Setup: func(fs *flag.FlagSet, s *state) func() error {
		return func() error {
			if s.config == nil {
				return fmt.Errorf("not logged in")
			}

			user, err := s.client.GetCurrentUser(context.Background())
			if err != nil {
				return fmt.Errorf("failed to fetch current user: %w", err)
			}

			fmt.Printf("%s (%s %s)\n", user.Login, user.FirstName, user.LastName)
			return nil
		}
	},
}
