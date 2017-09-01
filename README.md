# perfops-cli

A simple command line tool to access with ProspectOne [PerfOps API](http://docs.perfops.net/).

## Setup

Download the latest binary from the [releases](https://github.com/ProspectOne/perfops-cli/releases) page.

If you are interested in building `perfops-cli` from source, you can install
it via `go get`:

```sh
go get -u github.com/ProspectOne/perfops-cli
```

## Usage

```
$ perfops-cli -h
perfops-cli is a tool to interact with the PerfOps API.

Usage:
  perfops-cli [flags]
  perfops-cli [command]

Available Commands:
  help        Help about any command
  ping        Run a ping test on target

Flags:
  -h, --help         help for perfops-cli
  -K, --key string   The PerfOps API key (default is $PERFOPS_API_KEY)
  -v, --version      Prints the version information of perfops-cli

Use "perfops-cli [command] --help" for more information about a command.
```

## Feedback

Feedback is greatly appreciated.

## Contributing

Contributions are greatly appreciated. The maintainers actively manage the
issues list, and try to highlight issues suitable for newcomers. The project
follows the typical GitHub pull request model. See
[CONTRIBUTING.md](CONTRIBUTING.md) for more details. Before starting any
work, please either comment on an existing issue, or file a new one.
