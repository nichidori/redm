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

			fmt.Printf("ID:        %d\n", user.ID)
			fmt.Printf("Login:     %s\n", user.Login)
			fmt.Printf("Name:      %s %s\n", user.FirstName, user.LastName)
			fmt.Printf("Email:     %s\n", user.Mail)
			return nil
		}
	},
}
