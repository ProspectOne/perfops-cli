package cmd

import (
	"github.com/spf13/cobra"
	"testing"
)

func TestInitListCmd(t *testing.T) {
	parent := &cobra.Command{}
	listCmd.ResetFlags()
	initListCmd(parent)
	flags := listCmd.HasAvailableLocalFlags()
	if flags == true {
		t.Fatalf("expected that command doesn't have flags")
	}
}

func TestRunListCmd(t *testing.T) {
	testCases := map[string]struct {
		dataType    string
		url         string
		expectedErr interface{}
	}{
		"Countries":  {"countries", "/analytics/dns/countries", nil},
		"Cities":     {"cities", "/analytics/dns/city", nil},
		"Wrong type": {"country", "/analytics/dns/country", "no data with type 'country'"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tr := &recordingTransport{}
			c, err := newTestPerfopsClient(tr)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}

			err = runListCmd(c, tc.dataType)

			if tc.expectedErr != nil {
				if tc.expectedErr != err.Error() {
					t.Fatalf("expected %v error; got %v", tc.expectedErr, err)
				}

				// No reason to continue if we check for error case
				return
			}

			if got, exp := tr.req.URL.Path, tc.url; got != exp {
				t.Fatalf("expected %v; got %v", exp, got)
			}
		})
	}
}
