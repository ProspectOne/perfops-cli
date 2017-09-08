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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newPerfOpsClient()
			if err != nil {
				return err
			}
			return runPing(c, args[0], from, pingLimit)
		},
	}

	pingLimit int
)

func initPingCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(pingCmd)
	pingCmd.Flags().IntVarP(&pingLimit, "limit", "L", 1, "The limit")
}

func runPing(c *perfops.Client, target, from string, limit int) error {
	ctx := context.Background()
	return internal.RunTest(ctx, target, from, limit, c.Run.Ping, c.Run.PingOutput)
}
