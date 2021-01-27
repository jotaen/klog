# klog

klog is a file format and a command line tool for time tracking.

The idea behind klog is to store data in plain text files
in a simple and human-readable format.
The notation is similar to how you would write down the information
into a physical notebook using pen and paper.
By using the klog command line tool you can search and evaluate your
data.

 ✦  [**Download**](https://www.github.com/jotaen/klog/releases) – Get the latest version of the command line tool

 ✦  [Guide](docs/Guide.md) – Learn how to use klog

 ✦  [Specification](docs/Specification.md) – Understand the file format in-depth

## klog in a nutshell

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
You can capture short summaries about your activities along with your data
(only if you want, of course), which can help you later on to make sense
of what you did back in the day.
And by adding tags you are able to run fine-granular evaluations.

When stored in a file (e.g. `times.klg`) you can use the klog command
line tool to interact with it. For instance, you could evaluate the
resulting total time like this:

```shell
$ klog evaluate --today times.klg
Total: 7h45m
(In 1 record)
```

For MacOS users there is a menu bar widget bundled into the command line tool
that allows quick and convenient access to the most important functionalities.
Just run `klog widget` to start it up and take it from there.

Learn more about klog and how to use it by reading the [guide](docs/Guide.md).

## Current state: v1.0-rc

As of January 2021, the current state is the v1.0 release candidate (v1.0-rc).
Unless there are any obvious bugs or mistakes popping up,
the first stable version shall be timely released as is.

The main goals throughout the first version are:

- to validate the basic idea behind the file format
  (as this is the central pillar of klog)
- to learn more about different use-cases
- to make the command line tool capable of manipulating files
  (e.g. adding new entries)
- to find out how useful the MacOS widget is and whether it would
  be worth to provide such a graphical interface for other platforms

If you have an idea or thought you’d like to bounce off,
feel free to open an [issue on Github](https://github.com/jotaen/klog/issues) for further discussion.

## About

klog is free and open-source software distributed under the [MIT license](LICENSE.txt).
