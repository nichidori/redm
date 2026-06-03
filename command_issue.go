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
			resp, err := s.client.GetIssues(context.Background(), redmineapi.IssueFilter{})
			if err != nil {
				return fmt.Errorf("failed to fetch issues: %w", err)
			}

			if len(resp.Issues) == 0 {
				fmt.Println("No issues found.")
				return nil
			}

			fmt.Print(FixLength("ID", 4))
			fmt.Print(" | ")
			fmt.Print(FixLength("Project", 12))
			fmt.Print(" | ")
			fmt.Print(FixLength("Subject", 32))
			fmt.Print(" | ")
			fmt.Print(FixLength("Priority", 12))
			fmt.Print(" | ")
			fmt.Print(FixLength("Assignee", 16))
			fmt.Print(" | ")
			fmt.Print(FixLength("Status", 12))
			fmt.Print(" | ")
			fmt.Print(FixLength("Progress", 8))
			fmt.Println()

			for _, p := range resp.Issues {
				fmt.Print(FixLength(strconv.Itoa(p.ID), 4))
				fmt.Print(" | ")
				fmt.Print(FixLength(p.Project.Name, 12))
				fmt.Print(" | ")
				fmt.Print(FixLength(p.Subject, 32))
				fmt.Print(" | ")
				fmt.Print(FixLength(p.Priority.Name, 12))
				fmt.Print(" | ")
				fmt.Print(FixLength(p.AssignedTo.Name, 16))
				fmt.Print(" | ")
				fmt.Print(FixLength(p.Status.Name, 12))
				fmt.Print(" | ")
				fmt.Print(FixLength(strconv.Itoa(p.DoneRatio)+"%", 8))
				fmt.Println()
			}

			return nil
		}
	},
}
