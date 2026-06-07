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

var CommandUpdate = Command{
	Name:        "update",
	Description: "Update an issue",
	Setup: func(fs *flag.FlagSet, s *state) func() error {
		return func() error {
			if s.config == nil {
				return fmt.Errorf("not logged in")
			}

			args := fs.Args()
			if len(args) < 1 {
				return fmt.Errorf("usage: redm update <issue-id>")
			}

			issueID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid issue ID: %s", args[0])
			}

			issue, err := s.client.GetIssue(context.Background(), issueID)
			if err != nil {
				return fmt.Errorf("failed to fetch issue: %w", err)
			}

			reader := bufio.NewReader(os.Stdin)

			type field struct {
				label   string
				current func(redmineapi.Issue) string
				update  func(*redmineapi.UpdateIssueRequest, *bufio.Reader) error
			}

			fields := []field{
				{"Project", func(i redmineapi.Issue) string { return i.Project.Name }, func(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
					projects, err := s.client.GetProjects(context.Background(), 0, 100)
					if err != nil {
						return fmt.Errorf("failed to fetch projects: %w", err)
					}
					p, err := SelectOption(rd, "project", projects.Projects, func(p redmineapi.Project) string { return p.Name })
					if err != nil {
						return err
					}
					r.ProjectID = p.ID
					return nil
				}},
				{"Category", func(i redmineapi.Issue) string {
					if i.Category != nil {
						return i.Category.Name
					}
					return ""
				}, func(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
					projects, err := s.client.GetProjects(context.Background(), 0, 100)
					if err != nil {
						return fmt.Errorf("failed to fetch projects: %w", err)
					}
					var proj *redmineapi.Project
					for i := range projects.Projects {
						if projects.Projects[i].ID == issue.Project.ID {
							proj = &projects.Projects[i]
							break
						}
					}
					if proj == nil {
						return fmt.Errorf("project not found in fetched list")
					}
					if len(proj.IssueCategories) == 0 {
						fmt.Println("No categories available for this project")
						return nil
					}
					cat, err := SelectOption(rd, "category", proj.IssueCategories, func(c redmineapi.IDName) string { return c.Name })
					if err != nil {
						return err
					}
					r.CategoryID = cat.ID
					return nil
				}},
				{"Tracker", func(i redmineapi.Issue) string { return i.Tracker.Name }, func(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
					trackers, err := s.client.GetTrackers(context.Background())
					if err != nil {
						return fmt.Errorf("failed to fetch trackers: %w", err)
					}
					t, err := SelectOption(rd, "tracker", trackers.Trackers, func(t redmineapi.Tracker) string { return t.Name })
					if err != nil {
						return err
					}
					r.TrackerID = t.ID
					return nil
				}},
				{"Subject", func(i redmineapi.Issue) string { return i.Subject }, func(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
					fmt.Print("New Subject: ")
					v, err := rd.ReadString('\n')
					if err != nil {
						return err
					}
					r.Subject = strings.TrimSpace(v)
					return nil
				}},
				{"Description", func(i redmineapi.Issue) string { return i.Description }, func(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
					fmt.Print("New Description: ")
					v, err := rd.ReadString('\n')
					if err != nil {
						return err
					}
					r.Description = strings.TrimSpace(v)
					return nil
				}},
				{"Start Date", func(i redmineapi.Issue) string { return i.StartDate }, func(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
					fmt.Print("New Start Date (YYYY-MM-DD): ")
					v, err := rd.ReadString('\n')
					if err != nil {
						return err
					}
					r.StartDate = strings.TrimSpace(v)
					return nil
				}},
				{"Due Date", func(i redmineapi.Issue) string { return i.DueDate }, func(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
					fmt.Print("New Due Date (YYYY-MM-DD): ")
					v, err := rd.ReadString('\n')
					if err != nil {
						return err
					}
					r.DueDate = strings.TrimSpace(v)
					return nil
				}},
				{"Status", func(i redmineapi.Issue) string { return i.Status.Name }, func(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
					statuses, err := s.client.GetIssueStatuses(context.Background())
					if err != nil {
						return fmt.Errorf("failed to fetch statuses: %w", err)
					}
					st, err := SelectOption(rd, "status", statuses.IssueStatuses, func(st redmineapi.IssueStatus) string { return st.Name })
					if err != nil {
						return err
					}
					r.StatusID = st.ID
					return nil
				}},
				{"Progress", func(i redmineapi.Issue) string { return strconv.Itoa(i.DoneRatio) + "%" }, func(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
					fmt.Print("New Progress (0-100): ")
					v, err := rd.ReadString('\n')
					if err != nil {
						return err
					}
					ratio, err := strconv.Atoi(strings.TrimSpace(v))
					if err != nil {
						return fmt.Errorf("invalid progress value: %s", strings.TrimSpace(v))
					}
					if ratio < 0 || ratio > 100 {
						return fmt.Errorf("progress must be between 0 and 100")
					}
					r.DoneRatio = ratio
					return nil
				}},
			}

			selected, err := SelectOption(reader, "field to update", fields, func(f field) string {
				return FixLength(f.label, 12) + " : " + f.current(*issue)
			})
			if err != nil {
				return err
			}
			fmt.Println()

			req := &redmineapi.UpdateIssueRequest{}
			if err := selected.update(req, reader); err != nil {
				return err
			}

			if err := s.client.UpdateIssue(context.Background(), issueID, req); err != nil {
				return fmt.Errorf("failed to update issue: %w", err)
			}

			fmt.Printf("Issue #%d updated\n", issueID)
			return nil
		}
	},
}
