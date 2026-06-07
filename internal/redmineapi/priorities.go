package redmineapi

import "context"

type IssuePriority struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

type IssuePrioritiesResponse struct {
	IssuePriorities []IssuePriority `json:"issue_priorities"`
}

func (c *Client) GetIssuePriorities(ctx context.Context) (*IssuePrioritiesResponse, error) {
	return doGet[IssuePrioritiesResponse](c, ctx, "/enumerations/issue_priorities.json", nil)
}
