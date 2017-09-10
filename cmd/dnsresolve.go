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
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ProspectOne/perfops-cli/cmd/internal"
	"github.com/ProspectOne/perfops-cli/perfops"
)

var (
	dnsResolveCmd = &cobra.Command{
		Use:   "resolve [target]",
		Short: "Resolve a DNS record on target",
		Long:  `Resolve a DNS record on target.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newPerfOpsClient()
			if err != nil {
				return err
			}
			return chkRunError(runDNSResolve(c, args[0], dnsResolveType, dnsResolveDNSServer, from))
		},
	}

	dnsResolveType      string
	dnsResolveDNSServer string
)

func initDNSResolveCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(dnsResolveCmd)
	dnsResolveCmd.Flags().StringVarP(&dnsResolveType, "type", "T", "", "The DNS query type. On of: A, AAAA, CNAME, MX, NAPTR, NS, PTR, SOA, SPF, SRV, TXT.")
	dnsResolveCmd.Flags().StringVarP(&dnsResolveDNSServer, "dns-server", "S", "", "The DNS server to use to query for the test. You can use 127.0.0.1 to use the local resolver for location based benchmarking.")
}

func runDNSResolve(c *perfops.Client, target, queryType, dnsServer, from string) error {
	ctx := context.Background()
	dnsResolveReq := &perfops.DNSResolveRequest{
		Target:    target,
		Param:     queryType,
		DNSServer: dnsServer,
		Location:  from,
	}

	spinner := internal.NewSpinner()
	fmt.Println("")
	spinner.Start()

	testID, err := c.Run.DNSResolve(ctx, dnsResolveReq)
	if err != nil {
		spinner.Stop()
		return err
	}

	var output *perfops.DNSResolveOutput
	for {
		select {
		case <-time.After(250 * time.Millisecond):
		}

		if output, err = c.Run.DNSResolveOutput(ctx, testID); err != nil {
			spinner.Stop()
			return err
		}
		if output.IsFinished() {
			break
		}
	}

	spinner.Stop()

	for _, item := range output.Items {
		n := item.Result.Node
		fmt.Printf("Node%d, %s, %s\n%s\n", n.ID, n.City, n.Country.Name, strings.Join(item.Result.Output, "\n"))
	}
	return nil
}
