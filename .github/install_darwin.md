# Install klog

In order to install the downloaded klog binary on your system, follow these steps:

1. Make [MacOS “Gatekeeper”](https://support.apple.com/en-us/HT202491) trust the executable:
   - Either right-click on the binary in the Finder, and select “Open“
   - Or remove the “quarantine” flag from the binary via the CLI:
     `xattr -d com.apple.quarantine klog`
2. Copy it to a location that’s covered by your `$PATH` environment variabke, e.g.
   `mv klog /usr/local/bin/klog` (might require `sudo`)

In order to not miss any updates you can either subscribe to the release
notifications on [Github](https://github.com/jotaen/klog) (at the top right:
“Watch”→“Custom”→“Releases”), or you occasionally check by running
`klog version`.

For other install options, see [the documentation website](https://klog.jotaen.net/#get-klog).
