package lolzpayment

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	baseUrl = "https://api.lzt.market/"
)

type Client struct {
	token   string
	baseURL string

	userID   int
	username string

	httpClient *http.Client
}

type Options struct {
	// Token from your lolz app with access to the market.
	Token string

	// ClientTimeout field is optional. Default is 30s.
	ClientTimeout time.Duration

	// BaseURL field is optional. Default is https://api.lzt.market/.
	BaseURL string

	// ProxyURL field is optional. Example: "http://user:pass@ip:port" or "socks5://ip:port"
	ProxyURL string
}

// NewClient creates a new lolz_payment client to interact with api.
func NewClient(options Options) (*Client, error) {
	c := &Client{
		token: options.Token,
	}

	clientTimeout := 30 * time.Second
	if options.ClientTimeout != 0 {
		clientTimeout = options.ClientTimeout
	}

	c.baseURL = baseUrl
	if options.BaseURL != "" {
		c.baseURL = options.BaseURL
	}

	transport := &http.Transport{
		MaxIdleConns:       100,
		IdleConnTimeout:    90 * time.Second,
		DisableCompression: true,
	}

	// Добавляем прокси если он указан
	if options.ProxyURL != "" {
		proxyURL, err := url.Parse(options.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	c.httpClient = &http.Client{
		Transport: transport,
		Timeout:   clientTimeout,
	}

	profile, err := c.Me(context.Background())
	if err != nil {
		return nil, err
	}

	c.userID = profile.UserID
	c.username = profile.Username

	return c, nil
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
