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
		RunE:  runLatency,
	}

	latencyFrom  string
	latencyLimit int
)

func initlatencyCmd() {
	rootCmd.AddCommand(latencyCmd)
	latencyCmd.Flags().StringVarP(&latencyFrom, "from", "F", "", "A continent, region (e.g eastern europe), country, US state or city")
	latencyCmd.Flags().IntVarP(&latencyLimit, "limit", "L", 1, "The limit")
}

func runLatency(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	c, err := perfops.NewClient(perfops.WithAPIKey(apiKey))
	if err != nil {
		return err
	}
	return internal.RunTest(ctx, c, args[0], latencyFrom, latencyLimit, func(ctx context.Context, c *perfops.Client, req *perfops.RunRequest) (perfops.TestID, error) {
		return c.Run.Latency(ctx, req)
	}, func(ctx context.Context, c *perfops.Client, latencyID perfops.TestID) (*perfops.RunOutput, error) {
		return c.Run.LatencyOutput(ctx, latencyID)
	})
}
