# PerfOps cli - Global network testing and benchmarking
[![Build Status](https://semaphoreci.com/api/v1/projects/77896bab-6c47-4549-8018-05f07b60d941/1495977/badge.svg)](https://semaphoreci.com/prospectone/perfops-cli)

A simple command line tool to interact with hundreds of servers around the world. Run benchmarks and debug your infrastructure without leaving your console. [More information](https://perfops.net/cli)

## [Install instructions](https://github.com/ProspectOne/perfops-cli/blob/master/INSTALL.md)

## Usage

Help screen
```
$ perfops -h
perfops is a tool to interact with the PerfOps API.

Usage:
  perfops [flags]
  perfops [command]

Available Commands:
  curl        Run a curl test on a domain name or IP address
  dnsperf     Find the time it takes to resolve a DNS record on a target
  help        Help about any command
  latency     Run a ICMP latency test on a domain name or IP address
  mtr         Run a MTR test on a domain name or IP address
  ping        Run a ping test on a domain name or IP address
  resolve     Resolve a DNS record on a domain name
  traceroute  Run a traceroute test on a domain name or IP address

Flags:
      --debug             Enables debug output
  -F, --from string       A continent, region (e.g eastern europe), country, US state or city
  -h, --help              help for perfops
  -J, --json              Print the result of a command in JSON format
  -K, --key string        The PerfOps API key (default is $PERFOPS_API_KEY)
  -N, --nodeid intSlice   A comma separated list of node IDs to run a test from
  -v, --version           Prints the version information of perfops

Use "perfops [command] --help" for more information about a command.
```

## Examples

Ping google.com from a random server in Eastern Europe
```
perfops ping --from "eastern europe" google.com
Node111, Moscow, Russian Federation
PING google.com (173.194.222.113) 56(84) bytes of data.
64 bytes from 173.194.222.113: icmp_seq=1 ttl=50 time=11.6 ms
64 bytes from 173.194.222.113: icmp_seq=2 ttl=50 time=11.4 ms
64 bytes from 173.194.222.113: icmp_seq=3 ttl=50 time=11.4 ms

--- google.com ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 602ms
rtt min/avg/max/mdev = 11.433/11.513/11.650/0.157 ms
```

Traceroute to google.com from a server located in New York
```
 perfops traceroute --from "New York" google.com
Node15, New York City, United States
traceroute to google.com (172.217.10.46), 20 hops max, 60 byte packets
 1  vl223-ar-02.nyc-ny.atlantic.net (45.58.33.35)  0.432 ms  0.420 ms
 2  vl30-ar-01.nyc-ny.as6364.net (45.58.33.1)  0.452 ms  0.411 ms
 3  te0-0-1-1.rcr11.ewr04.atlas.cogentco.com (38.104.44.141)  1.153 ms  1.145 ms
 4  154.24.52.17 (154.24.52.17)  1.142 ms te0-3-0-4.rcr21.ewr02.atlas.cogentco.com (154.24.9.9)  1.042 ms
 5  be2390.rcr23.jfk01.atlas.cogentco.com (154.54.80.189)  1.502 ms be2600.rcr23.jfk01.atlas.cogentco.com (154.54.40.29)  1.438 ms
 6  be2896.ccr41.jfk02.atlas.cogentco.com (154.54.84.201)  2.397 ms  2.193 ms
 7  be3294.ccr31.jfk05.atlas.cogentco.com (154.54.47.218)  2.319 ms  2.422 ms
 8  tata.jfk05.atlas.cogentco.com (154.54.12.18)  1.997 ms  1.955 ms
 9  if-ae-12-2.tcore1.N75-New-York.as6453.net (66.110.96.5)  2.256 ms  2.314 ms
10  72.14.195.232 (72.14.195.232)  2.125 ms  2.112 ms
11  * *
12  216.239.62.169 (216.239.62.169)  1.621 ms 216.239.62.171 (216.239.62.171)  1.501 ms
13  lga34s13-in-f14.1e100.net (172.217.10.46)  1.826 ms  1.857 ms
```

Check ICMP latency from 9 servers located in Europe

```
perfops latency --from europe --limit 9 google.com
Node92, Arezzo, Italy
7.705
Node242, Meppel, Netherlands
2.753
Node215, Nottingham, United Kingdom
9.861
Node85, Kiev, Ukraine
15.332
Node196, Riga, Latvia
47.940
Node244, Zürich, Switzerland
12.591
Node194, Nuremberg, Germany
3.697
Node259, Luxembourg, Luxembourg
7.928
Node76, Vilnius, Lithuania
24.506
```

## Setup

If you are interested in building `perfops` from source, you can install
it via `go get`:

```sh
go get -u github.com/ProspectOne/perfops-cli -o perfops
```

## Feedback

Feedback is greatly appreciated.

## Contributing

Contributions are greatly appreciated. The maintainers actively manage the
issues list, and try to highlight issues suitable for newcomers. The project
follows the typical GitHub pull request model. See
[CONTRIBUTING.md](CONTRIBUTING.md) for more details. Before starting any
work, please either comment on an existing issue, or file a new one.
