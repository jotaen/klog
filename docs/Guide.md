# Guide

- [File format](#file-format)
    - [Tracking time](#tracking-time)
    - [Summary](#summary)
    - [Tagging / categorising](#tagging--categorising)
    - [Day shifting](#day-shifting)
    - [Open-ended time ranges](#open-ended-time-ranges)
    - [Should-total](#should-total)
    - [FAQ](#faq)
- [Command line tool](#command-line-tool)
- [Menu bar widget (MacOS)](#menu-bar-widget-macos)

## File format
A `.klg` file could look like this:

```klog
2019-07-22
    13:00 - 14:30 Workout
    2h30m Reading books

2019-07-25
Chores and housework
    1h
    11:23 - 12:46
```

It can contain any number of sections (called *records* in klog)
that each consists of a date, times and
(optionally) a summary of how the time was spent.

Records are separated by one blank line between them.
The first line of a record must be a date (formatted either
`YYYY-MM-DD` or `YYYY/MM/DD`).

### Tracking time
Entries for tracking time appear one per line and are indented by one level.

```klog
2019-07-22
Both entries below are worth 1 hour each,
resulting in a total of 2 hours for this day.
    8:00 - 9:00
    1h
```

You can either record a duration (e.g. `1h`, `2h33m`, `48m`)
or a time range (e.g. `12:32 - 17:20` or `8:45am - 1:30pm`).

### Summary
A summary is free text that can appear…

- underneath the date,
  in which case it is supposed to refer to the entire record
- behind entries on the same line,
  in which case it is only supposed to refer to that very entry

Summaries are optional.

### Tagging / categorising
Summaries can contain `#hashtags` that allow for more fine-granular
filtering of the data.

```klog
2019-07-22
If a tag appears in the #overall summary,
it applies to all time entries.
    4h Otherwise it only applies to the respective #entry
    5h
```

Here, the grand total for the tag `#overall` would be `9h`;
for the tag `#entry` it would be `4h`.
(And if you filter by both tags, it’s `9h`.)

### Day shifting
Sometimes you start an activity in the evening and then stop
it after midnight.

```klog
2019-07-26
Last day of the week!
    <23:30 - 8:00 Night shift
    22:30 - 0:30> Watching movies
```

You can “shift over” a time to the previous or subsequent day
by adding the `<` prefix or the `>` suffix respectively.

When filtering records, keep in mind that these entries are still
associated with the date they are recorded under, so the grand total
for the above date is `10h30m`.

### Open-ended time ranges
In case you started an activity that is still ongoing you
can denote an open-ended time range by replacing the second
value with a question mark:

```klog
2019-07-26
Just started my work day
    8:30 - ?
```

There can only be one open-ended range per record and
when evaluating the record they are skipped.

### Should-total
There are use-cases where you have a certain overall time goal
that you want to achieve.
For example, let’s say you are supposed to work 7½ hours per day:

```klog
2019-07-26 (7h30m!)
    8:00 - 16:00 Work
    -45m lunch break
```

This “should-total” is a mere meta-property. It can be used during
evaluation in order to diff the actual times against it.

### FAQ

- **Is it possible to use to-the-second precision,
  like `1h10m30s` or `8:23:49`?**
  No, this is not supported.
  The reason is that it would effectively prohibit mixing values
  with and without seconds, which leads to a lot of hassle.
  Keep in mind, klog is for time-tracking activities, it’s not a stopwatch.
- **How can I capture timezone information?**
  You cannot.
  In case you are affected by a timezone change or
  a switch to daylight saving time
  you need to account for that yourself.
  (Realistically, this doesn’t happen all too often anyway.)
- **Can there be multiple records for the same date in one file?**
  Yes, as many as you want.
- **Will it lead to problems if I track more than 24 hours per day,
  or if the resulting total is a negative value?**
  No, klog doesn’t care about that.
  (There are actually legitimate use-cases for this.)

## Command line tool
The command line tool currently allows you to pretty print and
evaluate files in your terminal. In order to learn about its usage
you should run `klog --help`.

For example, if you want to evaluate all records in `sport.klg`
from 2018 that are tagged with `#workout`, you would do:

```
$ klog eval --after=2018-01-01 --before=2018-12-31 --tag=workout sport.klg
```

Or if you want an ongoing counter of the current day to be displayed:

```
$ klog eval --today --diff --live worktimes.klg
```

Pro-tip: most shells have native support for glob patterns, so in case
you want to organise your records in multiple files (e.g. one file per month)
you can evaluate them all at once by doing `klog eval *.klg`.

## Menu bar widget (MacOS)
On MacOS you can launch a menu bar widget by running

```
$ klog widget --file=worktimes.klg
```

It displays an ongoing counter of the current day and a summary
of the entire file in your menu bar / system tray.
