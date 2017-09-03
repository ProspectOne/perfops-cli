package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/ProspectOne/perfops-cli/perfops"
)

type (
	runFunc       func(ctx context.Context, c *perfops.Client, req *perfops.RunRequest) (perfops.TestID, error)
	runOutputFunc func(ctx context.Context, c *perfops.Client, pingID perfops.TestID) (*perfops.RunOutput, error)
)

// RunTest runs an MTR or ping testm retrives its output and presents it to the user.
func RunTest(ctx context.Context, client *perfops.Client, target, location string, limit int, runTest runFunc, getRunOutput runOutputFunc) error {
	runReq := &perfops.RunRequest{
		Target:   target,
		Location: location,
		Limit:    limit,
	}

	spinner := NewSpinner()
	fmt.Println("")
	spinner.Start()

	testID, err := runTest(ctx, client, runReq)
	if err != nil {
		spinner.Stop()
		return err
	}

	var output *perfops.RunOutput
	for {
		select {
		case <-time.After(250 * time.Millisecond):
		}

		if output, err = getRunOutput(ctx, client, testID); err != nil {
			spinner.Stop()
			return err
		}
		if output.IsFinished() {
			break
		}
	}

	spinner.Stop()

	for _, item := range output.Items {
		n := item.Result.Node
		fmt.Printf("Node%d, %s, %s\n%s\n", n.ID, n.City, n.Country.Name, item.Result.Output)
	}
	return nil
}
