## Quick install

### Ubuntu/Debian

```
curl -s https://packagecloud.io/install/repositories/p1/perfops/script.deb.sh | sudo bash
apt-get install perfops
```

### RHEL/CentOS/Fedora

```
curl -s https://packagecloud.io/install/repositories/p1/perfops/script.rpm.sh | sudo bash
yum install perfops
```

## Manual

### deb

```
curl -L https://packagecloud.io/p1/perfops/gpgkey | sudo apt-key add -
```

Make sure to replace ubuntu and trusty in the config below with your Linux distribution and version:
```
deb https://packagecloud.io/p1/perfops/ubuntu/ trusty main
deb-src https://packagecloud.io/p1/perfops/ubuntu/ trusty main
apt-get update
```

Install
```
apt-get install perfops
```

### rpm

Create a file named /etc/yum.repos.d/perfops.repo that contains the repository configuration below.

Make sure to replace el and 6 in the config below with your Linux distribution and version:
```
[p1_perfops]
name=p1_perfops
baseurl=https://packagecloud.io/p1/perfops/el/6/$basearch
repo_gpgcheck=1
gpgcheck=0
enabled=1
gpgkey=https://packagecloud.io/p1/perfops/gpgkey
```

Install
```
yum update
yum install perfops
```


### zypper

Create a file named /etc/zypp/repos.d/perfops.repo that contains the repository configuration below.

```
[p1_perfops]
name=p1_perfops
baseurl=https://packagecloud.io/p1/perfops/opensuse/13.2/$basearch
enabled=1
repo_gpgcheck=1
pkg_gpgcheck=0
gpgkey=https://packagecloud.io/p1/perfops/gpgkey
autorefresh=1
type=rpm-md
```

Update your local zypper cache by running
```
zypper --gpg-auto-import-keys refresh perfops
```

Install
```
zypper install perfops
```
