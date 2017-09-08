package cmd

import (
	"context"

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

	pingLimit int
)

func initPingCmd() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().IntVarP(&pingLimit, "limit", "L", 1, "The limit")
}

func runPing(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	c, err := perfops.NewClient(perfops.WithAPIKey(apiKey))
	if err != nil {
		return err
	}
	return internal.RunTest(ctx, c, args[0], from, pingLimit, func(ctx context.Context, c *perfops.Client, req *perfops.RunRequest) (perfops.TestID, error) {
		return c.Run.Ping(ctx, req)
	}, func(ctx context.Context, c *perfops.Client, pingID perfops.TestID) (*perfops.RunOutput, error) {
		return c.Run.PingOutput(ctx, pingID)
	})
}
