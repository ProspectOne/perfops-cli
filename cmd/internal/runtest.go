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
	"bytes"
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

		mu  sync.Mutex
		buf bytes.Buffer
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

	if outputJSON {
		f.StartSpinner()
	}
	var o *perfops.RunOutput
	for {
		select {
		case <-time.After(50 * time.Millisecond):
		}
		if o, err = res.Output(); err != nil {
			return err
		}
		if !outputJSON && o != nil {
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
	spinner := f.s.Step()
	if !output.IsFinished() {
		f.Printf("%s", spinner)
		if len(output.Items) > 1 {
			finished := 0
			for _, item := range output.Items {
				if item.Result.IsFinished() {
					finished++
				}
			}
			f.Printf(" %d/%d", finished, len(output.Items))
		}
		f.Printf("\n")
	}
	for _, item := range output.Items {
		r := item.Result
		n := r.Node
		if item.Result.Message == "" {
			o := r.Output
			if o == "-2" {
				o = "The command timed-out. It either took too long to execute or we could not connect to your target at all."
			}
			f.Printf("Node%d, AS%d, %s, %s\n%s\n", n.ID, n.AsNumber, n.City, n.Country.Name, o)
		} else if r.Message != "NO DATA" {
			f.Printf("Node%d, AS%d, %s, %s\n%s\n", n.ID, n.AsNumber, n.City, n.Country.Name, r.Message)
		}
		if !item.Result.IsFinished() {
			f.Printf("%s\n", spinner)
		}
	}
	f.Flush(!output.IsFinished())
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
	f.mu.Lock()
	defer f.mu.Unlock()
	return fmt.Fprintf(&f.buf, format, a...)
}

// Flush writes to w and resets the buffer. It should be called after the
// last call to Printf to ensure that any data buffered in the Writer is
// written to output.
// Any incomplete escape sequence at the end is considered complete for
// formatting purposes.
// An error is returned if the contents of the buffer cannot be written
// to the underlying output stream
func (f *Formatter) Flush(limit bool) error {
	cols, rows := termSize()

	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.buf.Bytes()) == 0 {
		return nil
	}
	defer f.buf.Reset()

	termStartOfRow(f.w)
	if limit {
		out := string(f.buf.Bytes())
		width := 0
		lines := 0
		for _, r := range out {
			width++
			if r == '\n' || width > cols {
				width = 0
				lines++
				if lines == rows {
					break
				}
				if _, err := f.w.Write([]byte("\n")); err != nil {
					return err
				}
			} else if _, err := f.w.Write([]byte(string(r))); err != nil {
				return err
			}
		}
	} else if _, err := f.w.Write(f.buf.Bytes()); err != nil {
		return err
	}
	return f.w.Flush()
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
