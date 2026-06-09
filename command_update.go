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

type updateFieldHandlers struct {
	client *redmineapi.Client
	issue  *redmineapi.Issue
}

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
				current string
				update  func(*redmineapi.UpdateIssueRequest, *bufio.Reader) error
			}

			h := &updateFieldHandlers{client: s.client, issue: issue}

			catName := ""
			if issue.Category != nil {
				catName = issue.Category.Name
			}

			fields := []field{
				{"Project", issue.Project.Name, h.updateProject},
				{"Category", catName, h.updateCategory},
				{"Tracker", issue.Tracker.Name, h.updateTracker},
				{"Subject", issue.Subject, h.updateSubject},
				{"Description", issue.Description, h.updateDescription},
				{"Start Date", issue.StartDate, h.updateStartDate},
				{"Due Date", issue.DueDate, h.updateDueDate},
				{"Status", issue.Status.Name, h.updateStatus},
				{"Progress", strconv.Itoa(issue.DoneRatio) + "%", h.updateProgress},
			}

			for _, cf := range issue.CustomFields {
				val := ""
				if cf.Value != nil {
					val = fmt.Sprintf("%v", cf.Value)
				}
				fields = append(fields, field{
					cf.Name, val, func(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
						fmt.Print("New value: ")
						v, err := rd.ReadString('\n')
						if err != nil {
							return err
						}
						v = strings.TrimSpace(v)
						if v == "" {
							return nil
						}
						r.CustomFields = []redmineapi.CustomField{{ID: cf.ID, Value: v}}
						return nil
					},
				})
			}

			selected, err := SelectOption(reader, "field to update", fields, func(f field) string {
				return FixLength(f.label, 16) + " : " + f.current
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

func (h *updateFieldHandlers) updateProject(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
	projects, err := h.client.GetProjects(context.Background(), 0, 100)
	if err != nil {
		return fmt.Errorf("failed to fetch projects: %w", err)
	}
	p, err := SelectOption(rd, "project", projects.Projects, func(p redmineapi.Project) string { return p.Name })
	if err != nil {
		return err
	}
	r.ProjectID = p.ID
	return nil
}

func (h *updateFieldHandlers) updateCategory(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
	projects, err := h.client.GetProjects(context.Background(), 0, 100)
	if err != nil {
		return fmt.Errorf("failed to fetch projects: %w", err)
	}
	var proj *redmineapi.Project
	for i := range projects.Projects {
		if projects.Projects[i].ID == h.issue.Project.ID {
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
}

func (h *updateFieldHandlers) updateTracker(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
	trackers, err := h.client.GetTrackers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch trackers: %w", err)
	}
	t, err := SelectOption(rd, "tracker", trackers.Trackers, func(t redmineapi.Tracker) string { return t.Name })
	if err != nil {
		return err
	}
	r.TrackerID = t.ID
	return nil
}

func (h *updateFieldHandlers) updateSubject(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
	fmt.Print("New Subject: ")
	v, err := rd.ReadString('\n')
	if err != nil {
		return err
	}
	r.Subject = strings.TrimSpace(v)
	return nil
}

func (h *updateFieldHandlers) updateDescription(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
	fmt.Print("New Description: ")
	v, err := rd.ReadString('\n')
	if err != nil {
		return err
	}
	r.Description = strings.TrimSpace(v)
	return nil
}

func (h *updateFieldHandlers) updateStartDate(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
	fmt.Print("New Start Date (YYYY-MM-DD): ")
	v, err := rd.ReadString('\n')
	if err != nil {
		return err
	}
	r.StartDate = strings.TrimSpace(v)
	return nil
}

func (h *updateFieldHandlers) updateDueDate(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
	fmt.Print("New Due Date (YYYY-MM-DD): ")
	v, err := rd.ReadString('\n')
	if err != nil {
		return err
	}
	r.DueDate = strings.TrimSpace(v)
	return nil
}

func (h *updateFieldHandlers) updateStatus(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
	statuses, err := h.client.GetIssueStatuses(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch statuses: %w", err)
	}
	st, err := SelectOption(rd, "status", statuses.IssueStatuses, func(st redmineapi.IssueStatus) string { return st.Name })
	if err != nil {
		return err
	}
	r.StatusID = st.ID
	return nil
}

func (h *updateFieldHandlers) updateProgress(r *redmineapi.UpdateIssueRequest, rd *bufio.Reader) error {
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
}
