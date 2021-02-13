# klog – File Format Specification

klog is a file format for tracking time.

It is free and open-source software distributed under the MIT-License.

> **Current state**: version 1 RFC (request for comments)
>
> This is a draft for the first version of the klog file format.
> While the basic structure will likely remain as it is,
> there still might be minor additions or corrections necessary.
> Time will tell when it’s good to go to be finalised as version 1.
> 
> In case you have comments or thoughts, please file an issue
> in the [klog repository](https://github.com/jotaen/klog)
> so that they can be discussed.

## Preface

The keywords “MUST”, “MUST NOT”, “SHOULD”, “SHOULD NOT”, “RECOMMENDED”, “NOT RECOMMENDED” and “MAY”
in this document are to be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119).

Whenever a word has special meaning in klog, it is formatted in *italics*.

Other technical terms are surrounded by “quotes”. These are defined at the end of this specification.

## I. Records

A *record* is a self-contained data structure that contains time-tracking information.

Each *record* MUST appear as one consecutive block in the file,
without any “blank lines” appearing within.

The first line of a *record* MUST start with a *date*.
On the same line there MAY follow a list of *properties*.
This list MUST be enclosed in “parentheses” and
separated by one “space” from the *date*
(there MAY also be multiple “spaces”).
If there are multiple *properties* they MUST be separated
by a `,` followed by a “space”.
The list of *properties* MUST NOT be empty.

A *summary* MAY appear on the subsequent lines.
Any amount of *entries* MAY appear afterwards.

### Date
A *date* is a day that is representable in the Gregorian calendar.

Each *record* MUST contain a *date*.

It MUST be formatted according to one of the following patterns:
- `YYYY-MM-DD` (RECOMMENDED),
- `YYYY/MM/DD`

(Where `Y` is a digit to denote the year, `M` the month, `D` the day.)

### Properties
*Property* is an abstract term that denotes additional information
or configuration of a *record*.
A *should-total* is an instance of a *property*.

### Should-Total
A *should-total* is a *property* to denote the targeted total time of a *record*.

A *should-total* MUST be a *duration* value followed by a `!`,
e.g. `8h!` or `5h30m!`.
(That implies that a negative value MAY be used.)

### Summary
A *summary* is user-provided text for holding arbitrary information.
There are two places where *summary* text MAY appear in *records*:

- Underneath the *date*:
  In this case the *summary* is considered to be associated with the entire *record*.
  The *summary* MAY span multiple lines.
  Each of its lines MUST NOT start with “whitespace”.
- Behind *entries*:
  In this case the *summary* is only considered to be referring to the corresponding *entry*.
  The *summary* text follows the *entry* on the same line,
  and it ends at the end of that line.
  It MUST be separated from the *entry* by one “space”
  (there MAY be multiple “spaces”).

### Tags
The purpose of *tags* is to help categorise records and entries.

Any amount of *tags* MAY appear anywhere within *summaries*.
A *tag* MUST be a sequence of “letters”, “digits” or the `_` character,
preceded by a single `#` character,
e.g. `#gym`, `#24hours` or `#home_office`.

### Entry
*Entry* is an abstract term for time-related data.
A *range* and a *duration* are instances of *entries*.

Each *entry* MUST appear on its own line and
MUST be indented in one of the following ways:
- by four “spaces” (RECOMMENDED)
- by two or three “spaces”
- by one “tab”

### Time
A *time* is a value that represents a point in time throughout a day
as it would be displayed by a wall clock (which divides a day into
24 hours and every hour into 60 minutes),
e.g. `9:00`, `23:18`, `6:30am`, `9:23pm` 

A *time* value MUST consist of both an hour part and a minute part.
Single-digit hour parts MAY be padded with a `0`.
The minute part MUST always contain two digits.

As default, *times* are to be interpreted as 24-hour clock values.
An `am` or `pm` suffix MAY be used to denote that the value is
to be interpreted as 12-hour clock value.

### Range
A *range* is an *entry* that represents the time span between two points in time.

It MUST consist of two values that denote the start and the end.
Start and end MUST be written in chronological order.

There MUST be a `-` between the two values.
There MAY appear “spaces” on either side of the `-`;
for this case it is RECOMMENDED to use exactly one “space” on both sides.

The start value MUST be a *time*.
It MAY be prefixed with a `<` to indicate that
this *time* is referring to the day before the *record’s* date,
e.g. `<23:00`.

The end value MUST be either a *time* or a placeholder for a *time*.
- In case the end value is a *time* it MAY be suffixed with a `>` to indicate
  that this *time* is referring to the day after the *record’s* date,
  e.g. `0:30>`.
- In case the end value is a placeholder the *range* is considered to be *open-ended*,
  which means that the end *time* is not determined yet.
  The placeholder MUST be denoted by a `?`, e.g. `9:00 - ?`;
  the `?` MAY be repeated, e.g. `9:00 - ???`.
  An *open-ended range* MUST NOT occur more than once per record.

### Duration
A *duration* is an *entry* that represents a period of time.
It contains an amount of hours and/or an amount of minutes.
(So it MUST either contain one of these two or both.)
The hour part MUST be written first.
Examples are: `1h`, `5m`, `4h12m`, `-8h30m`.

The hour part MUST be an “integer”
which MUST be immediately followed by the character `h`.
It MAY be `0h`.
It MAY be greater than `24h`,
e.g. `50h`.
If the hour part is missing, a value of `0h` is assumed.

The minute part MUST be an “integer”
which MUST be immediately followed by the character `m`.
It MAY be `0m`.
When the hour part is present,
the minute part MUST NOT be greater than `59m`,
e.g. `1h59m`;
otherwise it MAY be greater than `59m`,
e.g. `119m`
(it is RECOMMENDED to break this up, though).
If the minute part is missing, a value of `0m` is assumed.

The *duration* as a whole is a signed value:
That means it is either positive (i.e. adding to the total time)
or negative (i.e. deducting from the total time).
As default a *duration* is positive,
which MAY be indicated by a leading `+` character,
e.g. `+4h12m`.
If the *duration* is supposed to be negative, it MUST be preceded by a `-` character.

## II. Organizing records in files

A file MAY hold any amount of *records*.
Apart from that it MUST NOT contain anything
but what is allowed by this specification.

There MUST appear one “blank line” between subsequent *records*;
additional “blank lines” MAY appear.

The *records* don’t have to appear in any order.
There MAY exist multiple *records* for the same day.
These are treated as distinct.

### Technical remarks

The file extension MUST be `.klg`, e.g. `times.klg`.

The file encoding MUST be UTF-8.

Newlines MUST be encoded with the LF linefeed character (escape sequence `\n`).
There SHOULD be a “newline” at the end of the file.

## III. Appendix

### Glossary of technical terms

- “space”: The character ` ` (U+0020)
- “tab”: The tab character (U+0009), escape sequence `\t`
- “whitespace”: A “space”, a “tab”, or another character that appears blank
- “parenthesis”: The opening and closing parentheses `(` and `)` (U+0028 and U+0029)
- “newline”: The LF linefeed character (U+0010), escape sequence `\n`
- “blank line”: A line that only contains whitespace characters
- “letter”: A character as defined by the Unicode letter category, regex `\p{L}`
- “digit”: Any of 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
- “integer”: An unsigned number without fractional component
