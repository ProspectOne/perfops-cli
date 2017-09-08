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
		Use:   "dns-resolve [target]",
		Short: "Resolve a DNS record on target",
		Long:  `Resolve a DNS record on target.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newPerfOpsClient()
			if err != nil {
				return err
			}
			return runDNSResolve(c, args[0], dnsResolveParam, dnsResolveDNSServer, from)
		},
	}

	dnsResolveParam     string
	dnsResolveDNSServer string
)

func initDNSResolveCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(dnsResolveCmd)
	dnsResolveCmd.Flags().StringVarP(&dnsResolveParam, "param", "P", "", "The DNS query type. On of: A, AAAA, CNAME, MX, NAPTR, NS, PTR, SOA, SPF, SRV, TXT.")
	dnsResolveCmd.Flags().StringVarP(&dnsResolveDNSServer, "dns-server", "S", "", "The DNS server to use to query for the test. You can use 127.0.0.1 to use the local resolver for location based benchmarking.")
}

func runDNSResolve(c *perfops.Client, target, param, dnsServer, from string) error {
	ctx := context.Background()
	dnsResolveReq := &perfops.DNSResolveRequest{
		Target:    target,
		Param:     param,
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
