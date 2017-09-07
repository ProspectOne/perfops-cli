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
		RunE:  runDNSResolve,
	}

	dnsResolveFrom      string
	dnsResolveParam     string
	dnsResolveDNSServer string
)

func initDNSResolveCmd() {
	rootCmd.AddCommand(dnsResolveCmd)
	dnsResolveCmd.Flags().StringVarP(&dnsResolveFrom, "from", "F", "", "A continent, region (e.g eastern europe), country, US state or city")
	dnsResolveCmd.Flags().StringVarP(&dnsResolveParam, "param", "P", "", "The DNS query type. On of: A, AAAA, CNAME, MX, NAPTR, NS, PTR, SOA, SPF, SRV, TXT.")
	dnsResolveCmd.Flags().StringVarP(&dnsResolveDNSServer, "dns-server", "S", "", "The DNS server to use to query for the test. You can use 127.0.0.1 to use the local resolver for location based benchmarking.")
}

func runDNSResolve(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	c, err := perfops.NewClient(perfops.WithAPIKey(apiKey))
	if err != nil {
		return err
	}
	dnsResolveReq := &perfops.DNSResolveRequest{
		Target:    args[0],
		Param:     dnsResolveParam,
		DNSServer: dnsResolveDNSServer,
		Location:  dnsResolveFrom,
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
