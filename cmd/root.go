// Copyright 2017 The PerfOps-CLI Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/ProspectOne/perfops-cli/perfops"
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
	initLatencyCmd(rootCmd)
	initMTRCmd(rootCmd)
	initPingCmd(rootCmd)
	initTracerouteCmd(rootCmd)
	initDNSResolveCmd(rootCmd)
	return rootCmd.Execute()
}

func initRootCmd() {
	envAPIKey := os.Getenv("PERFOPS_API_KEY")

	rootCmd.PersistentFlags().StringVarP(&apiKey, "key", "K", envAPIKey, "The PerfOps API key (default is $PERFOPS_API_KEY)")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Prints the version information of perfops")

	rootCmd.PersistentFlags().StringVarP(&from, "from", "F", "", "A continent, region (e.g eastern europe), country, US state or city")
}

// newPerfOpsClient returns a perfops.Client object initialized with the
// API key.
func newPerfOpsClient() (*perfops.Client, error) {
	return perfops.NewClient(perfops.WithAPIKey(apiKey))
}
