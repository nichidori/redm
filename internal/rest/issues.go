package rest

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

type Issue struct {
	ID             int               `json:"id"`
	Project        IDName            `json:"project"`
	Tracker        IDName            `json:"tracker"`
	Status         IDName            `json:"status"`
	Priority       IDName            `json:"priority"`
	Author         IDName            `json:"author"`
	AssignedTo     *IDName           `json:"assigned_to"`
	Category       *IDName           `json:"category"`
	FixedVersion   *IDName           `json:"fixed_version,omitempty"`
	Parent         *struct{ ID int } `json:"parent,omitempty"`
	Subject        string            `json:"subject"`
	Description    string            `json:"description"`
	StartDate      string            `json:"start_date"`
	DueDate        string            `json:"due_date"`
	DoneRatio      int               `json:"done_ratio"`
	IsPrivate      bool              `json:"is_private"`
	EstimatedHours *float64          `json:"estimated_hours"`
	CustomFields   []CustomField     `json:"custom_fields"`
	CreatedOn      time.Time         `json:"created_on"`
	UpdatedOn      time.Time         `json:"updated_on"`
	ClosedOn       *time.Time        `json:"closed_on"`
}

type IssuesResponse struct {
	Issues []Issue `json:"issues"`
	PaginationMeta
}

// IssueFilter holds optional query parameters for listing issues.
// Zero values are omitted from the query string.
type IssueFilter struct {
	ProjectID    int
	TrackerID    int
	StatusID     string
	AssignedToID string
	Offset       int
	Limit        int
	Sort         string
}

func (f IssueFilter) toQuery() url.Values {
	q := url.Values{}
	if f.ProjectID > 0 {
		q.Set("project_id", strconv.Itoa(f.ProjectID))
	}
	if f.TrackerID > 0 {
		q.Set("tracker_id", strconv.Itoa(f.TrackerID))
	}
	if f.StatusID != "" {
		q.Set("status_id", f.StatusID)
	}
	if f.AssignedToID != "" {
		q.Set("assigned_to_id", f.AssignedToID)
	}
	if f.Offset > 0 {
		q.Set("offset", strconv.Itoa(f.Offset))
	}
	if f.Limit > 0 {
		q.Set("limit", strconv.Itoa(f.Limit))
	}
	if f.Sort != "" {
		q.Set("sort", f.Sort)
	}
	return q
}

func (c *Client) GetIssues(ctx context.Context, filter IssueFilter) (*IssuesResponse, error) {
	return doGet[IssuesResponse](c, ctx, "/issues.json", filter.toQuery())
}
