package app

import (
	"github.com/jotaen/klog/klog"
	tf "github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/parser"
	"strconv"
	"strings"
)

// TextSerialiser is a specialised parser.Serialiser implementation for the terminal.
type TextSerialiser struct {
	DecimalDuration bool
	Styler          tf.Styler
}

func NewSerialiser(styler tf.Styler, decimal bool) TextSerialiser {
	return TextSerialiser{
		DecimalDuration: decimal,
		Styler:          styler,
	}
}

func (cs TextSerialiser) duration(d klog.Duration, withSign bool) string {
	if cs.DecimalDuration {
		return strconv.Itoa(d.InMinutes())
	}
	if withSign {
		return d.ToStringWithSign()
	}
	return d.ToString()
}

func (cs TextSerialiser) Date(d klog.Date) string {
	return cs.Styler.Props(tf.StyleProps{Color: tf.TEXT, IsUnderlined: true}).Format(d.ToString())
}

func (cs TextSerialiser) ShouldTotal(d klog.Duration) string {
	return cs.Styler.Props(tf.StyleProps{Color: tf.PURPLE}).Format(cs.duration(d, false))
}

func (cs TextSerialiser) Summary(s parser.SummaryText) string {
	txt := s.ToString()
	summaryStyler := cs.Styler.Props(tf.StyleProps{Color: tf.TEXT_SUBDUED})
	txt = klog.HashTagPattern.ReplaceAllStringFunc(txt, func(h string) string {
		return cs.Styler.Props(tf.StyleProps{Color: tf.TEXT_SUBDUED, IsBold: true}).FormatAndRestore(
			h, summaryStyler,
		)
	})
	return summaryStyler.Format(txt)
}

func (cs TextSerialiser) Range(r klog.Range) string {
	return cs.Styler.Props(tf.StyleProps{Color: tf.BLUE_DARK}).Format(r.ToString())
}

func (cs TextSerialiser) OpenRange(or klog.OpenRange) string {
	return cs.Styler.Props(tf.StyleProps{Color: tf.BLUE_LIGHT}).Format(or.ToString())
}

func (cs TextSerialiser) Duration(d klog.Duration) string {
	var c tf.Colour = tf.GREEN
	if strings.HasPrefix(d.ToStringWithSign(), "-") {
		c = tf.RED
	}
	return cs.Styler.Props(tf.StyleProps{Color: c}).Format(cs.duration(d, false))
}

func (cs TextSerialiser) SignedDuration(d klog.Duration) string {
	var c tf.Colour = tf.GREEN
	if strings.HasPrefix(d.ToStringWithSign(), "-") {
		c = tf.RED
	}
	return cs.Styler.Props(tf.StyleProps{Color: c}).Format(cs.duration(d, true))
}

func (cs TextSerialiser) Time(t klog.Time) string {
	return cs.Styler.Props(tf.StyleProps{Color: tf.BLUE_LIGHT}).Format(t.ToString())
}
