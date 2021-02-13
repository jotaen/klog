# klog

klog is a plain-text file format and a command line tool for time tracking.

 ✦  [**Documentation**](https://klog.jotaen.net) – **Learn how to use klog**

 ✦  [Changelog](https://github.com/jotaen/klog/blob/main/CHANGELOG.md) – See the latest changes

 ✦  [Specification](Specification.md) – Understand the file format in-depth

## Get klog

### MacOS
1. [**Download**](https://www.github.com/jotaen/klog/releases) and unzip
2. Right-click on the binary and select “Open“ (due to [Gatekeeper](https://support.apple.com/en-us/HT202491))
3. Copy to path, e.g. `mv klog /usr/local/bin/klog` (might require `sudo`)

### Linux
1. [**Download**](https://www.github.com/jotaen/klog/releases) and unzip
2. Copy to path, e.g. `mv klog /usr/local/bin/klog` (might require `sudo`)

### Windows
There is no native binary for Windows, but you can use the Linux binary on the [Windows Subsystem for Linux](https://docs.microsoft.com/en-us/windows/wsl/install-win10).

## Contribute

This repository contains the specification of the klog file format
as well as the sources of the command line tool.

- **Command line tool**: if you have ideas, run into a problem,
  or just want to bounce off some feedback, feel invited to open an
  [issue on Github](https://github.com/jotaen/klog/issues) so that we can discuss it.
- **File format**: current state is RFC (request for comments) for version 1.
  Please see the [Specification](Specification.md) for further details.

The version numbers of the file format and the CLI tool are independent of each other. 

## About

klog is free and open-source software distributed under the [MIT license](LICENSE.txt).
