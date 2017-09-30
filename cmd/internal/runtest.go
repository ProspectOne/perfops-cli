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
	"time"

	"github.com/ProspectOne/perfops-cli/perfops"
)

type (
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

	spinner := NewSpinner()
	fmt.Println("")
	spinner.Start()

	testID, err := runTest(ctx, runReq)
	spinner.Stop()
	if err != nil {
		return err
	}

	if debug && !outputJSON {
		fmt.Printf("Test ID: %v\n", testID)
	}

	var output *perfops.RunOutput
	printedIDs := map[string]bool{}
	for {
		spinner.Start()
		select {
		case <-time.After(500 * time.Millisecond):
		}

		output, err = runOutput(ctx, testID)
		spinner.Stop()
		if err != nil {
			return err
		}

		if !outputJSON {
			PrintPartialOutput(output, printedIDs)
		}
		if output.IsFinished() {
			break
		}
	}
	if outputJSON {
		PrintOutputJSON(output)
	}
	return nil
}

// PrintPartialOutput prints run items that have been data.
func PrintPartialOutput(output *perfops.RunOutput, printedIDs map[string]bool) {
	for _, item := range output.Items {
		if printedIDs[item.ID] {
			continue
		}
		r := item.Result
		n := r.Node
		if item.Result.Message == "" {
			printedIDs[item.ID] = true
			fmt.Printf("Node%d, %s, %s\n%s\n", n.ID, n.City, n.Country.Name, r.Output)
		} else if r.Message != "NO DATA" {
			printedIDs[item.ID] = true
			fmt.Printf("Node%d, %s, %s\n%s\n", n.ID, n.City, n.Country.Name, r.Message)
		}
	}
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
