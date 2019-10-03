package cmd

import (
	"context"
	"github.com/ProspectOne/perfops-cli/cmd/internal"
	"github.com/ProspectOne/perfops-cli/perfops"
	"github.com/spf13/cobra"
	"net/http"
)

var (
	citiesCmd = &cobra.Command{
		Use:     "cities",
		Short:   "Get a list of cities where PerfOps nodes are present",
		Long:    `Get a list of cities where PerfOps nodes are present`,
		Example: `perfops cities`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newPerfOpsClient()
			if err != nil {
				return err
			}
			return chkRunError(runCitiesCmd(c))
		},
	}
)

func initCitiesCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(citiesCmd)
}

func runCitiesCmd(c *perfops.Client) error {
	var res *[]perfops.City

	ctx := context.Background()
	u := c.BasePath + "/analytics/dns/city"

	f := internal.NewFormatter(debug && !outputJSON)
	f.StartSpinner()

	req, _ := http.NewRequest("GET", u, nil)
	req = req.WithContext(ctx)

	err := c.DoRequest(req, &res);
	f.StopSpinner()

	if err != nil {
		return err
	}

	internal.PrintOutputJSON(res)

	return nil
}
