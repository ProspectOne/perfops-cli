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
	addCommonFlags(curlCmd)

	curlCmd.Flags().BoolVarP(&curlHead, "head", "I", false, "Fetch the headers only")
	curlCmd.Flags().BoolVarP(&curlInsecure, "insecure", "k", false, "Allow curl to proceed for server connections considered insecure")
	curlCmd.Flags().BoolVarP(&curlHTTP2, "http2", "", false, "Use HTTP version 2")
	curlCmd.Flags().IntVarP(&curlLimit, "limit", "L", 1, "The maximum number of nodes to use")

	parentCmd.AddCommand(curlCmd)
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

	f := internal.NewFormatter(debug && !outputJSON)

	f.StartSpinner()
	testID, err := c.Run.Curl(ctx, curlReq)
	f.StopSpinner()
	if err != nil {
		return err
	}

	res := &internal.RunOutputResult{}
	go func() {
		for {
			select {
			case <-time.After(200 * time.Millisecond):
			}
			output, err := c.Run.CurlOutput(ctx, testID)
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
			internal.PrintOutput(f, o)
		}
		if o != nil && o.IsFinished() {
			break
		}
	}
	if outputJSON {
		f.StopSpinner()
		internal.PrintOutputJSON(o)
	}
	return nil
}
