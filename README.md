# klog

klog is a file format and a command line tool for time tracking.

The idea behind klog is to store data in plain text files
in a simple and human-readable format.
The notation is similar to how you would write down the information
into a physical notebook using pen and paper.
By using the klog command line tool you can search and evaluate your
data.

 ✦  [**Download**](https://www.github.com/jotaen/klog/releases) – Get the latest version of the command line tool

 ✦  [Documentation](https://klog.jotaen.net) – Learn how to use klog

 ✦  [Specification](Specification.md) – Read the file format specification 

## klog in a nutshell

![Screen recording of the demo described below](https://klog.jotaen.net/demo.gif)

Let’s say you started a new job and want to use klog for tracking work times:

```klog
2018-03-24
First day at my new job
    8:30-17:00
    -45m Lunch break
```

At the 24th of March 2018 you came to the office at 8:30
in the morning and went home at 17:00 in the afternoon.
Somewhere in between there was a 45-minute lunch break.

As you see, klog supports different notations to record times.
You can capture short summaries about your activities along with the
time data (only if you want, of course), which can help you later
on to make sense of what you did back in the day.
And by adding tags you are able to run more fine-granular evaluations.

When stored in a file (e.g. `worktimes.klg`) you can use the klog command
line tool to interact with it. For instance, you could evaluate the
resulting total time like this:

```
$ klog total --today worktimes.klg
Total: 7h45m
(In 1 record)
```

For MacOS users there is an experimental menu bar widget bundled into the command line tool
that allows quick and convenient access to the most important functionalities.
Just run `klog widget` to start it up and take it from there.

Learn more about klog and how to use it by reading the [documentation](https://klog.jotaen.net).

## Contribute

This repository contains the specification of the klog file format
as well as the sources of the command line tool.
While the CLI tool is based on the file format, it still evolves independently.

- **File format**: current state is RFC (request for comments) for version 1.
  Please see the [Specification](Specification.md) for further details.
- **Command line tool**: if you have ideas, run into a problem,
  or just want to bounce off some feedback, feel invited to open an
  [issue on Github](https://github.com/jotaen/klog/issues) so that we can discuss it.

## About

klog is free and open-source software distributed under the [MIT license](LICENSE.txt).
