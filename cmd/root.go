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
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

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
		Use:          "perfops",
		Short:        "perfops is a tool to interact with the PerfOps API",
		Long:         `perfops is a tool to interact with the PerfOps API.`,
		Example:      `perfops traceroute --from "New York" google.com`,
		SilenceUsage: true,
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
	debug       bool

	from    string
	nodeIDs []int

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
	initCurlCmd(rootCmd)
	return rootCmd.Execute()
}

func initRootCmd() {
	cobra.OnInitialize(initConfig)
	apiKey = os.Getenv("PERFOPS_API_KEY")

	rootCmd.PersistentFlags().StringVarP(&apiKey, "key", "K", "", "The PerfOps API key (default is $PERFOPS_API_KEY)")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Prints the version information of perfops")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "Enables debug output")

	rootCmd.PersistentFlags().StringVarP(&from, "from", "F", "", "A continent, region (e.g eastern europe), country, US state or city")
	rootCmd.PersistentFlags().IntSliceVarP(&nodeIDs, "nodeid", "N", []int{}, "A comma separated list of node IDs to run a test from")
}

// newPerfOpsClient returns a perfops.Client object initialized with the
// API key.
func newPerfOpsClient() (*perfops.Client, error) {
	return perfops.NewClient(perfops.WithAPIKey(apiKey))
}

func initConfig() {
	if apiKey == "" {
		apiKey = os.Getenv("PERFOPS_API_KEY")
	}
}

// requireTarget returns an error if no target is specified.
func requireTarget() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("no target specified")
		}
		return nil
	}
}

func chkRunError(err error) error {
	type argNamer interface {
		ArgName() string
	}
	if perfops.IsUnauthorized(err) {
		// A bit of a hack...
		err = errors.New("The API token was declined. Please correct it or do not send a token to use the free plan.")
	} else if namer, ok := err.(argNamer); ok {
		rootCmd.SetHelpFunc(invalidArgHelp(namer.ArgName()))
		err = flag.ErrHelp
	}
	return err
}

func invalidArgHelp(name string) func(cmd *cobra.Command, args []string) {
	if name == "limit" {
		return func(cmd *cobra.Command, args []string) {
			cmd.Println("For free users the maximum allowed number nodes for a single test is 20. Please change your limit.")
		}
	}
	return func(cmd *cobra.Command, args []string) {
		cmd.Println("Missing or invalid arguments:")
		cmd.Flags().VisitAll(func(f *flag.Flag) {
			if f.Name == name {
				cmd.Println(fmtShortUsage(f))
			} else if _, ok := f.Annotations[cobra.BashCompOneRequiredFlag]; ok {
				cmd.Println(fmtShortUsage(f))
			}
		})
	}
}

func fmtShortUsage(f *flag.Flag) string {
	line := ""
	if f.Shorthand != "" && f.ShorthandDeprecated == "" {
		line = fmt.Sprintf("  -%s, --%s", f.Shorthand, f.Name)
	} else {
		line = fmt.Sprintf("  --%s", f.Name)
	}
	varname, _ := flag.UnquoteUsage(f)
	if varname != "" {
		line += " " + varname
	}
	return line
}
