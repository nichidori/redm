package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	"github.com/nichidori/redm/internal/redmineapi"
)

var CommandIssue = Command{
	Name:        "issue",
	Description: "Lists all issues",
	Setup: func(fs *flag.FlagSet, s *state) func() error {
		return func() error {
			if s.config == nil {
				return fmt.Errorf("not logged in")
			}

			resp, err := s.client.GetIssues(context.Background(), redmineapi.IssueFilter{})
			if err != nil {
				return fmt.Errorf("failed to fetch issues: %w", err)
			}

			if len(resp.Issues) == 0 {
				fmt.Println("No issues found.")
				return nil
			}

			issues := resp.Issues
			cols := []Column{
				NewColumn("ID", 4, func(i int) string { return strconv.Itoa(issues[i].ID) }),
				NewColumn("Project", 16, func(i int) string { return issues[i].Project.Name }),
				NewColumn("Subject", 32, func(i int) string { return issues[i].Subject }),
				NewColumn("Priority", 12, func(i int) string { return issues[i].Priority.Name }),
				NewColumn("Assignee", 16, func(i int) string { return issues[i].AssignedTo.Name }),
				NewColumn("Status", 12, func(i int) string { return issues[i].Status.Name }),
				NewColumn("Progress", 8, func(i int) string { return strconv.Itoa(issues[i].DoneRatio) + "%" }),
			}
			PrintTable(cols, len(issues))

			return nil
		}
	},
}
