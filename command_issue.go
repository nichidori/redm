package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/nichidori/redm/internal/redmineapi"
)

var CommandIssue = Command{
	Name:        "issue",
	Description: "List issues",
	Setup: func(fs *flag.FlagSet, s *state) func() error {
		var (
			project      string
			status       string
			allAssignees bool
		)

		fs.StringVar(&project, "p", "", "Filter by project ID or name")
		fs.StringVar(&status, "s", "", "Filter by status ID or name")
		fs.BoolVar(&allAssignees, "a", false, "Show issues for all assignees")

		return func() error {
			if s.config == nil {
				return fmt.Errorf("not logged in")
			}

			var filterProjectID int
			if project != "" {
				if id, err := strconv.Atoi(project); err == nil {
					filterProjectID = id
				} else {
					projects, err := s.client.GetProjects(context.Background(), 0, 100)
					if err != nil {
						return fmt.Errorf("failed to fetch projects: %w", err)
					}
					found := false
					for _, p := range projects.Projects {
						if strings.EqualFold(p.Name, project) {
							filterProjectID = p.ID
							found = true
							break
						}
					}
					if !found {
						return fmt.Errorf("project not found: %s", project)
					}
				}
			}

			var filterStatusID string
			if status != "" {
				if _, err := strconv.Atoi(status); err == nil {
					filterStatusID = status
				} else {
					statuses, err := s.client.GetIssueStatuses(context.Background())
					if err != nil {
						return fmt.Errorf("failed to fetch issue statuses: %w", err)
					}
					found := false
					for _, st := range statuses.IssueStatuses {
						if strings.EqualFold(st.Name, status) {
							filterStatusID = strconv.Itoa(st.ID)
							found = true
							break
						}
					}
					if !found {
						return fmt.Errorf("status not found: %s", status)
					}
				}
			}

			filterAssignedToID := "me"
			if allAssignees {
				filterAssignedToID = ""
			}

			resp, err := s.client.GetIssues(context.Background(), redmineapi.IssueFilter{
				ProjectID:    filterProjectID,
				AssignedToID: filterAssignedToID,
				StatusID:     filterStatusID,
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
				NewColumn("Subject", 40, func(i int) string { return issues[i].Subject }),
				NewColumn("Tracker", 12, func(i int) string { return issues[i].Tracker.Name }),
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
