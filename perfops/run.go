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

	// Ping represents the parameters for a ping request.
	Ping struct {
		// Target name
		Target string `json:"target"`
		// List of nodes ids, comma separated
		Nodes string `json:"nodes,omitempty"`
		// Countries names, comma separated
		Location string `json:"location,omitempty"`
		// Max number of nodes
		Limit int `json:"limit,omitempty"`
	}

	// PingID represents the ID of a ping test.
	PingID string

	// Continent contains information about a continent.
	Continent struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		ISO  string `json:"iso"`
	}

	// Country contains information about a country.
	Country struct {
		ID         int        `json:"id"`
		Name       string     `json:"name"`
		ISO        string     `json:"iso"`
		ISONumeric string     `json:"iso_numeric"`
		Continent  *Continent `json:"continent,omitempty"`
	}

	// Node contains informatin about a test node.
	Node struct {
		ID        int      `json:"id"`
		Latitude  float64  `json:"latitude"`
		Longitude float64  `json:"longitude"`
		City      string   `json:"city"`
		SubRegion string   `json:"sub_region"`
		Country   *Country `json:"country,omitempty"`
	}

	// PingResult represents the result of a ping.
	PingResult struct {
		Node   *Node  `json:"node,omitempty"`
		Output string `json:"output,omitempty"`
	}

	// PingItem represents an item of a ping output.
	PingItem struct {
		ID     string      `json:"id,omitempty"`
		Result *PingResult `json:"result,omitempty"`
	}

	// PingOutput represents the response of ping output calls.
	PingOutput struct {
		ID        string      `json:"id,omitempty"`
		Requested string      `json:"requested,omitempty"`
		Finished  string      `json:"finished"`
		Items     []*PingItem `json:"items,omitempty"`
	}
)

// Ping runs a ping test.
func (s *RunService) Ping(ctx context.Context, ping *Ping) (PingID, error) {
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
	return PingID(raw.ID), nil
}

// PingOutput returns the full ping output under a test ID.
func (s *RunService) PingOutput(ctx context.Context, pingID PingID) (*PingOutput, error) {
	u := s.client.BasePath + "/run/ping/" + string(pingID)
	req, _ := http.NewRequest("GET", u, nil)
	var v *PingOutput
	if err := s.client.do(req, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IsFinished returns a value indicating whether the whole output is
// complete or not.
func (o *PingOutput) IsFinished() bool {
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
