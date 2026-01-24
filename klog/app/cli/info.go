package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/prettify"
)

type Info struct {
	Spec      bool `name:"spec" help:"Print the .klg file format specification."`
	License   bool `name:"license" help:"Print license / copyright information."`
	About     bool `name:"about" help:"Print meta information about klog."`
	Filtering bool `name:"filtering" help:"Print documentation for using filter expressions."`
}

func (opt *Info) Run(ctx app.Context) app.Error {
	text, err := func() (string, app.Error) {
		if opt.Spec {
			return ctx.Meta().Specification + "\n", nil
		} else if opt.License {
			return ctx.Meta().License + "\n", nil
		} else if opt.About {
			return INTRO_SUMMARY, nil
		} else if opt.Filtering {
			return `klog filter expressions are a generic way for filtering data for evaluation purposes. Use the --filter flag to specify a filter expression, e.g.:

    klog total --filter='2025-04 && #work' mytimes.klg

Wrap the filter expression in single quotes to avoid undesired shell word splitting or substitution. Filter expressions consist of operands for matching the data that shall be included in the filter result. Operands can be combined via logical operators and grouped via parentheses.

Examples:
    2025-04-20 || 2020-04-21
        All entries at either 2025-04-20 or 2020-04-21.
    2025-04 && !#work
        All entries in April 2025 that don’t match tag #work.
    2025-04-15...2025-05-30 && (#gym || #run)
        All entries since 2025-04-15 until 2025-05-30 (inclusive), that match either tags #gym or #run.

Operators:
    (  )
        Group operands with parentheses. Example: (#foo || #bar) && 2020-04
    &&
    ||
        Combine operands logically, either as AND or OR. Note: you cannot mix different operators within the same group.
        Examples:
            #foo && 2020-04
            #foo && (2020-04 || 2020-05)
    !
        Negate operands to exclude data. Example: !2020-04 && !#foo

Operands:
    YYYY-MM-DD
        Entries at that date. Example: 2025-04-30
    YYYY-MM
        Entries in that month. Example: 2025-04
    YYYY-Wxx
        Entries in that week. Example: 2025-W34
    YYYY-Qx
        Entries in that quarter. Example: 2025-Q4
    YYYY
        Entries in that year. Example: 2025
    YYYY-MM-DD...YYYY-MM-DD
    YYYY-MM-DD...
    ...YYYY-MM-DD
        Entries within that date range.
        Examples:
            2025-04-30...2025-05-14
            2025-04-30...
    #tag
    #tag=value
        Entries matching that a tag. Note that tags can either be specified in the entry summary or in the record summary. In the latter case, the entry “inherits” the record tags.
        Examples: #work || #project=467 || #project='#312'
    type:xxx
        Entries of that type, where xxx can be either:
        range, open-range, duration, duration-positive, duration-negative
        Example: type:duration
`, nil
		} else {
			return "", app.NewErrorWithCode(
				app.GENERAL_ERROR,
				"No flag specified",
				"Run with `--help` for more info",
				nil,
			)
		}
	}()

	if err != nil {
		return err
	}

	ctx.Print(prettify.Reflower.Reflow(text, ""))
	return nil
}
