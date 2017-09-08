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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newPerfOpsClient()
			if err != nil {
				return err
			}
			return runMTR(c, args[0], from, mtrLimit)
		},
	}

	mtrLimit int
)

func initMTRCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(mtrCmd)
	mtrCmd.Flags().IntVarP(&mtrLimit, "limit", "L", 1, "The limit")
}

func runMTR(c *perfops.Client, target, from string, limit int) error {
	ctx := context.Background()
	return internal.RunTest(ctx, target, from, limit, c.Run.MTR, c.Run.MTROutput)
}
