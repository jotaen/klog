# Changelog (command line tool)

## v2.5
- **[ FEATURE ]** Restructure output of `klog now`, especially when
  using the `--diff` flag
- **[ FEATURE ]** Use distinct exit codes for different error cases
- **[ FEATURE ]** Introduce `--quiet` flag to retrieve raw output
- **[ FEATURE ]** Extend help texts, improve error messages 
- **[ FIX ]** Fix formatting issues of error output

## v2.4
- **[ FEATURE ]** Automatically create a new record when doing
  `klog start` or `klog track` if there is no record yet
- **[ FEATURE ]** Allow wildcard searching in tags by appending `...`,
  e.g. `--tag=foo...` would match the tag `#foobar`
- **[ FIX ]** `klog stop` now also looks for open ranges of the
  previous day and closes them with a shifted end time
- **[ FIX ]** `klog stop --summary=""` doesn’t fail if the existing
  entry summary was empty

## v2.3
- **[ FEATURE ]** Add `--summary`/`-s` flag for `start` and
  `stop` subcommands
- **[ FEATURE ]** If `KLOG_DEBUG` environment variable is set,
  print more verbose error output
- **[ FIX ]** Ensure that reading from stdin works on Windows
- **[ FIX ]** Display a more helpful error message on Windows
  to explain the quirks with `bookmark set`

## v2.2
- **[ FEATURE ]** Provide `--no-style` option to disable output
  formatting (i.e. no colours, underlined, bold, etc.)
- **[ FIX ]** Make sure that output formatting works on Windows
  across all Terminals.

## v2.1
- **[ FEATURE ]** Provide native Windows binary

## v2.0
- **[ BREAKING ]** Make `--after` and `--before` filters exclusive
- **[ FEATURE ]** Add commands for manipulating files:
  - `create` for creating a new record
  - `track` for adding an entry to a record
  - `start` to track an open-ended time range
  - `stop` to close an open-ended time range
- **[ FEATURE ]** Add `--since` and `--until` filters (inclusive)
- **[ FEATURE ]** Add `--period` filter (e.g. `--period=2015` for 
  all in 2015, or `--period=2015-04` for all in April 2015).

## v1.6
- **[ FEATURE ]** Add `json` subcommand that allows users to build
  programmatic extensions
- **[ FEATURE ]** Support Windows line endings (`\r\n`)
- **[ FEATURE ]** Add `bookmark unset` command for clearing current selection
- **[ FEATURE ]** Check stdin for input (to allow shell piping)

## v1.5
- **[ FIX ]** Fix the ongoing time counter in `klog now --follow`

## v1.4
- **[ FIX ]** Fix the ongoing time counter in the MacOS widget

## v1.3
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

## v1.2
- **[ INFO ]** Provided more helpful error messages
- **[ FIX ]** Fix unhandled error with experimental `template` subcommand
  (introduced in v1.1)

## v1.1
- **[ INFO ]** Introduced hidden and experimental `template` subcommand,
  see https://github.com/jotaen/klog/pull/12
- **[ FIX ]** If a duration consists hours and minutes,
  the minutes cannot be greater than `59m`, e.g. `1h59m`
- **[ FIX ]** Ensure there is a final blank line when `print`-ing
- **[ FIX ]** Improve error messages regarding the bookmark subcommand

## v1.0
- **[ BREAKING ]** Renamed subcommand `eval` to `total`.
  (This wording is more inline with the documentation and
  therefore more intuitive.)
- **[ FEATURE ]** Added subcommand `report` that generates a
  calendar overview
- **[ FEATURE ]** Added subcommand `tags` that shows the total
  times aggregated by tags
- **[ FEATURE ]** Added subcommand `bookmark` (a file that
  is used by default when no input files are specified)
