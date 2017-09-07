package perfops

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
)

type (
	// RunService defines the interface for the run API
	RunService service

	// TestID represents the ID of an MTR or ping test.
	TestID string

	// RunRequest represents the parameters for a ping request.
	RunRequest struct {
		// Target name
		Target string `json:"target"`
		// List of nodes ids, comma separated
		Nodes string `json:"nodes,omitempty"`
		// Countries names, comma separated
		Location string `json:"location,omitempty"`
		// Max number of nodes
		Limit int `json:"limit,omitempty"`
	}

	// RunResult represents the result of an MTR or ping run.
	RunResult struct {
		Node   *Node  `json:"node,omitempty"`
		Output string `json:"output,omitempty"`
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

	// DNSResolveRequest represents the parameters for a DNS resolve request.
	DNSResolveRequest struct {
		Target    string   `json:"target,omitempty"`
		Param     string   `json:"param,omitempty"`
		DNSServer string   `json:"dnsServer,omitempty"`
		Nodes     []string `json:"nodes,omitempty"`
		Location  string   `json:"location,omitempty"`
	}

	// DNSResolveResult represents the result of a DNS resolve output.
	DNSResolveResult struct {
		DNSServer string   `json:"dnsServer,omitempty"`
		Output    []string `json:"output,omitempty"`
		Node      *Node    `json:"node,omitempty"`
	}

	// DNSResolveItem respresents an item of a DNS resolve output.
	DNSResolveItem struct {
		ID     string            `json:"id,omitempty"`
		Result *DNSResolveResult `json:"result,omitempty"`
	}

	// DNSResolveOutput represents the response of a DNS resolve output call.
	DNSResolveOutput struct {
		ID        string            `json:"id,omitempty"`
		Requested string            `json:"requested,omitempty"`
		Finished  string            `json:"finished"`
		Items     []*DNSResolveItem `json:"items,omitempty"`
	}
)

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

// DNSResolve resolves a DNS record.
func (s *RunService) DNSResolve(ctx context.Context, resolve *DNSResolveRequest) (TestID, error) {
	if !isValidTarget(resolve.Target) {
		return "", errors.New("target invalid")
	}
	if resolve.Param == "" {
		return "", errors.New("param invalid")
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
func (s *RunService) DNSResolveOutput(ctx context.Context, resolveID TestID) (*DNSResolveOutput, error) {
	u := s.client.BasePath + "/run/dns-resolve/" + string(resolveID)
	req, _ := http.NewRequest("GET", u, nil)
	var v *DNSResolveOutput
	err := s.client.do(req, &v)
	return v, err
}

// IsFinished returns a value indicating whether the whole output is
// complete or not.
func (o *RunOutput) IsFinished() bool {
	return o.Finished == "true"
}

// IsFinished returns a value indicating whether the whole output is
// complete or not.
func (o *DNSResolveOutput) IsFinished() bool {
	return o.Finished == "true"
}

// isValidTarget checks if a string is a valid target, i.e., a public
// domain name or an IP address.
func isValidTarget(s string) bool {
	if len(s) == 0 {
		return false
	}
	if ip := net.ParseIP(s); ip != nil {
		return true
	}
	// Assume domain name and require at least one level above TLD
	i := strings.LastIndex(s, ".")
	return i > 0 && len(s)-i > 1
}

func (s *RunService) doPostRunRequest(ctx context.Context, path string, runReq *RunRequest) (TestID, error) {
	if !isValidTarget(runReq.Target) {
		return "", errors.New("target invalid")
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
