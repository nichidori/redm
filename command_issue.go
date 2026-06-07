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
	Description: "List issues",
	Setup: func(fs *flag.FlagSet, s *state) func() error {
		projectID := fs.Int("project", 0, "Filter issues by project ID")
		assigneeID := fs.String("assignee", "", "Filter by assignee ID (or \"me\" for current user)")
		statusID := fs.String("status", "", "Filter by status ID (or \"open\", \"closed\", \"*\")")

		return func() error {
			if s.config == nil {
				return fmt.Errorf("not logged in")
			}

			resp, err := s.client.GetIssues(context.Background(), redmineapi.IssueFilter{
				ProjectID:    *projectID,
				AssignedToID: *assigneeID,
				StatusID:     *statusID,
			})
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
				NewColumn("Assignee", 12, func(i int) string { return issues[i].AssignedTo.Name }),
				NewColumn("Start Date", 12, func(i int) string { return issues[i].StartDate }),
				NewColumn("Due Date", 12, func(i int) string { return issues[i].DueDate }),
				NewColumn("Status", 12, func(i int) string { return issues[i].Status.Name }),
				NewColumn("Progress", 8, func(i int) string { return strconv.Itoa(issues[i].DoneRatio) + "%" }),
			}
			PrintTable(cols, len(issues))

			return nil
		}
	},
}
