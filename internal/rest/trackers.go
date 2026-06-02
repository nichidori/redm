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
	return doGet[TrackersResponse](c, ctx, "/trackers.json", nil)
}
