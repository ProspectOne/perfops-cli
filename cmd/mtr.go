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
		RunE:  runMTR,
	}

	mtrLimit int
)

func initMTRCmd() {
	rootCmd.AddCommand(mtrCmd)
	mtrCmd.Flags().IntVarP(&mtrLimit, "limit", "L", 1, "The limit")
}

func runMTR(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	c, err := perfops.NewClient(perfops.WithAPIKey(apiKey))
	if err != nil {
		return err
	}
	return internal.RunTest(ctx, c, args[0], from, mtrLimit, func(ctx context.Context, c *perfops.Client, req *perfops.RunRequest) (perfops.TestID, error) {
		return c.Run.MTR(ctx, req)
	}, func(ctx context.Context, c *perfops.Client, pingID perfops.TestID) (*perfops.RunOutput, error) {
		return c.Run.MTROutput(ctx, pingID)
	})
}
