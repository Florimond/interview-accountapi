package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/Florimond/interview-accountapi/service/contracts"
)

// Client is a struct that represents a client to the API.
type Client struct {
	BaseURL string //url.URL
	http    *http.Client
}

// Error represents an HTTP error with a code.
type Error interface {
	error
	Code() int
}

// NewClient creates a new client.
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		BaseURL: baseURL,
		http: &http.Client{
			Timeout: timeout,
		},
	}
}

type errorResponse struct {
	HTTPCode int    `json:"code"`
	Message  string `json:"error_message"`
}

// Error complies with golang error interface
func (r *errorResponse) Error() string {
	return r.Message
}

// Code returns an HTTP response code
func (r *errorResponse) Code() int {
	return r.HTTPCode
}

// Response represents a parsed response
type Response struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
}

// As decodes the response
func (r *Response) As(v interface{}) error {
	return json.Unmarshal(r.Data, v)
}

func (c *Client) sendRequest(req *http.Request) (*Response, error) {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/vnd.api+json")

	ctx := req.Context()
	attempts := 0

retry:
	for {
		res, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}

		switch {
		case res.StatusCode == 0 || res.StatusCode == 500 || res.StatusCode > 501:
			select {
			case <-time.After(time.Duration(math.Pow(float64(attempts), 2.0)*50) * time.Millisecond): // Retry immediately the first time
				attempts++
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		case res.StatusCode >= 200 && res.StatusCode <= 299:
			return parseResponse(res)
		case res.StatusCode >= 400 && res.StatusCode <= 599:
			return nil, parseError(res)
		default:
			break retry
		}
	}

	return nil, fmt.Errorf("unknown error")
}

// parseResponse parses a successful response
func parseResponse(res *http.Response) (*Response, error) {
	defer res.Body.Close()
	response := &Response{
		Code: res.StatusCode,
	}

	if err := json.NewDecoder(res.Body).Decode(response); err != nil && err != io.EOF { // io.EOF means an empty body
		return nil, err
	}
	return response, nil
}

// parseError parses an error response
func parseError(res *http.Response) error {
	defer res.Body.Close()
	response := &errorResponse{
		HTTPCode: res.StatusCode,
	}

	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return err
	}
	return response
}

// formatOptions formats a set of URL options
func formatOptions(options []string) string {
	opts, hasOpts := "", false
	if options != nil && len(options) > 0 {
		for _, option := range options {
			if !hasOpts {
				hasOpts = true
				opts += "?"
			} else {
				opts += "&"
			}
			opts += option
		}
	}
	return opts
}

// Trim removes both suffix and prefix
func trim(v string) string {
	return strings.TrimSuffix(strings.TrimPrefix(v, "/"), "/")
}

func (c *Client) makeRequest(ctx context.Context, method, url string, urlOptions []string, body io.Reader) (*http.Request, error) {
	fullURL := fmt.Sprintf("%s/%s%s", trim(c.BaseURL), trim(url), formatOptions(urlOptions))
	if ctx != nil {
		return http.NewRequestWithContext(ctx, method, fullURL, body)
	}
	return http.NewRequest(method, fullURL, body)
}

// FindByID finds a document by its id
func (c *Client) FindByID(ctx context.Context, provider contracts.Provider, id string) (*Response, error) {
	url := fmt.Sprintf("%s/%s", provider.Path(), id)
	req, err := c.makeRequest(ctx, "GET", url, nil, nil)
	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

// List returns a list of documents for a provider
func (c *Client) List(ctx context.Context, provider contracts.Provider, options ...string) (*Response, error) {
	req, err := c.makeRequest(ctx, "GET", provider.Path(), options, nil)
	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

// Delete deletes a document by its id
func (c *Client) Delete(ctx context.Context, provider contracts.Provider, id string, version uint) (*Response, error) {
	url := fmt.Sprintf("%s/%s?version=%d/", provider.Path(), id, version)
	req, err := c.makeRequest(ctx, "DELETE", url, nil, nil)
	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

// Create creates a document
func (c *Client) Create(ctx context.Context, provider contracts.Provider, doc interface{}) (*Response, error) {
	body := struct {
		Data interface{} `json:"data"`
	}{
		Data: doc,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := c.makeRequest(ctx, "POST", provider.Path(), nil, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

// IsHTTPError checks if an error is an HTTP error or not
func IsHTTPError(r error) (int, bool) {
	err, ok := r.(Error)
	if !ok {
		return 0, false
	}
	return err.Code(), true
}
