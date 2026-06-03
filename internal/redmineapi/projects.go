package redmineapi

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

type Project struct {
	ID                int           `json:"id"`
	Name              string        `json:"name"`
	Identifier        string        `json:"identifier"`
	Description       string        `json:"description"`
	Status            int           `json:"status"`
	IsPublic          bool          `json:"is_public"`
	InheritMembers    bool          `json:"inherit_members"`
	Homepage          string        `json:"homepage,omitempty"`
	Parent            *IDName       `json:"parent,omitempty"`
	CustomFields      []CustomField `json:"custom_fields,omitempty"`
	IssueCustomFields []IDName      `json:"issue_custom_fields,omitempty"`
	IssueCategories   []IDName      `json:"issue_categories,omitempty"`
	CreatedOn         time.Time     `json:"created_on"`
	UpdatedOn         time.Time     `json:"updated_on"`
}

type ProjectsResponse struct {
	Projects []Project `json:"projects"`
	PaginationMeta
}

// GetProjects returns all projects the user has access to.
// Use offset/limit for pagination (default: offset=0, limit=25, max=100).
func (c *Client) GetProjects(ctx context.Context, offset, limit int) (*ProjectsResponse, error) {
	query := url.Values{}
	query.Set("include", "issue_custom_fields,issue_categories")
	if offset > 0 {
		query.Set("offset", strconv.Itoa(offset))
	}
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	return doGet[ProjectsResponse](c, ctx, "/projects.json", query)
}
