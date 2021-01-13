# klog

**Time tracking with plain text files.**

[**Download**](https://www.github.com/jotaen/klog/releases) – Get the latest version of the command line tool

[Guide](docs/Guide.md) – Learn how to use klog

[Specification](docs/Specification.md) – Understand the file format in-depth

## How does it work?

Tracking time with klog is based on the idea of storing the data in
plain text files that are formatted in a human-readable way, as an example:

```
2018-03-24
	| First day at my new job 🎉 
	| Mostly #onboarding and getting to know everyone
	8:00 - 17:15
	-1h | Lunch break
```

Let’s see: someone apparently started a new job here and wants to use
klog for tracking work times. At the 24th of March 2018 they started to
work at 8:00 in the morning and went home at 17:15 in the afternoon.
Somewhere in between they took a one-hour lunch break. This results
in a net total time of 8 hours and 15 minutes for that particular day.

You can store this record in a file (e.g. `times.klg`) and then use
the klog command line tool to run all sorts of evaluations, for instance:

```bash
$ klog total times.klg
Total time: 8h15m
```

You can also conveniently manipulate it, like so:

```bash
$ klog track 15m --date=2018-03-24 times.klg
Date: 2018-03-24 (2 entries)
Added entry: 15m
New total: 8h30m
```

For MacOS there is also a menu bar widget bundled into the command line tool
that allows quick and convenient access to the most important functionalities.
In order to start it up, just run `klog widget` and take it from there.

Learn more about klog and how to use it by reading the [guide](docs/Guide.md).

## About

klog is free and open-source software distributed under the [MIT license](LICENSE.txt).
