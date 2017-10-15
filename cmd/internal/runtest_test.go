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
	"errors"
	"fmt"
	"testing"

	"github.com/ProspectOne/perfops-cli/perfops"
)

func TestRunTest(t *testing.T) {
	runErr := errors.New("run")
	outputErr := errors.New("output")
	testCases := map[string]struct {
		run    runFunc
		output runOutputFunc
		err    error
	}{
		"run failed": {
			func(ctx context.Context, req *perfops.RunRequest) (perfops.TestID, error) {
				return perfops.TestID(""), runErr
			},
			nil,
			runErr,
		},
		"output failed": {
			func(ctx context.Context, req *perfops.RunRequest) (perfops.TestID, error) {
				return perfops.TestID("test-123"), nil
			},
			func(ctx context.Context, pingID perfops.TestID) (*perfops.RunOutput, error) {
				return nil, outputErr
			},
			outputErr,
		},
		"succeeded": {
			func(ctx context.Context, req *perfops.RunRequest) (perfops.TestID, error) {
				return perfops.TestID("test-123"), nil
			},
			func(ctx context.Context, pingID perfops.TestID) (*perfops.RunOutput, error) {
				var output *perfops.RunOutput
				err := json.Unmarshal([]byte(`{"id": "9072a72f762b876525ca4c9153af9983","items": [{"id": "edca088e43bde5453b961f6210723157","result": {"output": "Start: Thu Jul 27 15:59:05 2017                Loss%   Snt   Last   Avg  Best  Wrst StDev\n  1.|-- 172.18.0.1                 0.0%     2    0.0   0.1   0.0   0.1   0.0\n  2.|-- 10.0.2.2                   0.0%     2    0.2   0.2   0.2   0.2   0.0\n  3.|-- 192.168.0.1                0.0%     2    1.3   1.5   1.3   1.6   0.0\n  4.|-- ???                       100.0     2    0.0   0.0   0.0   0.0   0.0\n  5.|-- 80.81.194.168              0.0%     2   25.0  23.4  21.8  25.0   2.0\n  6.|-- 80.81.194.52               0.0%     2   41.4  41.4  41.4  41.4   0.0\n  7.|-- 104.44.80.143              0.0%     2   40.9  40.5  40.0  40.9   0.0\n  8.|-- ???                       100.0     2    0.0   0.0   0.0   0.0   0.0\n  9.|-- ???                       100.0     2    0.0   0.0   0.0   0.0   0.0\n 10.|-- ???                       100.0     2    0.0   0.0   0.0   0.0   0.0\n 11.|-- 13.107.21.200              0.0%     2   39.8  40.1  39.8  40.3   0.0\n","node": {"id": 5,"latitude": 50.110781326572834,"longitude": 8.68984222412098,"country": {"id": 116,"name": "Germany","continent": {"id": 3,"name": "Europe","iso": "EU"},"iso": "DE","iso_numeric": "276"},"city": "Frankfurt","sub_region": "Western Europe"}}}],"requested": "bing.com","finished": "true"}`), &output)
				return output, err
			},
			nil,
		},
	}
	ctx := context.Background()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := RunTest(ctx, "target", "location", []int{}, 1, false, false, tc.run, tc.output)
			if err != tc.err {
				t.Fatalf("expected %v; got %v", tc.err, err)
			}
		})
	}
}

func TestPrintPartialOutput(t *testing.T) {
	testCases := map[string]struct {
		output     func() *perfops.RunOutput
		printedIDs map[string]bool
		exp        string
	}{
		"all": {
			func() *perfops.RunOutput {
				var o *perfops.RunOutput
				json.Unmarshal([]byte(`{"id":"706fc55e3377104da01f05569e35a30b","items":[{"id":"bba072473bd034d432df03730e116686","result":{"output":"121","node":{"id":27,"latitude":22.28512548314,"longitude":114.17507171631,"country":{"id":195,"name":"Hong Kong","continent":{"id":2,"name":"Asia","iso":"AS"},"iso":"HK","iso_numeric":"344"},"city":"Hong Kong","sub_region":"Eastern Asia"},"time":1508061924.088372}}],"requested":"sendergram.com","finished":"true","elapsedTime":0.66500000000000004}`), &o)
				return o
			},
			map[string]bool{},
			"Node27, Hong Kong, Hong Kong\n121\n",
		},
		"partially printed": {
			func() *perfops.RunOutput {
				var o *perfops.RunOutput
				json.Unmarshal([]byte(`{"id":"706fc55e3377104da01f05569e35a30b","items":[{"id":"bba072473bd034d432df03730e116686","result":{"output":"121","node":{"id":27,"latitude":22.28512548314,"longitude":114.17507171631,"country":{"id":195,"name":"Hong Kong","continent":{"id":2,"name":"Asia","iso":"AS"},"iso":"HK","iso_numeric":"344"},"city":"Hong Kong","sub_region":"Eastern Asia"},"time":1508061924.088372}}],"requested":"sendergram.com","finished":"true","elapsedTime":0.66500000000000004}`), &o)
				return o
			},
			map[string]bool{"bba072473bd034d432df03730e116686": true},
			"",
		},
		"timeout": {
			func() *perfops.RunOutput {
				var o *perfops.RunOutput
				json.Unmarshal([]byte(`{"id":"706fc55e3377104da01f05569e35a30b","items":[{"id":"bba072473bd034d432df03730e116686","result":{"output":"-2","node":{"id":27,"latitude":22.28512548314,"longitude":114.17507171631,"country":{"id":195,"name":"Hong Kong","continent":{"id":2,"name":"Asia","iso":"AS"},"iso":"HK","iso_numeric":"344"},"city":"Hong Kong","sub_region":"Eastern Asia"},"time":1508061924.088372}}],"requested":"sendergram.com","finished":"true","elapsedTime":0.66500000000000004}`), &o)
				return o
			},
			map[string]bool{},
			"Node27, Hong Kong, Hong Kong\nThe command timed-out. It either took too long to execute or we could not connect to your target at all.\n",
		},
	}

	var b bytes.Buffer
	printf := func(format string, a ...interface{}) (n int, err error) {
		return fmt.Fprintf(&b, format, a...)
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			b.Reset()
			PrintPartialOutput(printf, tc.output(), tc.printedIDs)
			got := b.String()
			if got != tc.exp {
				t.Fatalf("expected %#v; got %#v", tc.exp, got)
			}
		})
	}
}
