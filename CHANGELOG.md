# Changelog (command line tool)

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
