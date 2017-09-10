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

package perfops

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func TestIsUnauthorized(t *testing.T) {
	err := &clientError{http.StatusUnauthorized, ""}
	if !IsUnauthorized(err) {
		t.Fatalf("expected IsUnauthorized; got %v", err)
	}
}

func TestNewClient(t *testing.T) {
	testClient := &http.Client{}
	testCases := map[string]struct {
		opts []func(c *Client) error
		chk  func(c *Client) (interface{}, interface{})
	}{
		"Default":        {[]func(c *Client) error{}, func(c *Client) (interface{}, interface{}) { return nil, nil }},
		"WithAPIKey":     {[]func(c *Client) error{WithAPIKey("abc")}, func(c *Client) (interface{}, interface{}) { return c.apiKey, "abc" }},
		"WithHTTPClient": {[]func(c *Client) error{WithHTTPClient(testClient)}, func(c *Client) (interface{}, interface{}) { return c.client, testClient }},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c, err := NewClient(tc.opts...)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if got, exp := tc.chk(c); got != exp {
				t.Fatalf("expected %v; got %v", exp, got)
			}
		})
	}
}

type roundTripper interface {
	RoundTrip(req *http.Request) (*http.Response, error)
}

func newTestClient(tr roundTripper) (*Client, error) {
	c := &http.Client{Transport: tr}
	return NewClient(WithHTTPClient(c))
}

type recordingTransport struct {
	req *http.Request
}

func (t *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.req = req
	return nil, errors.New("dummy impl")
}

type respondingTransport struct {
	resp *http.Response
}

func (t *respondingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.resp, nil
}

type testTransport struct {
	req  *http.Request
	resp *http.Response
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.req = req
	return t.resp, nil
}

func dummyReq(method string) *http.Request {
	return &http.Request{Method: method, Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1}
}

func dummyResp(status int, method, body string) *http.Response {
	return &http.Response{
		Status:     strconv.Itoa(status) + " Meep meep",
		StatusCode: status,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Request:    dummyReq(method),
		Header: http.Header{
			"Connection":     {"close"},
			"Content-Length": {strconv.Itoa(len(body))},
		},
		Close:         true,
		ContentLength: int64(len(body)),
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
	}
}
