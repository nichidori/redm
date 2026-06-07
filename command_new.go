package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nichidori/redm/internal/redmineapi"
)

var CommandNew = Command{
	Name:        "new",
	Description: "Create new issue",
	Setup: func(fs *flag.FlagSet, s *state) func() error {
		return func() error {
			if s.config == nil {
				return fmt.Errorf("not logged in")
			}

			// Fetch project list
			projectResp, err := s.client.GetProjects(context.Background(), 0, 25)
			if err != nil {
				return fmt.Errorf("failed to fetch projects: %w", err)
			}

			reader := bufio.NewReader(os.Stdin)

			// Prompt user to select project
			p, err := selectOption(
				reader,
				"project",
				projectResp.Projects,
				func(p redmineapi.Project) string { return p.Name },
			)
			if err != nil {
				return err
			}
			fmt.Printf("Project '%s' has been selected\n", p.Name)
			fmt.Println()

			// Fetch tracker list
			trackersResp, err := s.client.GetTrackers(context.Background())
			if err != nil {
				return fmt.Errorf("failed to fetch trackers: %w", err)
			}

			// Prompt user to select tracker
			t, err := selectOption(
				reader,
				"tracker",
				trackersResp.Trackers,
				func(t redmineapi.Tracker) string { return t.Name },
			)
			if err != nil {
				return err
			}
			fmt.Printf("Tracker '%s' has been selected\n", t.Name)
			fmt.Println()

			// Input issue name
			fmt.Print("Subject: ")
			subject, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			subject = strings.TrimSpace(subject)
			fmt.Println()

			// Input issue description
			fmt.Print("Description: ")
			description, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			description = strings.TrimSpace(description)
			fmt.Println()

			// Fetch current user
			user, err := s.client.GetCurrentUser(context.Background())
			if err != nil {
				return fmt.Errorf("failed to get current user: %w", err)
			}

			req := &redmineapi.CreateIssueRequest{
				ProjectID:    p.ID,
				TrackerID:    t.ID,
				Subject:      subject,
				Description:  description,
				AssignedToID: user.ID,
			}

			issue, err := s.client.CreateIssue(context.Background(), req)
			if err != nil {
				return fmt.Errorf("failed to create issue: %w", err)
			}

			fmt.Printf("Issue #%d created: %s\n", issue.ID, issue.Subject)
			return nil
		}
	},
}

func selectOption[T any](reader *bufio.Reader, label string, options []T, optionFormatter func(T) string) (T, error) {
	fmt.Printf("Select %s:\n", label)
	for i, o := range options {
		fmt.Printf("(%v) %s\n", i+1, optionFormatter(o))
	}

	fmt.Printf("Enter %s number: ", label)

	input, err := reader.ReadString('\n')
	if err != nil {
		var t T
		return t, err
	}

	idx, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		var t T
		return t, err
	}
	idx--

	if idx < 0 || idx >= len(options) {
		var t T
		return t, fmt.Errorf("invalid choice: must be between 1 and %d", len(options))
	}

	return options[idx], nil
}
