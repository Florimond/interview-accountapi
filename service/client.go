package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/Florimond/interview-accountapi/service/contracts"
	//"github.com/Florimond/interview-accountapi/client/account"
)

// Client is a struct that represents a client to the API.
type Client struct {
	BaseURL string //url.URL
	http    *http.Client

	//Account *account.Account
}

// TODO: options, WithTimeout()
// TODO: async and sync
// TODO: modules for each part? Account, etc Client.Account.Fetch, Client.User...
// TODO: version management
// TODO: context to cancel requests!

// NewClient creates a new client.
func NewClient(baseURL string, timeout time.Duration) *Client {

	// TODO strict url checking
	/*
		url, err := url.Parse(baseURL)
		if err != nil {
			panic(err)
		}*/

	return &Client{
		BaseURL: baseURL,
		http: &http.Client{
			Timeout: timeout,
		},
	}
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"error_message"`
}

type Response struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
}

// As decodes the response
func (r *Response) As(v interface{}) error {
	return json.Unmarshal(r.Data, v)
}

func (r *Response) IsErrorStatus() bool {
	return r.Code < http.StatusOK || r.Code >= http.StatusBadRequest
}

func (c *Client) sendRequest(req *http.Request) (*Response, error) {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/vnd.api+json")

	ctx := req.Context()
	attempts := 0
	var res *http.Response
	var err error
	for {
		log.Println("Attempting request...")
		res, err = c.http.Do(req)
		if err != nil {
			return nil, err
		}

		if res.StatusCode == 0 || res.StatusCode == 500 || res.StatusCode > 501 {
			log.Println("StatusCode ", res.StatusCode, " Waiting before retry...")
			select {
			case <-time.After(time.Duration(math.Pow(float64(attempts), 2.0)*50) * time.Millisecond): // Retry immediately the first time
				attempts++
			case <-ctx.Done():
				log.Println("Request cancelled.")
				return nil, ctx.Err()
			}
		} else {
			break
		}
	}

	defer res.Body.Close()

	var response Response
	response.Code = res.StatusCode

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
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

// WithPageNumber builds a string reprenting a url option for the page number of the results
func WithPageNumber(n uint) string {
	return fmt.Sprint("page[number]=", n)
}

// WithPageSize builds a string representing a url option for the page size of the results
func WithPageSize(s uint) string {
	return fmt.Sprint("page[size]=", s)
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
func (c *Client) Delete(ctx context.Context, provider contracts.Provider, id string) (*Response, error) {
	url := fmt.Sprintf("%s/%s", provider.Path(), id)
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
