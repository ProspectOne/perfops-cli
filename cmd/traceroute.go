package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/ProspectOne/perfops-cli/cmd/internal"
	"github.com/ProspectOne/perfops-cli/perfops"
)

var (
	tracerouteCmd = &cobra.Command{
		Use:   "traceroute [target]",
		Short: "Run a traceroute test on target",
		Long:  `Run a traceroute test on target.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runTraceroute,
	}

	tracerouteLimit int
)

func initTracerouteCmd() {
	rootCmd.AddCommand(tracerouteCmd)
	tracerouteCmd.Flags().IntVarP(&tracerouteLimit, "limit", "L", 1, "The limit")
}

func runTraceroute(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	c, err := perfops.NewClient(perfops.WithAPIKey(apiKey))
	if err != nil {
		return err
	}
	return internal.RunTest(ctx, c, args[0], from, tracerouteLimit, func(ctx context.Context, c *perfops.Client, req *perfops.RunRequest) (perfops.TestID, error) {
		return c.Run.Traceroute(ctx, req)
	}, func(ctx context.Context, c *perfops.Client, tracerouteID perfops.TestID) (*perfops.RunOutput, error) {
		return c.Run.TracerouteOutput(ctx, tracerouteID)
	})
}
