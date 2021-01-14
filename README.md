# klog

**Time tracking with plain text files.**

[**Download**](https://www.github.com/jotaen/klog/releases) – Get the latest version of the command line tool

[Guide](docs/Guide.md) – Learn how to use klog

[Specification](docs/Specification.md) – Understand the file format in-depth

## What is klog?

Time tracking with klog is based on the idea of storing data in
plain text files in a human-readable format. Here is an example:

```klg
2018-03-24
First day at my new job
Setup computer and started onboarding
	8:00 - 17:15
	-1h Lunch #break at Sushi place
```

In this example someone started a new job and uses klog for tracking
work times. At the 24th of March 2018 they came to the office at 8:00
in the morning and went home at 17:15 in the afternoon. Somewhere in
between there was a one-hour lunch break.

You can store this data in a file (e.g. `times.klg`) and use the
klog command line tool to run all sorts of evaluations. For instance,
you could check what the net total time is:

```bash
$ klog total
Total time: 8h15m
```

You can manipulate the data by editing the file manually, or you can use
the command line tool. Let’s say you wanted to add another 15 minutes:

```bash
$ klog track 15m --date=2018-03-24 times.klg
Date: 2018-03-24 (2 entries)
Added entry: 15m
New total: 8h30m
```

For MacOS users there is a menu bar widget bundled into the command line tool
that allows quick and convenient access to the most important functionalities.
Just run `klog widget` to start it up and take it from there.

Learn more about klog and how to use it by reading the [guide](docs/Guide.md).

## About

klog is free and open-source software distributed under the [MIT license](LICENSE.txt).
