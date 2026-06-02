package rest

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Identifier  string    `json:"identifier"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	IsPublic    bool      `json:"is_public"`
	Homepage    string    `json:"homepage,omitempty"`
	CreatedOn   time.Time `json:"created_on"`
	UpdatedOn   time.Time `json:"updated_on"`
}

type ProjectsResponse struct {
	Projects []Project `json:"projects"`
	PaginationMeta
}

// GetProjects returns all projects the user has access to.
// Use offset/limit for pagination (default: offset=0, limit=25, max=100).
func (c *Client) GetProjects(ctx context.Context, offset, limit int) (*ProjectsResponse, error) {
	query := url.Values{}
	if offset > 0 {
		query.Set("offset", strconv.Itoa(offset))
	}
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	return doGet[ProjectsResponse](c, ctx, "/projects.json", query)
}
