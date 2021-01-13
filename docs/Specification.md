# klog Specification

*klog* is a file format for personal time tracking.

## Preface

The key words “MUST”, “MUST NOT”, “SHOULD”, “SHOULD NOT”, and “MAY” in this document are to be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119).


## Concepts

### Record
A *record* 
It MUST contain a *date*;
it MAY contain *entries* (i.e. *ranges* and/or *durations*);

### Date
SHOULD be `YYYY-MM-DD`, MAY be `YYYY/MM/DD`

### Summary
A *summary* is text for holding user-provided information.
It MAY appear as a block above the *entries* of a record,
or it MAY appear inline at the end of *entries*.

*Summary* text MUST be commenced with a `|` character.

### Tags

### Entry
*Entry* is the general term for time-related data, i.e. either a *range* or a *duration*.
Each *entry* MAY 

### Time
A *time* is a value that represents a point in time throughout a day as it would be displayed by a digital wall clock.
It MUST consist of an hour part and a minute part.

The hour part MUST be a number between `0` and `23` (both inclusive).
Single digit numbers MAY be preceded with a leading `0`.

The minute part MUST be a number between `0` and `59` (both inclusive).
Single digit numbers MUST be preceded with a leading `0`.

### Range
...
<8:00 etc.

### Duration
A *duration* represents a time span.
It contains an amount of hours and/or an amount of minutes.
(So it MUST contain either one of these two or both.)
The hour part MUST be written first.
The two parts MAY be separated by one space character (` `).

The hour part MUST be an unsigned number
which MUST be immediately followed by the character `h`.
It MAY be `0h`.
(Though it SHOULD be omitted then.)
It MAY be greater than `24h`.
If the hour part is missing, a value of `0h` is assumed.

The minute part MUST be an unsigned number
which MUST be immediately followed by the character `m`.
It MAY be `0m`.
(Though it SHOULD be omitted then.)
It MAY be greater than `59m`.
(Though it SHOULD be broken down into hours then.)
If the hour part is missing, a value of `0m` is assumed.

The *duration* is always a signed value:
That means it is either positive (i.e. additional) or negative (i.e. deductible).
As default a *duration* is positive,
which MAY be indicated by a leading `+` character.
If the *duration* is supposed to be negative, it MUST be preceded by a `-` character.
There MUST NOT be whitespace between the sign and the *duration*.

#### Examples
- `1h`
- `5m`
- `4h12m`
- `-4h12m`
- `+4h12m` (leading `+` MAY appear)
- `4h 12m` (Space character MAY appear between parts)
- `0h 0m` (not recommended; SHOULD BE omitted altogether, if possible)
- `150m` (not recommended; SHOULD BE `2h30m`)


## File layout

In order to be persisted, *records* MUST be stored in plain text files.
The file extension MUST be `.klg`, e.g. `times.klg`.
The file encoding MUST be UTF-8.
Newlines MUST be encoded through the sequence `\n`.

A file MAY hold any amount of records (including none).
It MUST NOT contain anything but what is allowed by this specification.

Each *record* appears as one consecutive block without any blank lines within.
Subsequent *records* SHOULD be divided by one blank line;
there MAY be more than one blank line.
(A blank line is a line that is either empty or contains whitespace exclusively.)

The first line of a *record* MUST start with a *date*.
There MUST NOT appear any preceding whitespace.

The subsequent lines MAY contain the *summary* and/or the *entries*.
These MUST be indented; the indentation
SHOULD be one “tab” sequence (`\t`),
or MAY be two or more “space” characters (` `).
