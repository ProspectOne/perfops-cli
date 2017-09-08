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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newPerfOpsClient()
			if err != nil {
				return err
			}
			return runTraceroute(c, args[0], from, tracerouteLimit)
		},
	}

	tracerouteLimit int
)

func initTracerouteCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(tracerouteCmd)
	tracerouteCmd.Flags().IntVarP(&tracerouteLimit, "limit", "L", 1, "The limit")
}

func runTraceroute(c *perfops.Client, target, from string, limit int) error {
	ctx := context.Background()
	return internal.RunTest(ctx, target, from, limit, c.Run.Traceroute, c.Run.TracerouteOutput)
}
