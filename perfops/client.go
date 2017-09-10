// Copyright 2017 The PerfOps-CLI Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package perfops provides functionality to access Prospect One's PerfOps APIs.
package perfops

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
)

// UserAgent is the header string used to identify this package.
const UserAgent = "PerfOps API Go Client"

const (
	apiRoot    = "https://api.perfops.net"
	basePath   = apiRoot
	libVersion = "v1.0.0"
	userAgent  = UserAgent + "/" + libVersion + " (" + runtime.GOOS + "/" + runtime.GOARCH + ")"
)

type (
	// Client defines the API client interface.
	Client struct {
		client *http.Client
		common service

		BasePath  string // API endpoint base URL
		UserAgent string // optional additional User-Agent fragment
		apiKey    string

		Run *RunService
	}

	service struct {
		client *Client
	}

	clientError struct {
		Code int
		Text string
	}
)

// Error returns
func (err clientError) Error() string {
	if err.Text != "" {
		return fmt.Sprintf("%d: %s", err.Code, err.Text)
	}
	return fmt.Sprintf("%d", err.Code)
}

// IsUnauthorized returns a value indicating whether the error is an
// authorization error.
func (err *clientError) IsUnauthorized() bool {
	return err.Code == http.StatusUnauthorized
}

type isUnauthorizeder interface {
	IsUnauthorized() bool
}

// IsUnauthorized returns a value indicating whether the error is an
// authorization error.
func IsUnauthorized(err error) bool {
	if ce, ok := err.(isUnauthorizeder); ok {
		return ce.IsUnauthorized()
	}
	return false
}

// WithAPIKey sets the API key for the API client.
func WithAPIKey(key string) func(c *Client) error {
	return func(c *Client) error {
		c.apiKey = key
		return nil
	}
}

// WithHTTPClient sets the HTTP client for the API client.
func WithHTTPClient(client *http.Client) func(c *Client) error {
	return func(c *Client) error {
		if client == nil {
			return errors.New("HTTP client is nil")
		}
		c.client = client
		return nil
	}
}

// NewClient returns a new Client given an API key and options.
func NewClient(opts ...func(c *Client) error) (*Client, error) {
	c := &Client{
		client:   http.DefaultClient,
		BasePath: basePath,
	}
	c.common.client = c
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	c.Run = (*RunService)(&c.common)

	return c, nil
}

func (c *Client) userAgent() string {
	if c.UserAgent == "" {
		return userAgent
	}
	return c.UserAgent + " " + userAgent
}

func (c *Client) do(req *http.Request, v interface{}) error {
	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent())
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer closeBody(resp)
	if resp.StatusCode >= http.StatusBadRequest {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		//fmt.Printf("***\n%v\n***\n", string(b))
		return &clientError{Code: resp.StatusCode, Text: string(b)}
	}
	d := json.NewDecoder(resp.Body)
	return d.Decode(&v)
}

func newJSONReader(v interface{}) (io.Reader, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	return buf, nil
}

func closeBody(res *http.Response) {
	if res == nil || res.Body == nil {
		return
	}
	res.Body.Close()
}
