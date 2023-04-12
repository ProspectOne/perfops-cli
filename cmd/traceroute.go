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

	"github.com/spf13/cobra"

	"github.com/ProspectOne/perfops-cli/cmd/internal"
	"github.com/ProspectOne/perfops-cli/perfops"
)

var (
	tracerouteCmd = &cobra.Command{
		Use:     "traceroute [target]",
		Short:   "Run a traceroute test on a domain name or IP address",
		Long:    `Run a traceroute test on a target, e.g., google.com or 8.8.8.8.`,
		Example: `perfops traceroute --from "New York" google.com`,
		Args:    requireTarget(),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newPerfOpsClient()
			if err != nil {
				return err
			}
			return chkRunError(runTraceroute(c, args[0], from, nodeIDs, tracerouteLimit, tracerouteIpv6))
		},
	}

	tracerouteLimit int
	tracerouteIpv6  bool
)

func initTracerouteCmd(parentCmd *cobra.Command) {
	addCommonFlags(tracerouteCmd)
	tracerouteCmd.Flags().IntVarP(&tracerouteLimit, "limit", "L", 1, "The maximum number of nodes to use")
	tracerouteCmd.Flags().BoolVarP(&tracerouteIpv6, "ipv6", "6", false, "Use IPv6")
	parentCmd.AddCommand(tracerouteCmd)
}

func runTraceroute(c *perfops.Client, target, from string, nodeIDs []int, limit int, ipv6 bool) error {
	ctx := context.Background()
	ipversion := 4
	if ipv6 {
		ipversion = 6
	}
	return internal.RunTest(ctx, target, from, nodeIDs, limit, ipversion, debug, outputJSON, c.Run.Traceroute, c.Run.TracerouteOutput)
}
