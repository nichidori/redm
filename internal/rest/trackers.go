package rest

import "context"

type Tracker struct {
	ID                    int      `json:"id"`
	Name                  string   `json:"name"`
	DefaultStatus         IDName   `json:"default_status"`
	Description           string   `json:"description,omitempty"`
	EnabledStandardFields []string `json:"enabled_standard_fields,omitempty"`
}

type TrackersResponse struct {
	Trackers []Tracker `json:"trackers"`
}

func (c *Client) GetTrackers(ctx context.Context) (*TrackersResponse, error) {
	resp := &TrackersResponse{}
	if err := c.doGet(ctx, "/trackers.json", nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
