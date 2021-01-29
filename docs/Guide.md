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
A `.klg` file can contain any number of records
that each consists of a date, some time entries,
and optional summary texts.

```klog
2019-07-22
    13:00 - 14:30 Workout
    2h30m Reading books

2019-07-25
Chores and housework
    1h
    11:23 - 12:46
```

Records are separated by one blank line between them.
The first line of a record must be a date (formatted either
`YYYY-MM-DD` or `YYYY/MM/DD`).

### Tracking time
Entries are the actual time values that you track.
They appear one per line and are indented by one level.
(Indentation is either 1 tab or 2–4 spaces.)

```klog
2019-07-22
Both entries below are worth 1 hour each,
resulting in a total of 2 hours for this day.
    8:00 - 9:00
    1h
```

You can record a duration (e.g. `1h`, `-2h33m`, `48m`)
which represent an “amount” of time spent on something.
This can also be a negative value, in which case it will be
deducted from the overall total.

An entry can also be time range (e.g. `12:32 - 17:20` or `8:45am - 1:30pm`)
representing a time span between two points in time throughout a day.

### Summary
A summary is free text that can (optionally) appear:

- Underneath the date,
  in which case it is supposed to refer to the entire record.
  It can have multiple lines of text and is no indented.
- Behind entries on the same line,
  in which case it is only supposed to refer to that very entry.

### Tagging / categorising
Summaries can contain `#hashtags` that allow for more fine-granular
filtering of the data.

```klog
2019-07-22
If a tag appears in the #overall summary,
it applies to all time entries.
    4h Otherwise it only applies to an individual #entry
    5h
```

Here, the grand total for the tag `#overall` would be `9h`,
because the tag appears in the overall summary and therefore
all entries match it.
The grand total for the tag `#entry` would be `4h`, because
it only refers to one individual entry.

### Day shifting
Sometimes you start an activity in the evening and then stop
it after midnight.

```klog
2019-07-26
Last day of the week!
    <23:30 - 8:00 Night shift
    22:30 - 0:30> Watching movies
```

In this case you can “shift over” a time to the previous or to the next day
by adding the `<` prefix or the `>` suffix respectively.

When filtering records, keep in mind that these entries are still
associated with the date they are recorded under, so the grand total
for the above date `2019-07-26` is `10h30m`.

### Open-ended time ranges
In case you just begin an activity (without knowing when it will end)
you can already put it into your file as an open-ended time range.

```klog
2019-07-26
Just started my work day
    8:30 - ?
```

Open-ended time ranges are denoted by replacing the end time
with a question mark.
There can only be one open-ended range per record and it doesn’t
count towards the grand total.

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
The “should-total” always applies to the entire record with all its entries.

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
