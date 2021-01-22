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
	9:00 - 17:30
	-45m Lunch break
```

In this example someone started a new job and uses klog for tracking
work times. At the 24th of March 2018 they came to the office at 9:00
in the morning and went home at 17:30 in the afternoon. Somewhere in
between there was a 45-minute lunch break.

As you see, klog supports different notations to record times. It also
allows you to capture short summaries about your activities (only if you
want, of course), which can help you later on to make sense of what
you did back in the day. They also allow you to categorise your data.

When stored in a file (e.g. `times.klg`) you can use the klog command
line tool to interact with it. For instance, you could evaluate the
resulting total time like this:

```bash
$ klog evaluate --today times.klg
Total: 7h45m
(In 1 records)
```

For MacOS users there is a menu bar widget bundled into the command line tool
that allows quick and convenient access to the most important functionalities.
Just run `klog widget` to start it up and take it from there.

Learn more about klog and how to use it by reading the [guide](docs/Guide.md).

## Why use klog?

klog is for you if:

- you want to track times of your activities with to-the-minute precision
- you like to be able to preserve context by putting notes alongside the
  numbers
- you value the freedom and simplicity of plain text files

## About

klog is free and open-source software distributed under the [MIT license](LICENSE.txt).
