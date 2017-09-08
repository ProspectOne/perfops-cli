package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/ProspectOne/perfops-cli/cmd/internal"
	"github.com/ProspectOne/perfops-cli/perfops"
)

var (
	latencyCmd = &cobra.Command{
		Use:   "latency [target]",
		Short: "Run a latency test on target",
		Long:  `Run a latency test on target.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newPerfOpsClient()
			if err != nil {
				return err
			}
			return runLatency(c, args[0], from, latencyLimit)
		},
	}

	latencyLimit int
)

func initLatencyCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(latencyCmd)
	latencyCmd.Flags().IntVarP(&latencyLimit, "limit", "L", 1, "The limit")
}

func runLatency(c *perfops.Client, target, from string, limit int) error {
	ctx := context.Background()
	return internal.RunTest(ctx, target, from, limit, c.Run.Latency, c.Run.LatencyOutput)
}
