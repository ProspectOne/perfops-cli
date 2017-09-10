# Build scripts

The build scripts use the latest version tag (in the format of
`v\d+.\d+.\d+`, e.g., `v0.1.0`) to set the version in the binary and
packages.

## `build-all.sh`

`build-all.sh` builds the `perfops` binaries. The environment variable
`PERFOPS_BUILD_PLATFORMS` can be set to specify the platforms (default:
`"linux windows darwin"`) to build, while `PERFOPS_BUILD_ARCHS` specifies
the architectures to build for each platform (default: `"amd64"`).

## `build-pkgs.sh`

`build-pkgs.sh` build the DEB and RPM packages.

To push the packages to [packagecloud.io](https://packagecloud.io/) the
environment variable `PACKAGECLOUD_TOKEN` should be defined.
