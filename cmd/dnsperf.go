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
	dnsPerfCmd = &cobra.Command{
		Use:     "dnsperf [target]",
		Short:   "Find the time it takes to resolve a DNS record on a target",
		Long:    `Find the time it takes to resolve a DNS record on a target, e.g., google.com.`,
		Example: `perfops dnsperf bing.com`,
		Args:    requireTarget(),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newPerfOpsClient()
			if err != nil {
				return err
			}
			return chkRunError(runDNSPerf(c, args[0], dnsPerfDNSServer, from, nodeIDs, dnsPerfLimit, dnsPerfIpv6))
		},
	}

	dnsPerfDNSServer string
	dnsPerfLimit     int
	dnsPerfIpv6      bool
)

func initDNSPerfCmd(parentCmd *cobra.Command) {
	addCommonFlags(dnsPerfCmd)

	dnsPerfCmd.Flags().StringVarP(&dnsPerfDNSServer, "dns-server", "S", "", "The DNS server to use to query for the test. You can use 127.0.0.1 to use the local resolver for location based benchmarking.")
	dnsPerfCmd.Flags().IntVarP(&dnsPerfLimit, "limit", "L", 1, "The maximum number of nodes to use")
	dnsPerfCmd.Flags().BoolVarP(&dnsPerfIpv6, "ipv6", "6", false, "Use IPv6")

	parentCmd.AddCommand(dnsPerfCmd)
}

func runDNSPerf(c *perfops.Client, target, dnsServer, from string, nodeIDs []int, limit int, ipv6 bool) error {
	ctx := context.Background()
	ipversion := 4
	if ipv6 {
		ipversion = 6
	}
	dnsPerfReq := &perfops.DNSPerfRequest{
		Target:    target,
		DNSServer: dnsServer,
		Location:  from,
		Nodes:     nodeIDs,
		Limit:     limit,
		IPVersion: ipversion,
	}

	spinner := internal.NewSpinner()
	fmt.Println("")
	spinner.Start()

	testID, err := c.Run.DNSPerf(ctx, dnsPerfReq)
	spinner.Stop()
	if err != nil {
		return err
	}

	if debug && !outputJSON {
		fmt.Printf("Test ID: %v\n", testID)
	}

	var output *perfops.DNSTestOutput
	printedIDs := map[string]bool{}
	for {
		spinner.Start()
		select {
		case <-time.After(500 * time.Millisecond):
		}

		output, err = c.Run.DNSPerfOutput(ctx, testID)
		spinner.Stop()
		if err != nil {
			return err
		}

		if !outputJSON {
			printPartialDNSOutput(fmt.Printf, output, printedIDs, func(r *perfops.DNSTestResult) string {
				return r.PerfOutput()
			})
		}
		if output.IsFinished() {
			break
		}
	}
	if outputJSON {
		internal.PrintOutputJSON(output)
	}
	return nil
}
