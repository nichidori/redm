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

			projects := resp.Projects
			cols := []Column{
				NewColumn("ID", 4, func(i int) string { return strconv.Itoa(projects[i].ID) }),
				NewColumn("Name", 12, func(i int) string { return projects[i].Name }),
				NewColumn("Description", 32, func(i int) string { return projects[i].Description }),
				NewColumn("Last Update", 32, func(i int) string { return projects[i].UpdatedOn.String() }),
			}
			PrintTable(cols, len(projects))

			return nil
		}
	},
}
