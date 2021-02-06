package cli

import (
	"klog"
	. "klog/lib/jotaen/tf"
	"klog/parser"
)

var styler = parser.Serialiser{
	Date: func(d klog.Date) string {
		return Style{Color: "015", IsUnderlined: true}.Format(d.ToString())
	},
	ShouldTotal: func(d klog.Duration) string {
		return Style{Color: "213"}.Format(d.ToString())
	},
	Summary: func(s klog.Summary) string {
		txt := s.ToString()
		style := Style{Color: "249"}
		hashStyle := style.ChangedBold(true).ChangedColor("251")
		txt = klog.HashTagPattern.ReplaceAllStringFunc(txt, func(h string) string {
			return hashStyle.FormatAndRestore(h, style)
		})
		return style.Format(txt)
	},
	Range: func(r klog.Range) string {
		return Style{Color: "117"}.Format(r.ToString())
	},
	OpenRange: func(or klog.OpenRange) string {
		return Style{Color: "027"}.Format(or.ToString())
	},
	Duration: func(d klog.Duration, forceSign bool) string {
		f := Style{Color: "120"}
		if d.InMinutes() < 0 {
			f.Color = "167"
		}
		if forceSign {
			return f.Format(d.ToStringWithSign())
		}
		return f.Format(d.ToString())
	},
}
