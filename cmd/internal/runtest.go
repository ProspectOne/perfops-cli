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

package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/ProspectOne/perfops-cli/perfops"
	"github.com/gosuri/uilive"
)

type (
	terminalWriter interface {
		io.Writer
		Flush() error
	}

	// Formatter formats the run output.
	Formatter struct {
		printID bool
		s       *Spinner
		w       terminalWriter
	}

	// RunOutputResult collects the RunOutput and its error, if any, from
	// async calls.
	RunOutputResult struct {
		mu     sync.Mutex
		output *perfops.RunOutput
		err    error
	}

	runFunc       func(ctx context.Context, req *perfops.RunRequest) (perfops.TestID, error)
	runOutputFunc func(ctx context.Context, pingID perfops.TestID) (*perfops.RunOutput, error)
)

// RunTest runs an MTR or ping testm retrives its output and presents it to the user.
func RunTest(ctx context.Context, target, location string, nodeIDs []int, limit int, debug, outputJSON bool, runTest runFunc, runOutput runOutputFunc) error {
	runReq := &perfops.RunRequest{
		Target:   target,
		Location: location,
		Nodes:    nodeIDs,
		Limit:    limit,
	}

	f := NewFormatter(debug && !outputJSON)
	f.StartSpinner()
	testID, err := runTest(ctx, runReq)
	f.StopSpinner()
	if err != nil {
		return err
	}
	res := &RunOutputResult{}
	go func() {
		for {
			select {
			case <-time.After(200 * time.Millisecond):
			}
			output, err := runOutput(ctx, testID)
			res.SetOutput(output, err)
			if err != nil {
				break
			}
		}
	}()

	f.StartSpinner()
	var o *perfops.RunOutput
	for {
		select {
		case <-time.After(100 * time.Millisecond):
		}
		if o, err = res.Output(); err != nil {
			return err
		}
		if !outputJSON && o != nil {
			f.StopSpinner()
			PrintOutput(f, o)
		}
		if o != nil && o.IsFinished() {
			break
		}
	}
	if outputJSON {
		f.StopSpinner()
		PrintOutputJSON(o)
	}
	return nil
}

// PrintOutput prints run items that have been data.
func PrintOutput(f *Formatter, output *perfops.RunOutput) {
	if f.printID {
		f.Printf("Test ID: %v\n", output.ID)
	}
	for _, item := range output.Items {
		r := item.Result
		n := r.Node
		if item.Result.Message == "" {
			o := r.Output
			if o == "-2" {
				o = "The command timed-out. It either took too long to execute or we could not connect to your target at all."
			}
			f.Printf("Node%d, %s, %s\n%s\n", n.ID, n.City, n.Country.Name, o)
		} else if r.Message != "NO DATA" {
			f.Printf("Node%d, %s, %s\n%s\n", n.ID, n.City, n.Country.Name, r.Message)
		}
	}
	spinner := f.s.Step()
	if !output.IsFinished() {
		f.Printf("%s\n", spinner)
	}
	f.w.Flush()
}

// PrintOutputJSON marshals the output into JSON and prints the JSON.
func PrintOutputJSON(output interface{}) error {
	b, err := json.Marshal(output)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

// NewFormatter returns a new Formatter
func NewFormatter(printID bool) *Formatter {
	f := &Formatter{
		printID: printID,
		w:       uilive.New(),
		s:       NewSpinner(),
	}
	return f
}

// StartSpinner starts the spinner.
func (f *Formatter) StartSpinner() {
	f.s.Start()
}

// StopSpinner stops the spinner.
func (f *Formatter) StopSpinner() {
	f.s.Stop()
}

// Printf prints the arguments to the formatters writer.
func (f *Formatter) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(f.w, format, a...)
}

// SetOutput sets the output and the error.
func (g *RunOutputResult) SetOutput(o *perfops.RunOutput, err error) {
	g.mu.Lock()
	g.output = o
	g.err = err
	g.mu.Unlock()
}

// Output retruns the output and the error.
func (g *RunOutputResult) Output() (*perfops.RunOutput, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.output, g.err
}
