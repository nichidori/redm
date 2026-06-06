package redmineapi

import (
	"context"
	"time"
)

type User struct {
	ID              int        `json:"id"`
	Login           string     `json:"login"`
	FirstName       string     `json:"firstname"`
	LastName        string     `json:"lastname"`
	Mail            string     `json:"mail"`
	APIKey          string     `json:"api_key,omitempty"`
	Status          int        `json:"status,omitempty"`
	AvatarURL       string     `json:"avatar_url,omitempty"`
	CreatedOn       *time.Time `json:"created_on,omitempty"`
	UpdatedOn       *time.Time `json:"updated_on,omitempty"`
	LastLoginOn     *time.Time `json:"last_login_on,omitempty"`
	PasswdChangedOn *time.Time `json:"passwd_changed_on,omitempty"`
}

type UserResponse struct {
	User User `json:"user"`
}

func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	resp, err := doGet[UserResponse](c, ctx, "/users/current.json", nil)
	if err != nil {
		return nil, err
	}
	return &resp.User, nil
}
