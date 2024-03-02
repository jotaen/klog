# Changelog
**Summary of changes of the command line tool**

## Next Up
- **[ FEATURE ]** There is a new setting for the klog `config.ini` file,
  which allows to specify the colour theme of the terminal, so that klog
  can optimise its output colouring. The available options are: `dark` (the
  default), `light`, and `no_colour`. Run `klog config` to learn more.
- **[ FEATURE ]** Add two new rounding options: `12m` and `20m`.
  (E.g., when using `klog start --round`.)
- **[ FIX ]** Fix formatting bug of `klog print --with-totals` for
  multiline record summaries.

## v6.2 (2023-10-17)
- **[ FEATURE ]** Add new command `klog switch`, that stops a previously
  ongoing activity (open time range), and starts a new one.
- **[ FEATURE ]** `klog start --resume` now falls back to the previous
  record for determining the last entry summary.

## v6.1 (2023-05-01)
- **[ FEATURE ]** Add new flag `klog start --resume`, which takes over the
  summary of the last entry for the new open-ended entry.

## v6.0 (2023-03-06)
- **[ BREAKING ]** The default location of the klog config folder has moved!
  So far, that folder only contains the bookmark database, so if you don’t
  use bookmarks, you can ignore this change. In order to see or configure the
  location of the klog config folder, please run `klog info config-folder`
  (on the new release). The previous location was `~/.klog` on all systems,
  so you might have to manually move over the contents of that previous folder,
  and potentially adjust your dotfile management (if you have that).
- **[ FEATURE ]** Introduce optional, file-based configuration for general
  preferences such as the preferred date or time format, or default values for
  certain CLI flags. See `klog config` / `klog config --help` to learn more.
- **[ FEATURE ]** Display warning when using `--now` without there being any
  open-ended time range in the data.

## v5.4 (2022-11-23)
- **[ BREAKING ]** For `klog edit`, support if the `$EDITOR` variable
  contains additional flags, e.g. `vi -R` or `subl -w`. (If your editor
  path contains spaces, you now have to wrap it in quotes.)
- **[ BREAKING ]** Simplify logic of `klog pause` command; add `--extend`
  flag for extending a previous pause.
- **[ FIX ]** For `klog pause`, recover correctly after computer had
  been asleep.
- **[ FEATURE ]** For `klog pause`, display current record while pausing.
- **[ FEATURE ]** For `klog tags`, optionally display how many entries
  there are per tag via `--count`.
- **[ FEATURE ]** For `klog json`, provide `--now` flag.
- **[ FEATURE ]** For file manipulation commands (e.g. `klog track`),
  improve automatic detection of style preferences in the file.
- **[ INFO ]** Significantly improve parsing performance for large
  data inputs (i.e., for files with 1000+ records).

## v5.3 (2022-10-31)
- **[ FEATURE ]** Optionally amend `klog print` output with total
  values via the `--with-totals` flag.
- **[ FIX ]** Fix unhandled error edge-case in `klog report` command.

## v5.2 (2022-08-19)
- **[ FEATURE ]** Provide tab completion functionality for bash, zsh
  and fish shell. Run `klog completion` for setup instructions.
- **[ FIX ]** `klog edit` handles when the `$EDITOR` variable contains
   spaces, and it also fails when `$EDITOR` is invalid.

## v5.1 (2022-07-20)
- **[ FEATURE ]** Optionally print out totals as decimal values (in minutes)
  via the `--decimal` flag; e.g. `150` instead of `2h30m`.
- **[ FEATURE ]** Support `--now` on `klog tags` as well.
- **[ FEATURE ]** Allow setting a record summary via the `--summary` flag
  when using `klog create`.

## v5.0 (2022-04-13)
- **[ META ]** Release the klog file format specification into the public domain
  (under the CC0/OWFa dual license).
  Read it here: https://github.com/jotaen/klog/blob/main/Specification.md
- **[ FEATURE ]** Allow tags to (optionally) have values assigned to them,
  e.g. `#ticket=1764` or `#type=work`. The values can be quoted if
  they contain special characters: `#project="22/48.3"`.
- **[ FEATURE / BREAKING ]** Allow hyphens (`-`) to appear in tags, e.g. `#home-office`.
- **[ FEATURE ]** For the `--period` flag, additionally allow filtering
  by quarter (`YYYY-Qq`, e.g. `2022-Q1`) and week (`YYYY-Www`, e.g. `YYYY-W34`).

## v4.0 (2022-03-21)
- **[ FEATURE ]** Allow summaries behind entries to be continued on
  the next line (with increased indentation level), e.g.:
  ```
  2020-01-01
  Both of the following is fine:
      15:00-16:00 This is a very long text, so
          it can be continued on the next line.
      16:00-17:00
          Or, you can just start the entry summary
          on the next line, if you like.
  ```
  The CLI also handles this automatically when it encounters
  line breaks (`\n`), e.g. in the `--summary` flag value.
- **[ FEATURE ]** Add new command `klog pause` that “pauses”
  open-ended time ranges by adding a subsequent pause entry.
- **[ FEATURE ]** Provide rounding option for `klog start` and
  `klog stop`, which rounds times to the nearest multiple of
  5m, 10m, 15m, 30m, or 60m. E.g. for `--round=15m`: `8:03` -> `8:00`.
- **[ FEATURE ]** Add more shortcut filters, e.g. `--this-week`,
  `--last-month`, etc.
- **[ FEATURE ]** Embed the most recent part of the changelog for
  convenience, via `klog --changelog`.
- **[ BREAKING ]** Remove embedded macOS systray widget

## v3.3 (2022-01-30)
- **[ FEATURE ]** Allow times to be `24:00`, e.g. `22:00 - 24:00`.
- **[ FEATURE ]** Add `klog goto` command for opening the file explorer
  at the location of a file or bookmark.
- **[ FEATURE ]** Add `klog bookmark info` command.
- **[ FEATURE ]** When using the manipulation commands (`klog track`, etc.),
  conform to style preferences encountered in the file.
- **[ FEATURE ]** Add `--tomorrow` as shorthand flag for the next day’s date.
- **[ FEATURE ]** Improve warnings (which are shown for potential data problems).
- **[ FIX ]** Fix bug in week-based aggregation of `klog report --aggregate week`.

## v3.2 (2021-11-30)
- **[ BREAKING ]** Don’t allow mixing the indentation style within a
  record. (It might still differ *between* records, though.) For example: if
  the first entry is indented with a tab, then all further entries of that
  particular record have to be indented with a tab as well. In order to check
  that your existing files conform, you can parse all your `.klg` files at once
  via a wildcard lookup, in order to see whether any indentation-related
  errors are reported. On Linux, e.g.: `klog total ~/**/*.klg`.
- **[ FEATURE ]** Allow version check via `klog -v` (in addition
  to `klog --version` or `klog version`)
- **[ FEATURE ]** Embed specification and license in the binary
  (via `klog --spec` and `klog --license`)
- **[ FEATURE ]** Provide binaries for M1 Macs (ARM) for download.
- **[ FIX ]** Fix default sort order of `--sort` flag to be `asc`
- **[ INFO ]** Deprecate the embedded native widget (for MacOS). It will be
  removed in one of the next releases.

## v3.1 (2021-10-20)
- **[ FIX ]** Fix stdin processing on Windows

## v3.0 (2021-10-07)
- **[ FEATURE ]** Support multiple (named) bookmarks to quickly
  reference often-used files, e.g. `klog total @work`
- **[ FEATURE ]** Add additional evaluation options for `klog report`
  to aggregate the data by day, week, month, quarter or year
- **[ FEATURE ]** Add `klog edit` command for opening a file in an editor
  (Based on the `$EDITOR` variable.)
- **[ FEATURE ]** Allow value of `--sort` flag to be uppercase
  or lowercase (`ASC`/`asc` or `DESC`/`desc`)
- **[ FEATURE ]** Support `klog --version` in addition to `klog version`
- **[ FIX ]** Windows: don’t require admin privileges for setting bookmarks 

## v2.6 (2021-07-25)
- **[ INFO ]** Release first version of the file format
  specification (v1.0)
- **[ FIX ]** If a duration only contains a minute part,
  allow the value to be greater than 59, e.g. `120m`.

## v2.5 (2021-05-17)
- **[ BREAKING ]** Rename `klog now` to `klog today`; restructure the
  output, especially when using the `--diff`/`--now` flag
- **[ FEATURE ]** Use distinct exit codes for different error cases
- **[ FEATURE ]** Introduce `--quiet` flag to retrieve raw output
- **[ FEATURE ]** Extend help texts, improve error messages 
- **[ FIX ]** Fix formatting issues of error output

## v2.4 (2021-05-05)
- **[ FEATURE ]** Automatically create a new record when doing
  `klog start` or `klog track` if there is no record yet
- **[ FEATURE ]** Allow wildcard searching in tags by appending `...`,
  e.g. `--tag=foo...` would match the tag `#foobar`
- **[ FIX ]** `klog stop` now also looks for open ranges of the
  previous day and closes them with a shifted end time
- **[ FIX ]** `klog stop --summary=""` doesn’t fail if the existing
  entry summary was empty

## v2.3 (2021-04-28)
- **[ FEATURE ]** Add `--summary`/`-s` flag for `start` and
  `stop` subcommands
- **[ FEATURE ]** If `KLOG_DEBUG` environment variable is set,
  print more verbose error output
- **[ FIX ]** Ensure that reading from stdin works on Windows
- **[ FIX ]** Display a more helpful error message on Windows
  to explain the quirks with `bookmark set`

## v2.2 (2021-04-03)
- **[ FEATURE ]** Provide `--no-style` option to disable output
  formatting (i.e. no colours, underlined, bold, etc.)
- **[ FIX ]** Make sure that output formatting works on Windows
  across all Terminals.

## v2.1 (2021-03-19)
- **[ FEATURE ]** Provide native Windows binary

## v2.0 (2021-03-16)
- **[ BREAKING ]** Make `--after` and `--before` filters exclusive
- **[ FEATURE ]** Add commands for manipulating files:
  - `create` for creating a new record
  - `track` for adding an entry to a record
  - `start` to track an open-ended time range
  - `stop` to close an open-ended time range
- **[ FEATURE ]** Add `--since` and `--until` filters (inclusive)
- **[ FEATURE ]** Add `--period` filter (e.g. `--period=2015` for 
  all in 2015, or `--period=2015-04` for all in April 2015).

## v1.6 (2021-03-06)
- **[ FEATURE ]** Add `json` subcommand that allows users to build
  programmatic extensions
- **[ FEATURE ]** Support Windows line endings (`\r\n`)
- **[ FEATURE ]** Add `bookmark unset` command for clearing current selection
- **[ FEATURE ]** Check stdin for input (to allow shell piping)

## v1.5 (2021-02-16)
- **[ FIX ]** Fix the ongoing time counter in `klog now --follow`

## v1.4 (2021-02-16)
- **[ FIX ]** Fix the ongoing time counter in the MacOS widget

## v1.3 (2021-02-14)
- **[ BREAKING ]** Change structure of the bookmark subcommand
  (This is in order to account for the increasing number of operations)
- **[ FEATURE ]** Add subcommand `now` for displaying an ongoing total
  that takes open ranges into account (based on the time of execution)
- **[ FEATURE ]** Add `--now` flag to `total` and `report` to take
  open ranges into account optionally
- **[ FEATURE ]** Add subcommand `bookmark edit` for opening a bookmarked
  file in your $EDITOR
- **[ FEATURE ]** Allow to sort results in both directions
  (`--sort ASC` or `--sort DESC`)
- **[ FEATURE ]** Print warning when unclosed open ranges are detected
  in records before yesterday. (It’s probably always a mistake, if that occurs.)
  You can disable this check with the `--no-warn` flag.
- **[ FEATURE ]** Support filtering for `tags` and `reports`
- **[ FEATURE ]** Define shorthand flags, e.g. for `--now`, `--diff`
- **[ FIX ]** Don’t demand `.klg` file extensions for bookmarks

## v1.2 (2021-02-07)
- **[ INFO ]** Provided more helpful error messages
- **[ FIX ]** Fix unhandled error with experimental `template` subcommand
  (introduced in v1.1)

## v1.1 (2021-02-07)
- **[ INFO ]** Introduced hidden and experimental `template` subcommand,
  see https://github.com/jotaen/klog/pull/12
- **[ FIX ]** If a duration consists hours and minutes,
  the minutes cannot be greater than `59m`, e.g. `1h59m`
- **[ FIX ]** Ensure there is a final blank line when `print`-ing
- **[ FIX ]** Improve error messages regarding the bookmark subcommand

## v1.0 (2021-02-06)
- **[ BREAKING ]** Renamed subcommand `eval` to `total`.
  (This wording is more inline with the documentation and
  therefore more intuitive.)
- **[ FEATURE ]** Added subcommand `report` that generates a
  calendar overview
- **[ FEATURE ]** Added subcommand `tags` that shows the total
  times aggregated by tags
- **[ FEATURE ]** Added subcommand `bookmark` (a file that
  is used by default when no input files are specified)
