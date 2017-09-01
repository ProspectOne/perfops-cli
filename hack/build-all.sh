#!/usr/bin/env bash
# Copyright 2017 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will build perfops-cli and calculate hash for each
# (PERFOPS_BUILD_PLATFORMS, PERFOPS_BUILD_ARCHS) pair.
# PERFOPS_BUILD_PLATFORMS="linux" PERFOPS_BUILD_ARCHS="amd64" ./hack/build-all.sh
# can be called to build only for linux-amd64

set -e

VERSION=$(git describe --tags --dirty)
COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null)
DATE=$(date --iso-8601)

VERSION_PKG="github.com/ProspectOne/perfops-cli/cmd"
API_PKG="github.com/ProspectOne/perfops-cli/api"
GO_BUILD_CMD="go build -a -installsuffix cgo"
GO_BUILD_LDFLAGS="-s -w -X $VERSION_PKG.commitHash=$COMMIT_HASH -X $VERSION_PKG.buildDate=$DATE -X $VERSION_PKG.version=$VERSION -X $API_PKG.libVersion=$VERSION"

if [ -z "$PERFOPS_BUILD_PLATFORMS" ]; then
    PERFOPS_BUILD_PLATFORMS="linux windows darwin"
fi

if [ -z "$PERFOPS_BUILD_ARCHS" ]; then
    PERFOPS_BUILD_ARCHS="amd64"
fi

mkdir -p release

for OS in ${PERFOPS_BUILD_PLATFORMS[@]}; do
  for ARCH in ${PERFOPS_BUILD_ARCHS[@]}; do
    echo "Building for $OS/$ARCH"
    GOARCH=$ARCH GOOS=$OS CGO_ENABLED=0 $GO_BUILD_CMD -ldflags "$GO_BUILD_LDFLAGS"\
     -o "release/perfops-cli-$OS-$ARCH" .
    sha256sum "release/perfops-cli-$OS-$ARCH" > "release/perfops-cli-$OS-$ARCH".sha256
  done
done
