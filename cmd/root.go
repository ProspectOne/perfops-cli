package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// RootCmd is the root command of the application.
	RootCmd = &cobra.Command{
		Use:   "perfops-cli",
		Short: "perfops-cli is a tool to interact with the PerfOps API",
		Long:  `perfops-cli is a tool to interact with the PerfOps API.`,
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Printf(`perfops-cli:
 version:     %s
 build date:  %s
 git hash:    %s
 go version:  %s
 go compiler: %s
 platform:    %s/%s
`, version, buildDate, commitHash,
					runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
				return
			}
			cmd.Usage()
		},
	}

	// ApiKey defines the PerfOps API to use for API calls.
	apiKey      string
	showVersion bool

	// Version information set at build time
	version    = "devel"
	buildDate  string
	commitHash string
)

func init() {
	envAPIKey := os.Getenv("PERFOPS_API_KEY")

	RootCmd.PersistentFlags().StringVarP(&apiKey, "key", "K", envAPIKey, "The PerfOps API key (default is $PERFOPS_API_KEY)")
	RootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Prints the version information of perfops-cli")
}
