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
)

// MTR runs an MTR test.
func (s *RunService) MTR(ctx context.Context, mtr *RunRequest) (TestID, error) {
	if !isValidTarget(mtr.Target) {
		return "", errors.New("target invalid")
	}
	body, err := newJSONReader(mtr)
	if err != nil {
		return "", err
	}
	u := s.client.BasePath + "/run/mtr"
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

// MTROutput returns the full MTR output under a test ID.
func (s *RunService) MTROutput(ctx context.Context, mtrID TestID) (*RunOutput, error) {
	u := s.client.BasePath + "/run/mtr/" + string(mtrID)
	req, _ := http.NewRequest("GET", u, nil)
	var v *RunOutput
	if err := s.client.do(req, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Ping runs a ping test.
func (s *RunService) Ping(ctx context.Context, ping *RunRequest) (TestID, error) {
	if !isValidTarget(ping.Target) {
		return "", errors.New("target invalid")
	}
	body, err := newJSONReader(ping)
	if err != nil {
		return "", err
	}
	u := s.client.BasePath + "/run/ping"
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

// PingOutput returns the full ping output under a test ID.
func (s *RunService) PingOutput(ctx context.Context, pingID TestID) (*RunOutput, error) {
	u := s.client.BasePath + "/run/ping/" + string(pingID)
	req, _ := http.NewRequest("GET", u, nil)
	var v *RunOutput
	if err := s.client.do(req, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IsFinished returns a value indicating whether the whole output is
// complete or not.
func (o *RunOutput) IsFinished() bool {
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
