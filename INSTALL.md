# Install klog

In order to not miss any updates you can either subscribe to the release
notifications on [Github](https://github.com/jotaen/klog) (at the top right:
“Watch”→“Custom”→“Releases”), or you occasionally check by running
`klog version`.

For an archive of all klog releases, [see here](https://github.com/jotaen/klog/releases).

## MacOS
1. Download the latest version and unzip
   - [**Download for Intel**](https://github.com/jotaen/klog/releases/latest/download/klog-mac-intel.zip)
   - [**Download for M1 (ARM)**](https://github.com/jotaen/klog/releases/latest/download/klog-mac-arm.zip)
2. Right-click on the binary and select “Open“
   (due to [Gatekeeper](https://support.apple.com/en-us/HT202491))
3. Copy to path, e.g. `mv klog /usr/local/bin/klog` (might require `sudo`)

## Linux
1. [**Download**](https://github.com/jotaen/klog/releases/latest/download/klog-linux.zip)
   the latest version and unzip
2. Copy to path, e.g. `mv klog /usr/local/bin/klog` (might require `sudo`)

## Windows
1. [**Download**](https://github.com/jotaen/klog/releases/latest/download/klog-windows.zip)
   the latest version and unzip
2. Copy to path, e.g. to `C:\Windows\System32` (might require admin privileges)

By the way, as an alternative you can also use the Linux binary on
the [Windows Subsystem for Linux](https://docs.microsoft.com/en-us/windows/wsl/install-win10).

# Build klog from sources

Instead of downloading the binaries, you can also build klog yourself.

As prerequisite, you need to have the [Go compiler](https://golang.org/doc/install).
Please check the [`go.mod`](go.mod) file to see what Go version klog requires.

Fetch the sources:

```
git clone https://github.com/jotaen/klog.git
cd klog
```

In order to build the project, run:

```
go build src/app/cli/main/klog.go
```

This automatically resolves the dependencies and compiles the source code into an
executable for your platform.
