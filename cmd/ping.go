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
	pingCmd = &cobra.Command{
		Use:     "ping [target]",
		Short:   "Run a ping test on a domain name or IP address",
		Long:    `Run a ping test on a target, e.g., google.com or 8.8.8.8.`,
		Example: `perfops ping --from "eastern europe" google.com`,
		Args:    requireTarget(),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newPerfOpsClient()
			if err != nil {
				return err
			}
			return chkRunError(runPing(c, args[0], from, nodeIDs, pingLimit))
		},
	}

	pingLimit int
)

func initPingCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(pingCmd)
	pingCmd.Flags().IntVarP(&pingLimit, "limit", "L", 1, "The maximum number of nodes to use")
}

func runPing(c *perfops.Client, target, from string, nodeIDs []int, limit int) error {
	ctx := context.Background()
	return internal.RunTest(ctx, target, from, nodeIDs, limit, debug, c.Run.Ping, c.Run.PingOutput)
}
