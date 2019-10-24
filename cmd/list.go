package cmd

import (
	"errors"
	"fmt"
	"github.com/ProspectOne/perfops-cli/perfops"
	"github.com/spf13/cobra"
)

var (
	listTypesMap = map[string]func(client *perfops.Client) error{
		"countries": runCountriesCmd,
		"cities":    runCitiesCmd,
	}

	listCmd = &cobra.Command{
		Use:     "list [type]",
		Short:   "Get a locations where PerfOps nodes are present",
		Long:    `Get a locations where PerfOps nodes are present, e.g., 'list countries' or 'list cities'`,
		Example: `perfops list countries`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("specify the type of data You want to get")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newPerfOpsClient()
			if err != nil {
				return err
			}
			return chkRunError(runListCmd(c, args[0]))
		},
	}
)

func initListCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(listCmd)
}

func runListCmd(c *perfops.Client, dataType string) error {
	if f, ok := listTypesMap[dataType]; ok {
		err := f(c)
		if err != nil {
			return errors.New(fmt.Sprintf("error happened %v", err))
		}

		return nil
	}

	return errors.New(fmt.Sprintf("no data with type '%s'", dataType))
}
