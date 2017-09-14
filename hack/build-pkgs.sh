#!/bin/bash
# Copyright 2017 The PerfOps-CLI Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# This script will build perfops-cli DEB and RPM packages and push them to
# packagecloud.io.

set -e

VERSION=$(git tag -l "v*" --sort=-version:refname | head -1)
COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null)
PKG_VERSION="${VERSION:1}"

sudo apt-get update -y
sudo apt-get install ruby-dev build-essential rpm -y
sudo gem install --no-ri --no-rdoc fpm
sudo gem install --no-ri --no-rdoc package_cloud
sudo gem install --no-ri --no-rdoc rest-client

mkdir -p release/pkgs

PERFOPS_FPM_TARGETS="deb rpm"
for TARGET in ${PERFOPS_FPM_TARGETS[@]}; do
	echo "Building package $TARGET"
	fpm -t $TARGET --rpm-os linux -v $PKG_VERSION -s dir -p release/pkgs/ \
		-n perfops --license ALv2 --vendor ProspectOne -m "PerfOps Support <dak@prospectone.io>" \
		--description 'A simple command line tool to access the Prospect One PerfOps API.' \
		--url https://perfops.net/cli \
		release/perfops-linux-amd64=/usr/local/bin/perfops
done

RPM_DSTS="sles/11.4 sles/12.0 sles/12.1 sles/12.2
 opensuse/13.1 opensuse/13.2 opensuse/42.1 opensuse/42.2
 fedora/20 fedora/21 fedora/22 fedora/23 fedora/24 fedora/25 fedora/26
 el/5 el/6 el/7"

DEB_DSTS="debian/wheezy debian/jessie debian/stretch debian/buster
 ubuntu/trusty ubuntu/utopic ubuntu/vivid ubuntu/wily ubuntu/xenial ubuntu/yakkety ubuntu/zesty
 raspbian/wheezy raspbian/jessie raspbian/stretch raspbian/buster"

if [[ $PACKAGECLOUD_TOKEN ]]; then
	for PKG in release/pkgs/*; do
		if [[ $PKG == *.deb ]]; then
			for DST in ${DEB_DSTS[@]}; do
				echo "Uploading $PKG to $DST"
				package_cloud push p1/perfops/$DST $PKG --skip-errors
			done
		elif [[ $PKG == *.rpm ]]; then
			for DST in ${RPM_DSTS[@]}; do
				echo "Uploading $PKG to $DST"
				package_cloud push p1/perfops/$DST $PKG --skip-errors
			done
		fi
	done

	ruby hack/prune-pkgs.rb
fi
