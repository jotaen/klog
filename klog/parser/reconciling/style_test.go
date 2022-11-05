package reconciling

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser"
	"github.com/jotaen/klog/klog/parser/txt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDefaultStyle(t *testing.T) {
	assert.Equal(t, &style{
		lineEnding:                          styleProp[string]{"\n", false},
		indentation:                         styleProp[string]{"    ", false},
		dateUseDashes:                       styleProp[bool]{true, false},
		timeUse24HourClock:                  styleProp[bool]{true, false},
		rangesUseSpacesAroundDash:           styleProp[bool]{true, false},
		openRangeAdditionalPlaceholderChars: styleProp[int]{0, false},
	}, defaultStyle())
}

func TestDetectsStyleFromMinimalFile(t *testing.T) {
	rs, bs := parseOrPanic("2000-01-01")
	s := determine(rs[0], bs[0])
	assert.Equal(t, &style{
		lineEnding:                          styleProp[string]{"\n", false},
		indentation:                         styleProp[string]{"    ", false},
		dateUseDashes:                       styleProp[bool]{true, true},
		timeUse24HourClock:                  styleProp[bool]{true, false},
		rangesUseSpacesAroundDash:           styleProp[bool]{true, false},
		openRangeAdditionalPlaceholderChars: styleProp[int]{0, false},
	}, s)
}

func TestDetectCanonicalStyle(t *testing.T) {
	rs, bs := parseOrPanic("2000-01-01\nTest\n    8:00 - ?\n")
	s := determine(rs[0], bs[0])
	assert.Equal(t, &style{
		lineEnding:                          styleProp[string]{"\n", true},
		indentation:                         styleProp[string]{"    ", true},
		dateUseDashes:                       styleProp[bool]{true, true},
		timeUse24HourClock:                  styleProp[bool]{true, true},
		rangesUseSpacesAroundDash:           styleProp[bool]{true, true},
		openRangeAdditionalPlaceholderChars: styleProp[int]{0, true},
	}, s)
}

func TestDetectsCustomStyle(t *testing.T) {
	rs, bs := parseOrPanic("2000/01/01\r\nTest\r\n\t8:00am-?????\r\n")
	s := determine(rs[0], bs[0])
	assert.Equal(t, &style{
		lineEnding:                          styleProp[string]{"\r\n", true},
		indentation:                         styleProp[string]{"\t", true},
		dateUseDashes:                       styleProp[bool]{false, true},
		timeUse24HourClock:                  styleProp[bool]{false, true},
		rangesUseSpacesAroundDash:           styleProp[bool]{false, true},
		openRangeAdditionalPlaceholderChars: styleProp[int]{4, true},
	}, s)
}

func TestElectStyle(t *testing.T) {
	rs, bs := parseOrPanic(
		"2001-05-19\n\t1:00 - 2:00\n\n",
		"2001/05/19\r\n  1:00am-2:00pm\r\n\r\n",
		"2001-05-19\n   1:00am-2:00pm\n   2:00pm-3:00pm\n\n",
		"2001/05/19\r\n  1:00 - 2:00\r\n\r\n",
		"2001-05-19\r\n    1:00am-???\r\n\r\n",
	)
	result := elect(*defaultStyle(), rs, bs)
	assert.Equal(t, &style{
		lineEnding:                          styleProp[string]{"\r\n", true},
		indentation:                         styleProp[string]{"  ", true},
		dateUseDashes:                       styleProp[bool]{true, true},
		timeUse24HourClock:                  styleProp[bool]{false, true},
		rangesUseSpacesAroundDash:           styleProp[bool]{false, true},
		openRangeAdditionalPlaceholderChars: styleProp[int]{2, true},
	}, result)
}

func TestElectStyleDoesNotOverrideSetPreferences(t *testing.T) {
	majorityRs, majorityBs := parseOrPanic(
		"2001-05-19\n\t1:00 - 2:00\n\n",
		"2001/05/19\r\n  1:00am-2:00pm\r\n\r\n",
		"2001-05-19\n   1:00am-2:00pm\n   2:00pm-3:00pm\n\n",
		"2001/05/19\r\n  1:00 - 2:00\r\n\r\n",
		"2001-05-19\r\n    1:00am-2:00pm\r\n\r\n",
	)
	winnerRs, winnerBs := parseOrPanic("2018/01/01\n\t8:00 - 9:00")
	result := elect(*determine(winnerRs[0], winnerBs[0]), majorityRs, majorityBs)
	assert.Equal(t, &style{
		lineEnding:                          styleProp[string]{"\n", true},
		indentation:                         styleProp[string]{"\t", true},
		dateUseDashes:                       styleProp[bool]{false, true},
		timeUse24HourClock:                  styleProp[bool]{true, true},
		rangesUseSpacesAroundDash:           styleProp[bool]{true, true},
		openRangeAdditionalPlaceholderChars: styleProp[int]{0, true},
	}, result)
}

func parseOrPanic(recordsAsText ...string) ([]klog.Record, []txt.Block) {
	rs, bs, err := parser.NewSerialParser().Parse(strings.Join(recordsAsText, ""))
	if err != nil {
		panic("Invalid data")
	}
	return rs, bs
}
