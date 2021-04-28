![klog logp](https://klog.jotaen.net/logo/klog-black-small.svg)

# klog

klog is a plain-text file format and a command line tool for time tracking.

 ✦  [**Documentation**](https://klog.jotaen.net) – **Learn how to use klog**

 ✦  [Changelog](https://github.com/jotaen/klog/blob/main/CHANGELOG.md) – See the latest changes

 ✦  [Specification](Specification.md) – Understand the file format in-depth

## Get klog

In order to not miss any updates you can either subscribe to the release notifications on Github
(at the top right: “Watch”→“Custom”→“Releases”), or you occasionally run `klog version`.

### MacOS
1. [**Download**](https://www.github.com/jotaen/klog/releases) and unzip
2. Right-click on the binary and select “Open“ (due to [Gatekeeper](https://support.apple.com/en-us/HT202491))
3. Copy to path, e.g. `mv klog /usr/local/bin/klog` (might require `sudo`)

### Linux
1. [**Download**](https://www.github.com/jotaen/klog/releases) and unzip
2. Copy to path, e.g. `mv klog /usr/local/bin/klog` (might require `sudo`)

### Windows
1. [**Download**](https://www.github.com/jotaen/klog/releases) and unzip
2. Copy to path, e.g. to `C:\Windows\System32` (might require admin privileges)

By the way, as an alternative you can also use the Linux binary on
the [Windows Subsystem for Linux](https://docs.microsoft.com/en-us/windows/wsl/install-win10).

## Contribute

If you have questions, feature ideas, or just want to bounce off some feedback
feel invited to [start a discussion](https://github.com/jotaen/klog/discussions).
In case you run into a bug please [file an issue](https://github.com/jotaen/klog/issues).
(When in doubt just go for an issue.)

This repository contains the sources of the command line tool as well as
the [specification](Specification.md) of the klog file format. Note that the
version numbers of both are independent of each other.

## Build klog from sources

As prerequisite, you need to have the [Go compiler](https://golang.org/doc/install).
Please check the [`src/go.mod`](src/go.mod) file to see what version klog requires. 
In order to build the project, navigate to the [`src/`](src) folder and run:

```
go build app/cli/main/klog.go
```

This automatically resolves the dependencies and compiles the source code into an
executable for your platform.

## About

klog is free and open-source software distributed under the [MIT license](LICENSE.txt).
