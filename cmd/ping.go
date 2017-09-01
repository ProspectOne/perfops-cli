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
	pingCmd = &cobra.Command{
		Use:   "ping [target]",
		Short: "Run a ping test on target",
		Long:  `Run a ping test on target.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runPing,
	}

	nodes []string
	from  string
	limit int
)

func init() {
	RootCmd.AddCommand(pingCmd)
	pingCmd.Flags().StringSliceVarP(&nodes, "nodes", "N", []string{}, "A list of node IDs")
	pingCmd.Flags().StringVarP(&from, "from", "F", "", "A continent, region (e.g eastern europe), country, US state or city")
	pingCmd.Flags().IntVarP(&limit, "limit", "L", 1, "The limit")
}

func runPing(cmd *cobra.Command, args []string) error {
	ping := &perfops.Ping{
		Target:   args[0],
		Nodes:    strings.Join(nodes, ","),
		Location: from,
		Limit:    limit,
	}

	ctx := context.Background()
	c, err := perfops.NewClient(perfops.WithAPIKey(apiKey))
	if err != nil {
		return err
	}

	spinner := internal.NewSpinner()
	fmt.Println("")
	spinner.Start()

	pingID, err := c.Run.Ping(ctx, ping)
	if err != nil {
		spinner.Stop()
		return err
	}

	var output *perfops.PingOutput
	for {
		select {
		case <-time.After(250 * time.Millisecond):
		}

		if output, err = c.Run.PingOutput(ctx, pingID); err != nil {
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
		fmt.Printf("Node%d, %s, %s\n%s\n", n.ID, n.City, n.Country.Name, item.Result.Output)
	}
	return nil
}
