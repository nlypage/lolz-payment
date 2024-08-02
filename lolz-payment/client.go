package lolzpayment

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	token   string
	baseURL string

	httpClient *http.Client
}

type Options struct {
	// Token from your lolz app with access to the market.
	Token string

	// ClientTimeout field is optional. Default is 30s.
	ClientTimeout time.Duration

	// BaseURL field is optional. Default is https://api.lzt.market/.
	BaseURL string
}

// NewClient creates a new lolz_payment client to interact with api.
func NewClient(options Options) *Client {
	c := &Client{
		token: options.Token,
	}
	clientTimeout := 30 * time.Second
	if options.ClientTimeout != 0 {
		clientTimeout = options.ClientTimeout
	}

	c.baseURL = "https://api.lzt.market/"
	if options.BaseURL != "" {
		c.baseURL = options.BaseURL
	}

	c.httpClient = &http.Client{
		Timeout: clientTimeout,
	}

	return c
}

func (c *Client) parseRequest(r *request) (err error) {
	// set request options from user

	err = r.validate()
	if err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.baseURL, r.endpoint)

	queryString := r.query.Encode()
	header := http.Header{}
	if r.header != nil {
		header = r.header.Clone()
	}
	header.Add("authorization", fmt.Sprintf("Bearer %s", c.token))

	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}

	r.fullURL = fullURL
	r.header = header
	return nil
}

func (c *Client) do(ctx context.Context, r *request) (data []byte, err error) {
	err = c.parseRequest(r)
	if err != nil {
		return []byte{}, err
	}
	req, err := http.NewRequest(r.method, r.fullURL, nil)
	if err != nil {
		return []byte{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header

	res, err := c.httpClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		cerr := res.Body.Close()
		// Only overwrite the retured error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("error response from the server with code %d: %s", res.StatusCode, data)
	}

	return data, nil
}
