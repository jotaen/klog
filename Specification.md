# klog – File Format Specification

**Version 1.4**

klog is a file format for tracking time.

## License

Per [Creative Commons CC0 1.0 Universal](https://creativecommons.org/publicdomain/zero/1.0/),
to the extent possible under law, the editors have waived all copyright and related or
neighbouring rights to this work.
In addition, as of March 2022, the editors have made this specification available under the
[Open Web Foundation Agreement 1.0](https://www.openwebfoundation.org/the-agreements/the-owf-1-0-agreements-granted-claims/owfa-1-0).

## Preface

The keywords “MUST”, “MUST NOT”, “SHOULD”, “SHOULD NOT”, “RECOMMENDED”, “NOT RECOMMENDED” and “MAY”
in this document are to be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119).

Whenever a word has special meaning in klog, it is formatted in *italics*.

Other technical terms are surrounded by “quotes”. These are defined in the appendix.

Character sequences that are wrapped in `backticks` are meant to be read exactly (character by character).

## I. Record

A *record* is a self-contained data structure that contains time-tracking information.

Each *record* MUST appear as one consecutive block in the file,
without any “blank lines” appearing within.

The first line of a *record* MUST start with a *date*.
On the same line there MAY follow a *should-total*,
which MUST be separated by one “space” from the *date*
(additional “spaces” MAY appear).

A *record summary* MAY appear on the subsequent lines.
Any amount of *entries* MAY appear afterwards,
where each MAY have an *entry summary*.

In order to indent a line, it MUST start with one of the following sequences:
- Four “spaces” (RECOMMENDED)
- Two or three “spaces”
- One “tab”

To signify the second level of indentation,
the indentation sequence MUST appear twice.

The indentation style MUST be uniform within *records*.
(It MAY differ between *records*, though.)
[^indst]

### Date
A *date* is a day in the calendar.

> Examples: `2020-01-01`, `1984-08-30`, `2004/12/24`

*Dates* MUST contain 4 “digits” that denote the year,
2 “digits” that denote the month,
and 2 “digits” that denote the day.
The parts MUST be separated by either a `-` (RECOMMENDED)
or a `/`.
The year part MUST be written first, then the month, then the day.

The combination of year, month and day MUST be representable by the Gregorian calendar.

### Should-Total
A *should-total* denotes the targeted total time of a *record*.

> Examples: `(8h!)`, `(5h15m!)`, `(-3h30m!)`

A *should-total* MUST be a *duration* value
followed by a `!`
and wrapped in “parentheses”.

### Summary
A *summary* is user-provided text for capturing arbitrary information
about a *record* or an *entry*. *Summaries* are optional.

#### Record Summary
The *record summary* is considered to be associated with the entire *record*.

It MUST appear underneath the *date*,
and it MAY span multiple lines.
Each of its lines MUST NOT start with “blank characters”.
[^resui]

#### Entry Summary
The *entry summary* is considered to be referring to one particular *entry*.

It MUST either start on the same line as the *entry*,
separated from it by one “space”;
or it MUST start on the subsequent line.

The *entry summary* MAY span multiple lines.
All lines following the *entry* line MUST be indented twice;
they also MUST NOT only consist of “blank characters”.
[^iwses]

#### Tag
The purpose of *tags* is to help categorise *records* and *entries*.

> Examples: `#gym`, `#home-office`, `#読む`, `#ticket=891`, `#project="22/48.3"`

Any amount of *tags* MAY appear anywhere within *summaries*.

A *tag* MUST consist of a *tag name*,
which MUST be preceded by a single `#` character.
The *tag name* MAY be followed by a `=` character
and a *tag value*.

The *tag name* MUST only contain
“letters”, “digits”, or the characters `_` or `-`.
It MUST be interpreted as if it was all lower-case.
[^csitn]

The *tag value* MAY be surrounded by a pair of matching quotes,
which MUST either be `"` (RECOMMENDED) or `'`.
- If the *tag value* is quoted, it MAY contain any character
  except for the respective quote character itself,
  or a “newline”.
  In case no matching closing quote appears on the same line,
  the *tag value* MUST be treated as absent.
  [^qutvl]
- If the *tag value* is not quoted, it MUST only contain
  “letters”, “digits”, or the characters `_` or `-`.

An empty *tag value* (e.g. `#tag=` or `#tag=""`)
MUST be treated the same as an absent *tag value* (e.g. `#tag`).

### Entry
*Entry* is an abstract term for time-related data.
*Durations*, *ranges* and *open ranges* are instances of *entries*.

> Examples (indentation omitted): `2h30m`, `-1h Lunch break`, `11:00 - 14:15`, `8:00am - 2:00pm Long day at #school`

Each *entry* MUST appear on its own line and
MUST be indented once.

A *summary* MAY be associated with an *entry* (see section Summary).

### Time
A *time* is a value that represents a point in time throughout a day
as it would be displayed by a wall clock (which divides a day into
24 hours and every hour into 60 minutes).

> Examples: `14:18`, `6:30am`, `01:00>`, `<23:00am`

*Time* values MUST contain an hour part and a minute part,
separated by a `:` in between.
The hour part MUST be written first.

As default, *times* are to be interpreted as 24-hour clock values.
An `am` or `pm` suffix MAY be used to denote that the value is
to be interpreted as 12-hour clock value.

The minute part MUST be between 0-59 (inclusive).
Single-figure minute parts MUST be padded with a `0`.

The hour part MUST either be between 0-24 (inclusive) when using the 24-hour clock,
or between 1-12 (inclusive) when using the 12-hour clock.
Single-figure hour parts MAY be padded with a `0`.

When using the 24-hour clock, if the hour part is `24`,
then the minute part MUST be `00`;
`<24:00` MUST be interpreted as `0:00`,
`24:00` MUST be interpreted as `0:00>`,
`24:00>` MUST NOT appear.

*Time* values MAY be *shifted* to the next or to the previous day:
- To associate the *time* with the day before the *record’s* *date*,
  a `<` prefix MUST be used,
  e.g. `<23:00`.
- To associate the *time* with the day after the *record’s* *date*,
  a `>` suffix MUST be used,
  e.g. `1:30>`.

### Range
A *range* is an *entry* that represents the time span between two points in time.

> Examples: `8:00 - 9:00`, `11:00am - 1:00pm`, `<23:40 - 3:12`, `0:30> - 4:00>`

*Ranges* MUST contain two *time* values that denote the start and the end.
Start *time* and end *time* MUST be written in chronological order.
They MAY be equal.

There MUST be a `-` between the two values.
There MAY appear “spaces” on either side of the `-`,
in which case it is RECOMMENDED to use exactly one “space” on both sides of the dash.

### Open range
An *open range* is an *entry*
that can be used to track the start *time* of an activity,
i.e. the end *time* is not determined yet.

> Examples: `05:17 - ?`, `4:00pm - ?`

*Open ranges* are formatted in the same way as *ranges*,
except that the end *time* MUST be replaced by a placeholder.
The placeholder MUST be denoted by the character `?`,
e.g. `9:00 - ?`. 
The `?` MAY be repeated, e.g. `9:00 - ???`.
[^plrep]
The placeholder MUST NOT be *shifted*.

*Open ranges* MUST NOT appear more than once per *record*.
[^oasor]

### Duration
A *duration* is an *entry* that represents a period of time.

> Examples: `1h`, `5m`, `4h12m`, `-8h30m`

*Durations* MUST contain an amount of hours and/or an amount of minutes.
(So they MUST either contain one of these two or both.)
The hour part MUST be written first.

The hour part MUST be an “integer”
which MUST be followed by the character `h`.
It MAY be `0h`.
It MAY be greater than `24h`,
e.g. `50h`.
If the hour part is missing, a value of `0h` is assumed.

The minute part MUST be an “integer”
which MUST be followed by the character `m`.
It MAY be `0m`.
When the hour part is present,
the minute part MUST NOT be greater than `59m`,
e.g. `1h59m`;
otherwise it MAY be greater than `59m`,
e.g. `119m`
(it is RECOMMENDED to break this up, though).
If the minute part is missing, a value of `0m` is assumed.

The *duration* as a whole is a signed value:
That means it is either positive (i.e. adding to the *total time*)
or negative (i.e. deducting from the *total time*).
By default, a *duration* is positive,
which MAY be indicated by a leading `+` character,
e.g. `+4h12m`.
If the *duration* is supposed to be negative, it MUST be preceded by a `-` character.

## II. Organising records in files

A file MAY hold any amount of *records*.

There MUST appear one “blank line” between subsequent *records*;
additional “blank lines” MAY appear.

*Records* MAY appear in any order in the file.

There MAY exist multiple *records* with the same *date*.

A file MUST NOT contain anything but what is allowed by this specification.
Otherwise, it SHOULD NOT be evaluated.
[^fcocr]

The file extension SHOULD be `.klg`, e.g. `times.klg`.
The file encoding MUST be UTF-8.

“Newlines” MUST be encoded with either the
linefeed character (LF, escape sequence `\n`),
or carriage return and linefeed character (CRLF, escape sequences `\r\n`).
These two styles SHOULD NOT be mixed within the same file.

There SHOULD be a “newline” at the end of the file.

## III. Evaluating data

### Total time
The resulting *total time* of a *record* MUST be computed by summing up its *entries*:
positive values add to the *total time*,
negative values deduct from it.
The resulting *total time* MAY be 0;
it MAY be negative;
it MAY be greater than 24 hours.

Overlapping *ranges* MUST each be counted fully.
E.g., the two *entries* `12:00 - 13:00` and `12:30 - 13:30` result in a *total time* of `2h`.

*Ranges* with *shifted times* MUST be fully counted towards
the *date* at which they appear in the *record*.
They MUST NOT be implicitly split across the two adjacent *dates*.

*Open ranges* MUST NOT be counted by default;
they MAY be factored in upon explicit request, though.

Multiple *records* with the same *date* MUST be treated as distinct
and MUST NOT be combined into a single *record*.

## IV. Appendix

### Glossary of technical terms

- “space”: The character ` ` (U+0020)
- “tab”: The tab character (U+0009, escape sequence `\t`)
- “blank character”: A “tab”, or a character as defined by the Unicode Space Separator category (Zs)
- “blank line”: A line that only contains “blank characters”
- “parenthesis”: The opening and closing parentheses `(` and `)` (U+0028 and U+0029)
- “letter”: A character as defined by the Unicode Letter category (L)
- “digit”: Any of 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
- “integer”: An unsigned number without fractional component
- “newline”: Either a linefeed (U+0010, escape sequence `\n`), or a carriage return and linefeed (U+0013 and U+0010, escape sequence `\r\n`)

### Changelog

#### Version 1.4
- Release the specification document under the CC0/OWFa license.
- Support for tags to (optionally) have values assigned to them.
- Allow hyphens (`-`) to appear in tags.
- Add footnotes to make context information explicit.

#### Version 1.3
- Specify additional rules for multiline entry summaries.

#### Version 1.2
- Allow times to be `24:00`.
- Some minor restructurings for enhanced clarity.

#### Version 1.1
- Add a constraint regarding the indentation that requires the indentation style
  to be uniform within a record.
- Remove technical term “whitespace”, since its meaning is ambiguous and the definition lacked clarity.
  Replace it with “blank character” and base the definition on the Unicode category.

### Footnotes

The following footnotes are purely informational,
to make contextual background information explicit.

[^indst]: The indentation must be uniform, otherwise the levels can’t be determined
    unambiguously. E.g., if 4 spaces are encountered at the beginning of the line, it would
    be unclear whether that is 2 * 2 spaces or 1 * 4 spaces.
[^resui]: Lines in the record summary can’t start with blank characters, to avoid that they
    might be visually confused with the (indented) entries. There is no strict technical
    reason for this, though.
[^iwses]: In contrast to record summaries, lines in entry summaries can start with blank
    characters. That is for allowing the user to vertically align the summary text on all
    entry lines. A by-effect of this rule is that there can never be a third indentation level.
[^csitn]: The character set that a tag is allowed to consist of is deliberately limited,
    so that tags can appear as natural words in the flow of a sentence. E.g.:
    `#Office day (#coding, #meetings)`. That’s also why tag names are to be interpreted
    as case-insensitive. (Tag values, on the other hand, are always to be interpreted literally.)
[^qutvl]: The main use-case for quoted tag values is for literal references, such as a project id,
    or a name: `#project="2022/7.2"` or `#call="Liz Jones"`. That’s also why tag values
    are always to be interpreted as case-sensitive (in contrast to tag names).
[^plrep]: The `?` placeholder in open ranges can be repeated, to allow users to visually
    align it with other entries. E.g. `8:00-?????` has the same width as `8:00-9:00`.
[^oasor]: Open ranges only being allowed to appear once per record has a mere practical motivation:
    it’s important for making interactions with tools easier. Otherwise, when stopping activities
    via a tool, it might be ambiguous which of the open ranges is meant.
[^fcocr]: By allowing a file to only contain records and nothing else, a klog file can effectively
    be perceived as a text-based database. That makes it easy to process files programmatically,
    because every record is a self-contained and strictly structured unit of data.
