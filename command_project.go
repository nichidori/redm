package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
)

var CommandProject = Command{
	Name:        "project",
	Description: "Lists all projects",
	Setup: func(fs *flag.FlagSet, s *state) func() error {
		return func() error {
			resp, err := s.client.GetProjects(context.Background(), 0, 25)
			if err != nil {
				return fmt.Errorf("failed to fetch projects: %w", err)
			}

			fmt.Print(FixLength("ID", 4))
			fmt.Print(" | ")
			fmt.Print(FixLength("Name", 10))
			fmt.Print(" | ")
			fmt.Print(FixLength("Description", 30))
			fmt.Print(" | ")
			fmt.Print(FixLength("Last Update", 30))
			fmt.Println()

			for _, p := range resp.Projects {
				fmt.Print(FixLength(strconv.Itoa(p.ID), 4))
				fmt.Print(" | ")
				fmt.Print(FixLength(p.Name, 10))
				fmt.Print(" | ")
				fmt.Print(FixLength(p.Description, 30))
				fmt.Print(" | ")
				fmt.Print(FixLength(p.UpdatedOn.String(), 30))
				fmt.Println()
			}

			return nil
		}
	},
}
