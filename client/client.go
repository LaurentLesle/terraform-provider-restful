package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// ErrNotFound is expected to be returned for `Read` when the resource with the specified id doesn't exist.
var ErrNotFound = errors.New("resource not found")

type Client struct {
	*resty.Client

	// Per request options
	createMethod string
	contentType  string
}

type ClientBuilder struct {
	BaseURL string
	Option  *Option
}

func (b ClientBuilder) Build(ctx context.Context) (*Client, error) {
	opt := b.Option
	if opt == nil {
		opt = &Option{}
	}

	client := resty.New()
	if opt.Security != nil {
		client = opt.Security.newClient(ctx)
	}

	if _, err := url.Parse(b.BaseURL); err != nil {
		return nil, err
	}

	client.SetBaseURL(b.BaseURL)

	createMethod := "POST"
	if opt.CreateMethod != "" {
		createMethod = opt.CreateMethod
	}

	contentType := "application/json"
	if opt.ContentType != "" {
		contentType = opt.ContentType
	}

	return &Client{
		Client:       client,
		createMethod: createMethod,
		contentType:  contentType,
	}, nil
}

func (c *Client) Create(path string, body interface{}) ([]byte, error) {
	req := c.R().SetBody(body)
	if c.contentType != "" {
		req = req.SetHeader("Content-Type", c.contentType)
	}
	switch c.createMethod {
	case "POST":
		resp, err := req.Post(path)
		if err != nil {
			return nil, err
		}
		// TODO: Support LRO
		if resp.StatusCode()/100 != 2 {
			return nil, fmt.Errorf("Unexpected response (%s - code: %d): %s", resp.Status(), resp.StatusCode(), string(resp.Body()))
		}
		return resp.Body(), nil
	case "PUT":
		panic("TBD")
	}
	return nil, fmt.Errorf("unknown create method: %s", c.createMethod)
}

func (c *Client) Read(path string) ([]byte, error) {
	resp, err := c.R().Get(path)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resp.StatusCode()/100 != 2 {
		return nil, fmt.Errorf("Unexpected response (%s - code: %d): %s", resp.Status(), resp.StatusCode(), string(resp.Body()))
	}
	return resp.Body(), nil
}

func (c *Client) Update(path string, body interface{}) ([]byte, error) {
	req := c.R().SetBody(body)
	if c.contentType != "" {
		req = req.SetHeader("Content-Type", c.contentType)
	}
	resp, err := req.Put(path)
	if err != nil {
		return nil, err
	}
	// TODO: Support LRO
	if resp.StatusCode()/100 != 2 {
		return nil, fmt.Errorf("Unexpected response (%s - code: %d): %s", resp.Status(), resp.StatusCode(), string(resp.Body()))
	}
	return resp.Body(), nil
}

func (c *Client) Delete(path string) ([]byte, error) {
	resp, err := c.R().Delete(path)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrNotFound
	}
	// TODO: Support LRO
	if resp.StatusCode()/100 != 2 {
		return nil, fmt.Errorf("Unexpected response (%s - code: %d): %s", resp.Status(), resp.StatusCode(), string(resp.Body()))
	}
	return resp.Body(), nil
}
