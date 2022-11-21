package lib

import (
	"github.com/jotaen/klog/klog"
	tf "github.com/jotaen/klog/klog/app/cli/lib/terminalformat"
	"github.com/jotaen/klog/klog/parser"
	"strconv"
	"strings"
)

// CliSerialiser is a specialised parser.Serialiser implementation for the terminal.
type CliSerialiser struct {
	Unstyled bool // -> No colouring/styling
	Decimal  bool // -> Decimal values rather than the canonical totals
}

var (
	Green     = tf.Style{Color: "120"}
	Red       = tf.Style{Color: "167"}
	BlueDark  = tf.Style{Color: "117"}
	BlueLight = tf.Style{Color: "027"}
	Subdued   = tf.Style{Color: "249"}
	Purple    = tf.Style{Color: "213"}
)

func (cs CliSerialiser) Format(s parser.Styler, t string) string {
	if cs.Unstyled {
		return t
	}
	return s.Format(t)
}

func (cs CliSerialiser) formatAndRestore(s tf.Style, prev tf.Style, t string) string {
	if cs.Unstyled {
		return t
	}
	return s.FormatAndRestore(t, prev)
}

func (cs CliSerialiser) duration(d klog.Duration, withSign bool) string {
	if cs.Decimal {
		return strconv.Itoa(d.InMinutes())
	}
	if withSign {
		return d.ToStringWithSign()
	}
	return d.ToString()
}

func (cs CliSerialiser) Date(d klog.Date) string {
	return cs.Format(tf.Style{Color: "015", IsUnderlined: true}, d.ToString())
}

func (cs CliSerialiser) ShouldTotal(d klog.Duration) string {
	return cs.Format(Purple, cs.duration(d, false))
}

func (cs CliSerialiser) Summary(s parser.SummaryText) string {
	txt := s.ToString()
	style := Subdued
	hashStyle := style.ChangedBold(true).ChangedColor("251")
	txt = klog.HashTagPattern.ReplaceAllStringFunc(txt, func(h string) string {
		return cs.formatAndRestore(hashStyle, style, h)
	})
	return cs.Format(style, txt)
}

func (cs CliSerialiser) Range(r klog.Range) string {
	return cs.Format(BlueDark, r.ToString())
}

func (cs CliSerialiser) OpenRange(or klog.OpenRange) string {
	return cs.Format(BlueLight, or.ToString())
}

func (cs CliSerialiser) Duration(d klog.Duration) string {
	f := Green
	if strings.HasPrefix(d.ToStringWithSign(), "-") {
		f = Red
	}
	return cs.Format(f, cs.duration(d, false))
}

func (cs CliSerialiser) SignedDuration(d klog.Duration) string {
	f := Green
	if strings.HasPrefix(d.ToStringWithSign(), "-") {
		f = Red
	}
	return cs.Format(f, cs.duration(d, true))
}

func (cs CliSerialiser) Time(t klog.Time) string {
	return cs.Format(BlueLight, t.ToString())
}
