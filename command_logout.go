package main

import (
	"flag"
	"fmt"

	"github.com/nichidori/redm/internal/config"
)

var CommandLogout = Command{
	Name:        "logout",
	Description: "Remove saved Redmine credentials",
	Setup: func(fs *flag.FlagSet, s *state) func() error {
		return func() error {
			if s.config == nil {
				return fmt.Errorf("not logged in")
			}

			if err := config.Delete(); err != nil {
				return fmt.Errorf("logout failed: %w", err)
			}

			s.config = nil
			fmt.Println("Logged out.")
			return nil
		}
	},
}
