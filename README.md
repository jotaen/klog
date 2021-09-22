![klog logp](https://klog.jotaen.net/logo/klog-black-small.svg)

# klog

klog is a plain-text file format and a command line tool for time tracking.

 ✦  [**Documentation**](https://klog.jotaen.net) – **Learn how to use klog**

 ✦  [Download](INSTALL.md) – Get the latest version

 ✦  [Changelog](https://github.com/jotaen/klog/blob/main/CHANGELOG.md) – See what’s new

 ✦  [Specification](Specification.md) – Understand the file format in-depth

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
Please check the [`go.mod`](go.mod) file to see what version klog requires. 
In order to build the project, run:

```
go build src/app/cli/main/klog.go
```

This automatically resolves the dependencies and compiles the source code into an
executable for your platform.

In order to re-use klog in your own application, you have the following options:
- Execute `klog json` as subprocess and process the output.
  (Run `klog json --help` to explore all available options.)
- If your application is written in Go, you can add this package as code dependency.
  Please keep in mind that there are no separate releases of the source code itself,
  so the APIs and code structure might change anytime. Therefore, it’s recommended
  to specify the commit hash when adding the dependency.

## About

klog is free and open-source software distributed under the [MIT license](LICENSE.txt).
