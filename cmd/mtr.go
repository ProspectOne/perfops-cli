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
	mtrCmd = &cobra.Command{
		Use:   "mtr [target]",
		Short: "Run a MTR test on target",
		Long:  `Run a MTR test on target.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newPerfOpsClient()
			if err != nil {
				return err
			}
			return runMTR(c, args[0], from, mtrLimit)
		},
	}

	mtrLimit int
)

func initMTRCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(mtrCmd)
	mtrCmd.Flags().IntVarP(&mtrLimit, "limit", "L", 1, "The limit")
}

func runMTR(c *perfops.Client, target, from string, limit int) error {
	ctx := context.Background()
	return internal.RunTest(ctx, target, from, limit, c.Run.MTR, c.Run.MTROutput)
}
