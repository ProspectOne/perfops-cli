// Copyright 2017 Prospect One https://prospectone.io/. All rights reserved.
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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// The maximum number of nodes allowed for requests without an API key.
const freeMaxNodeCap = 20

type (
	// RunService defines the interface for the run API
	RunService service

	// TestID represents the ID of an MTR or ping test.
	TestID string

	// NodeIDs represents a list of node IDs.
	NodeIDs []int

	// RunRequest represents the parameters for a ping request.
	RunRequest struct {
		// Target name
		Target string `json:"target"`
		// List of nodes ids, comma separated
		Nodes NodeIDs `json:"nodes,omitempty"`
		// Countries names, comma separated
		Location string `json:"location,omitempty"`
		// Max number of nodes
		Limit int `json:"limit,omitempty"`
	}

	// RunResult represents the result of an MTR or ping run.
	RunResult struct {
		Node     *Node       `json:"node,omitempty"`
		Output   interface{} `json:"output,omitempty"`
		Message  string      `json:"message,omitempty"`
		Finished interface{} `json:"finished"`
	}

	// RunItem represents an item of an MTR or ping output.
	RunItem struct {
		ID     string     `json:"id,omitempty"`
		Result *RunResult `json:"result,omitempty"`
	}

	// RunOutput represents the response of an MTR or ping output call.
	RunOutput struct {
		ID        string     `json:"id,omitempty"`
		Requested string     `json:"requested,omitempty"`
		Finished  string     `json:"finished"`
		Items     []*RunItem `json:"items,omitempty"`
	}

	// DNSPerfRequest represents the parameters for a DNS perf request.
	DNSPerfRequest struct {
		Target    string  `json:"target,omitempty"`
		DNSServer string  `json:"dnsServer,omitempty"`
		Nodes     NodeIDs `json:"nodes,omitempty"`
		Location  string  `json:"location,omitempty"`
		Limit     int     `json:"limit,omitempty"`
	}

	// DNSResolveRequest represents the parameters for a DNS resolve request.
	DNSResolveRequest struct {
		Target    string  `json:"target,omitempty"`
		Param     string  `json:"param,omitempty"`
		DNSServer string  `json:"dnsServer,omitempty"`
		Nodes     NodeIDs `json:"nodes,omitempty"`
		Location  string  `json:"location,omitempty"`
		Limit     int     `json:"limit,omitempty"`
	}

	// DNSTestResult represents the result of a DNS perf and DNS resolve output.
	DNSTestResult struct {
		DNSServer string          `json:"dnsServer,omitempty"`
		Node      *Node           `json:"node,omitempty"`
		Output    json.RawMessage `json:"output,omitempty"`
		Message   string          `json:"message,omitempty"`
	}

	// DNSTestItem respresents an item of a DNS perf and DNS resolve output.
	DNSTestItem struct {
		ID     string         `json:"id,omitempty"`
		Result *DNSTestResult `json:"result,omitempty"`
	}

	// DNSTestOutput represents the response of a DNS perf and DNS resolve output call.
	DNSTestOutput struct {
		ID        string         `json:"id,omitempty"`
		Requested string         `json:"requested,omitempty"`
		Finished  string         `json:"finished"`
		Items     []*DNSTestItem `json:"items,omitempty"`
	}

	// CurlRequest represents the parameters for a curl request.
	CurlRequest struct {
		Target   string  `json:"target,omitempty"`
		Head     bool    `json:"head"`
		Insecure bool    `json:"insecure,omitempty"`
		HTTP2    bool    `json:"http2,omitempty"`
		Nodes    NodeIDs `json:"nodes,omitempty"`
		Location string  `json:"location,omitempty"`
		Limit    int     `json:"limit,omitempty"`
	}

	argError struct {
		name string
	}
)

// Error returns the stirng representaiton of the error.
func (e *argError) Error() string {
	return fmt.Sprintf("invalid argument: %s", e.name)
}

// ArgName returns the name of the argument.
func (e *argError) ArgName() string {
	return e.name
}

// IsArgError retruns a value indicating whether the error represents
// a parameter error or not.
func IsArgError(err error) bool {
	type argNamer interface {
		ArgName() string
	}
	_, ok := err.(argNamer)
	return ok
}

// MarshalJSON returns the JSON encoding of NodeIDs, e.g., a comma
// separated list.
func (n NodeIDs) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	l := len(n) - 1
	b.WriteRune('"')
	for i, id := range n {
		b.WriteString(strconv.Itoa(id))
		if i < l {
			b.WriteRune(',')
		}
	}
	b.WriteRune('"')
	return b.Bytes(), nil
}

// UnmarshalJSON parses the JSON-encoded data.
func (n *NodeIDs) UnmarshalJSON(data []byte) error {
	if len(data) > 2 {
		ids := bytes.Split(bytes.Trim(data, `"`), []byte(","))
		*n = make(NodeIDs, len(ids))
		for i, bid := range ids {
			id, err := strconv.Atoi(string(bid))
			if err != nil {
				return err
			}
			(*n)[i] = id
		}
	} else {
		*n = make(NodeIDs, 0)
	}
	return nil
}

// Latency runs a latency test.
func (s *RunService) Latency(ctx context.Context, latency *RunRequest) (TestID, error) {
	return s.doPostRunRequest(ctx, "/run/latency", latency)
}

// LatencyOutput returns the full latency output under a test ID.
func (s *RunService) LatencyOutput(ctx context.Context, latencyID TestID) (*RunOutput, error) {
	return s.doGetRunOutput(ctx, "/run/latency/", latencyID)
}

// MTR runs an MTR test.
func (s *RunService) MTR(ctx context.Context, mtr *RunRequest) (TestID, error) {
	return s.doPostRunRequest(ctx, "/run/mtr", mtr)
}

// MTROutput returns the full MTR output under a test ID.
func (s *RunService) MTROutput(ctx context.Context, mtrID TestID) (*RunOutput, error) {
	return s.doGetRunOutput(ctx, "/run/mtr/", mtrID)
}

// Ping runs a ping test.
func (s *RunService) Ping(ctx context.Context, ping *RunRequest) (TestID, error) {
	return s.doPostRunRequest(ctx, "/run/ping", ping)
}

// PingOutput returns the full ping output under a test ID.
func (s *RunService) PingOutput(ctx context.Context, pingID TestID) (*RunOutput, error) {
	return s.doGetRunOutput(ctx, "/run/ping/", pingID)
}

// Traceroute runs a traceroute test.
func (s *RunService) Traceroute(ctx context.Context, ping *RunRequest) (TestID, error) {
	return s.doPostRunRequest(ctx, "/run/traceroute", ping)
}

// TracerouteOutput returns the full traceroute output under a test ID.
func (s *RunService) TracerouteOutput(ctx context.Context, pingID TestID) (*RunOutput, error) {
	return s.doGetRunOutput(ctx, "/run/traceroute/", pingID)
}

// DNSPerf finds the time it takes to resolve a DNS record.
func (s *RunService) DNSPerf(ctx context.Context, perf *DNSPerfRequest) (TestID, error) {
	if !isValidTarget(perf.Target) {
		return "", &argError{"target"}
	}
	if perf.DNSServer != "" && !isValidTarget(perf.DNSServer) {
		return "", &argError{"dns server"}
	}
	if !isValidLimit(s.client.apiKey, perf.Limit) {
		return "", &argError{"limit"}
	}

	body, err := newJSONReader(perf)
	if err != nil {
		return "", err
	}
	u := s.client.BasePath + "/run/dns-perf"
	req, _ := http.NewRequest("POST", u, body)
	req = req.WithContext(ctx)
	var raw struct {
		Error string
		ID    string `json:"id"`
	}
	if err = s.client.do(req, &raw); err != nil {
		return "", err
	}
	if raw.Error != "" {
		return "", errors.New(raw.Error)
	}
	return TestID(raw.ID), nil
}

// DNSPerfOutput returns the full DNS perf output under a test ID.
func (s *RunService) DNSPerfOutput(ctx context.Context, perfID TestID) (*DNSTestOutput, error) {
	u := s.client.BasePath + "/run/dns-perf/" + string(perfID)
	req, _ := http.NewRequest("GET", u, nil)
	var v *DNSTestOutput
	err := s.client.do(req, &v)
	return v, err
}

// DNSResolve resolves a DNS record.
func (s *RunService) DNSResolve(ctx context.Context, resolve *DNSResolveRequest) (TestID, error) {
	if !isValidTarget(resolve.Target) {
		return "", &argError{"target"}
	}
	if resolve.Param == "" {
		return "", &argError{"param"}
	}
	if !isValidTarget(resolve.DNSServer) {
		return "", &argError{"dns server"}
	}
	if !isValidLimit(s.client.apiKey, resolve.Limit) {
		return "", &argError{"limit"}
	}

	body, err := newJSONReader(resolve)
	if err != nil {
		return "", err
	}
	u := s.client.BasePath + "/run/dns-resolve"
	req, _ := http.NewRequest("POST", u, body)
	req = req.WithContext(ctx)
	var raw struct {
		Error string
		ID    string `json:"id"`
	}
	if err = s.client.do(req, &raw); err != nil {
		return "", err
	}
	if raw.Error != "" {
		return "", errors.New(raw.Error)
	}
	return TestID(raw.ID), nil
}

// DNSResolveOutput returns the full DNS resolve output under a test ID.
func (s *RunService) DNSResolveOutput(ctx context.Context, resolveID TestID) (*DNSTestOutput, error) {
	u := s.client.BasePath + "/run/dns-resolve/" + string(resolveID)
	req, _ := http.NewRequest("GET", u, nil)
	var v *DNSTestOutput
	err := s.client.do(req, &v)
	return v, err
}

// Curl runs a curl request.
func (s *RunService) Curl(ctx context.Context, curl *CurlRequest) (TestID, error) {
	if !isValidTarget(curl.Target) {
		return "", &argError{"target"}
	}
	if !isValidLimit(s.client.apiKey, curl.Limit) {
		return "", &argError{"limit"}
	}

	body, err := newJSONReader(curl)
	if err != nil {
		return "", err
	}
	u := s.client.BasePath + "/run/curl"
	req, _ := http.NewRequest("POST", u, body)
	req = req.WithContext(ctx)
	var raw struct {
		Error string
		ID    string `json:"id"`
	}
	if err = s.client.do(req, &raw); err != nil {
		return "", err
	}
	if raw.Error != "" {
		return "", errors.New(raw.Error)
	}
	return TestID(raw.ID), nil
}

// CurlOutput returns the full curl output under a test ID.
func (s *RunService) CurlOutput(ctx context.Context, curlID TestID) (*RunOutput, error) {
	u := s.client.BasePath + "/run/curl/" + string(curlID)
	req, _ := http.NewRequest("GET", u, nil)
	var v *RunOutput
	err := s.client.do(req, &v)
	return v, err
}

// IsFinished returns a value indicating whether the run result is
// complete or not.
func (r *RunResult) IsFinished() bool {
	if v, ok := r.Finished.(bool); ok {
		return v
	}
	if v, ok := r.Finished.(string); ok {
		return v == "true"
	}
	return false
}

// IsFinished returns a value indicating whether the whole output is
// complete or not.
func (o *RunOutput) IsFinished() bool {
	return o.Finished == "true"
}

// IsFinished returns a value indicating whether the whole output is
// complete or not.
func (o *DNSTestOutput) IsFinished() bool {
	return o.Finished == "true"
}

// PerfOutput returns the unmarshalled output for DNS perf requests.
func (r *DNSTestResult) PerfOutput() string {
	var o string
	if err := json.Unmarshal(r.Output, &o); err != nil {
		o = "-"
	}
	return o
}

// ResolveOutput returns the unmarshalled output for DNS resolve requests.
func (r *DNSTestResult) ResolveOutput() []string {
	var o []string
	if err := json.Unmarshal(r.Output, &o); err == nil {
		return o
	}
	var o2 string
	if err := json.Unmarshal(r.Output, &o2); err != nil {
		return []string{"-"}
	}
	return strings.Split(o2, "\n")
}

// isValidTarget checks if a string is a valid target, i.e., a public
// domain name or an IP address.
func isValidTarget(s string) bool {
	if s == "" {
		return false
	}
	if ip := net.ParseIP(s); ip != nil {
		return true
	}
	// Assume domain name and require at least one level above TLD
	i := strings.LastIndex(s, ".")
	if i == -1 || len(s)-1 == i {
		return false
	}
	// TLD may not start with a number
	if c := s[i+1]; c >= '0' && c <= '9' {
		return false
	}
	return true
}

// isValidLimit retruns a value indicating whether the limit is valid,
// e.g., for requests without an API key the limit is capped at 20.
func isValidLimit(apiKey string, limit int) bool {
	return apiKey != "" || limit <= freeMaxNodeCap
}

func (s *RunService) doPostRunRequest(ctx context.Context, path string, runReq *RunRequest) (TestID, error) {
	if !isValidTarget(runReq.Target) {
		return "", &argError{"target"}
	}
	if !isValidLimit(s.client.apiKey, runReq.Limit) {
		return "", &argError{"limit"}
	}

	body, err := newJSONReader(runReq)
	if err != nil {
		return "", err
	}
	u := s.client.BasePath + path
	req, _ := http.NewRequest("POST", u, body)
	req = req.WithContext(ctx)
	var raw struct {
		Error string
		ID    string `json:"id"`
	}
	if err = s.client.do(req, &raw); err != nil {
		return "", err
	}
	if raw.Error != "" {
		return "", errors.New(raw.Error)
	}
	return TestID(raw.ID), nil
}

func (s *RunService) doGetRunOutput(ctx context.Context, path string, testID TestID) (*RunOutput, error) {
	u := s.client.BasePath + path + string(testID)
	req, _ := http.NewRequest("GET", u, nil)
	var v *RunOutput
	err := s.client.do(req, &v)
	return v, err
}
