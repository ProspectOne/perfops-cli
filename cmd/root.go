package cmd

import (
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

const versionTmpl = `perfops:
	version:     %s
	build date:  %s
	git hash:    %s
	go version:  %s
	go compiler: %s
	platform:    %s/%s
`

var (
	// rootCmd is the root command of the application.
	rootCmd = &cobra.Command{
		Use:   "perfops",
		Short: "perfops is a tool to interact with the PerfOps API",
		Long:  `perfops is a tool to interact with the PerfOps API.`,
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				cmd.Printf(versionTmpl,
					version, buildDate, commitHash,
					runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
				return
			}
			cmd.Usage()
		},
	}

	// ApiKey defines the PerfOps API to use for API calls.
	apiKey      string
	showVersion bool

	from string

	// Version information set at build time
	version    = "devel"
	buildDate  string
	commitHash string
)

// Execute executes the root command.
func Execute() error {
	initRootCmd()
	initlatencyCmd()
	initMTRCmd()
	initPingCmd()
	initTracerouteCmd()
	initDNSResolveCmd()
	return rootCmd.Execute()
}

func initRootCmd() {
	envAPIKey := os.Getenv("PERFOPS_API_KEY")

	rootCmd.PersistentFlags().StringVarP(&apiKey, "key", "K", envAPIKey, "The PerfOps API key (default is $PERFOPS_API_KEY)")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Prints the version information of perfops")

	rootCmd.PersistentFlags().StringVarP(&from, "from", "F", "", "A continent, region (e.g eastern europe), country, US state or city")
}
