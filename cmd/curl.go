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

package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/ProspectOne/perfops-cli/cmd/internal"
	"github.com/ProspectOne/perfops-cli/perfops"
)

var (
	curlCmd = &cobra.Command{
		Use:     "curl [target]",
		Short:   "Run a curl test on a domain name or IP address",
		Long:    `Run a curl test on a target, e.g., google.com or 8.8.8.8.`,
		Example: `perfops curl --http2 bing.com`,
		Args:    requireTarget(),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newPerfOpsClient()
			if err != nil {
				return err
			}
			return chkRunError(runCurl(c, args[0], curlHead, curlInsecure, curlHTTP2, from, nodeIDs, curlLimit))
		},
	}

	curlHead     bool
	curlInsecure bool
	curlHTTP2    bool
	curlLimit    int
)

func initCurlCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(curlCmd)
	curlCmd.Flags().BoolVarP(&curlHead, "head", "I", true, "Fetch the headers only")
	curlCmd.Flags().BoolVarP(&curlInsecure, "insecure", "k", false, "Allow curl to proceed for server connections considered insecure")
	curlCmd.Flags().BoolVarP(&curlHTTP2, "http2", "", false, "Use HTTP version 2")
	curlCmd.Flags().IntVarP(&dnsResolveLimit, "limit", "L", 1, "The maximum number of nodes to use")
}

func runCurl(c *perfops.Client, target string, head, insecure, http2 bool, from string, nodeIDs []int, limit int) error {
	ctx := context.Background()
	curlReq := &perfops.CurlRequest{
		Target:   target,
		Head:     head,
		Insecure: insecure,
		HTTP2:    http2,
		Location: from,
		Nodes:    nodeIDs,
		Limit:    limit,
	}

	spinner := internal.NewSpinner()
	fmt.Println("")
	spinner.Start()

	testID, err := c.Run.Curl(ctx, curlReq)
	spinner.Stop()
	if err != nil {
		return err
	}

	if debug {
		fmt.Printf("Test ID: %v\n", testID)
	}

	var output *perfops.RunOutput
	printedIDs := map[string]bool{}
	for {
		spinner.Start()
		select {
		case <-time.After(500 * time.Millisecond):
		}

		output, err = c.Run.CurlOutput(ctx, testID)
		spinner.Stop()
		if err != nil {
			return err
		}

		internal.PrintPartialOutput(output, printedIDs)
		if output.IsFinished() {
			break
		}
	}
	return nil
}
