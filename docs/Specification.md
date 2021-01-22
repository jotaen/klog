# klog – File Format Specification (DRAFT)

klog is a file format for tracking times.
It is free and open-source software distributed under the MIT-License.

## Preface

The keywords “MUST”, “MUST NOT”, “SHOULD”, “SHOULD NOT”, “RECOMMENDED”, “NOT RECOMMENDED” and “MAY”
in this document are to be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119).

Whenever a word has special meaning in klog, it is formatted in *italics*.

Other technical terms are surrounded by “quotes”. These are defined at the end of this specification.

## I. Records

A *record* is a self-contained and atomic unit of data.

Each *record* MUST appear as one consecutive block in the file,
without any “blank lines” appearing within.

The first line of a *record* MUST start with a *date*.
On the same line there MAY follow a list of *properties*.
This list MUST be enclosed in “parentheses” and
preceded by at least one “space”.
If there are multiple *properties* they MUST be separated
by a `,` followed by a “space”.
The list of *properties* MUST NOT be empty.

A *summary* MAY appear on the subsequent lines.
Any amount of *entries* MAY appear afterwards.

### Date
A *date* is a day that is representable in the Gregorian calendar.

Each *record* MUST contain a *date*.

It MUST be either formatted according to one of the following patterns:
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

### Summary
A *summary* is user-provided text for holding arbitrary information.
There are two places where *summary* text MAY appear in *records*:

- After the *date*:
  In this case the *summary* is considered to be associated with the entire *record*.
  The *summary* MAY span multiple lines.
  Each of its lines MUST NOT start with “whitespace”.
- Behind *entries*:
  In this case the *summary* is only considered to be referring to the corresponding *entry*.
  The *summary* text follows the *entry* on the same line,
  and it ends at the end of that line.
  It MUST be separated from the *entry* by at least one “space”.
  It MUST start with a “letter” or a `#`.

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

There MUST be a `-` between the two values,
which MAY be surrounded by one “space” on each side.

The start value MUST be a *time*.
It MAY be prefixed with a `<` to indicate that
this *time* is referring to the day before the *record’s* date,
e.g. `<23:00`.

The end value MAY be a *time*;
it MAY also be substituted by one or more `?` to denote that the end is not determined yet.
This MUST NOT occur more than once per record.
If the end value is a *time* it MAY be suffixed with a `>` to indicate
that this *time* is referring to the day after the *record’s* date,
e.g. `0:30>`.

### Duration
A *duration* is an *entry* that represents a period of time.
It contains an amount of hours and/or an amount of minutes.
(So it MUST contain either one of these two or both.)
The hour part MUST be written first.
Examples are: `1h`, `5m`, `4h12m`, `-8h30m`.

The hour part MUST be an unsigned number
which MUST be immediately followed by the character `h`.
It MAY be `0h`.
It MAY be greater than `24h`,
e.g. `50h`.
If the hour part is missing, a value of `0h` is assumed.

The minute part MUST be an unsigned number
which MUST be immediately followed by the character `m`.
It MAY be `0m`.
It MAY be greater than `59m`,
e.g. `150m`;
it is generally RECOMMENDED to break this up, though.
If the hour part is missing, a value of `0m` is assumed.

While the hour and minute parts itself are unsigned,
the *duration* as a whole is always a signed value:
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

Subsequent *records* MUST be separated by one “blank line”;
there MAY be additional blank lines.

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
- “digit” Any of 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
