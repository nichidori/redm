package redmineapi

import "context"

type IssueStatus struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsClosed bool   `json:"is_closed"`
}

type IssueStatusesResponse struct {
	IssueStatuses []IssueStatus `json:"issue_statuses"`
}

func (c *Client) GetIssueStatuses(ctx context.Context) (*IssueStatusesResponse, error) {
	return doGet[IssueStatusesResponse](c, ctx, "/issue_statuses.json", nil)
}
