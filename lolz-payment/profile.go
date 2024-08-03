package lolzpayment

import (
	"context"
	"encoding/json"
	"net/http"
)

type Profile struct {
	// UserID is just a user id man, why are you reading this?
	UserID int `json:"user_id"`
	// Username is just a username man, why are you reading this?
	Username string `json:"username"`
}

// Me is a function to get profile information using https://api.lzt.market/me endpoint.
func (c *Client) Me(ctx context.Context) (*Profile, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "me",
	}

	data, err := c.do(ctx, r)
	if err != nil {
		return nil, err
	}

	var res Profile
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
