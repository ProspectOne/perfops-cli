// Copyright 2017 The PerfOps-CLI Authors. All rights reserved.
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

	"github.com/spf13/cobra"

	"github.com/ProspectOne/perfops-cli/cmd/internal"
	"github.com/ProspectOne/perfops-cli/perfops"
)

var (
	latencyCmd = &cobra.Command{
		Use:     "latency [target]",
		Short:   "Run a ICMP latency test on a domain name or IP address",
		Long:    `Run a ICMP latency test on a target, e.g., google.com or 8.8.8.8.`,
		Example: `perfops latency --from europe --limit 9 google.com`,
		Args:    requireTarget(),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newPerfOpsClient()
			if err != nil {
				return err
			}
			return chkRunError(runLatency(c, args[0], from, nodeIDs, latencyLimit))
		},
	}

	latencyLimit int
)

func initLatencyCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(latencyCmd)
	latencyCmd.Flags().IntVarP(&latencyLimit, "limit", "L", 1, "The maximum number of nodes to use")
}

func runLatency(c *perfops.Client, target, from string, nodeIDs []int, limit int) error {
	ctx := context.Background()
	return internal.RunTest(ctx, target, from, nodeIDs, limit, debug, c.Run.Latency, c.Run.LatencyOutput)
}
