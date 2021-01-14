package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"klog/datetime"
	datetime2 "klog/testutil/datetime"
	"testing"
)

func TestMinimalValidDocument(t *testing.T) {
	text := `
date: 2020-01-01
`
	w, errs := Parse(text)

	require.Nil(t, errs)
	require.NotNil(t, w)

	assert.Equal(t, "", w.Summary())
	assert.Nil(t, w.Durations())
	assert.Nil(t, w.Ranges())
	assert.Nil(t, w.OpenRange())
}

func TestParsingAllFieldsCorrectly(t *testing.T) {
	yaml := `
date: 2008-12-03
summary: Just a normal day
hours:
- start: 12:00
  end: 12:50
- start: 12:00
  end: 12:00
- start: 23:55 yesterday
  end: 09:05
- start: 19:12
  end: 1:59 tomorrow
- start: 10:15
- time: 2h
- time: 05h 03m
- time: -1h 12m
`
	w, errs := Parse(yaml)

	require.Equal(t, 0, len(errs))
	require.NotNil(t, w)

	assert.Equal(t, "Just a normal day", w.Summary())
	assert.Equal(t, []datetime.TimeRange{
		datetime2.Range_(datetime2.Time_(12, 00), datetime2.Time_(12, 50)),
		datetime2.Range_(datetime2.Time_(12, 00), datetime2.Time_(12, 00)),
		datetime2.Range_(datetime2.TimeYesterday_(23, 55), datetime2.Time_(9, 5)),
		datetime2.Range_(datetime2.Time_(19, 12), datetime2.TimeTomorrow_(1, 59)),
	}, w.Ranges())
	assert.Equal(t, datetime2.Time_(10, 15), w.OpenRange())
	assert.Equal(t, []datetime.Duration{
		datetime.NewDuration(2, 0),
		datetime.NewDuration(5, 3),
		datetime.NewDuration(-1, -12),
	}, w.Durations())
}

func TestMalformedYamlFails(t *testing.T) {
	yaml := `
date: 2005-05-01
foo
bar
`
	w, errs := Parse(yaml)
	assert.Nil(t, w)
	assert.Error(t, errs[0])
	assert.Contains(t, errs, parserError("MALFORMED_YAML", ""))
}

func TestAbsentRequiredPropertiesFails(t *testing.T) {
	yaml := `
summary: Just a normal day
`
	w, errs := Parse(yaml)
	assert.Nil(t, w)
	assert.Contains(t, errs, parserError("MISSING_REQUIRED_PROPERTY", "date"))
}

func TestTimeAndRangeCannotAppearTogether(t *testing.T) {
	yaml := `
date: 1999-12-31
hours:
- end: 8:00
- start: 10:00
  end: 11:00
  time: 1:00
- start: 9:00
  time: 10:00
- end: 9:00
  time: 10:00
`
	_, errs := Parse(yaml)
	assert.Equal(t, []ParserError{
		parserError("MALFORMED_HOURS", "hours"),
		parserError("MALFORMED_HOURS", "hours"),
		parserError("MALFORMED_HOURS", "hours"),
		parserError("MALFORMED_HOURS", "hours"),
	}, errs)
}

func TestMalformedValuesFails(t *testing.T) {
	yaml := `
date: 1999-12-31
hours:
- start: asdf
  end: 9:00
- start: 8:00
  end: asdf
- start: asdf
  end: asdf
- start: asdf
- time: asdf
`
	w, errs := Parse(yaml)
	assert.Equal(t, w, nil)
	assert.Equal(t, []ParserError{
		parserError("MALFORMED_TIME", "start: asdf"),
		parserError("MALFORMED_TIME", "end: asdf"),
		parserError("MALFORMED_TIME", "start: asdf"),
		parserError("MALFORMED_TIME", "start: asdf"),
		parserError("MALFORMED_DURATION", "time: asdf"),
	}, errs)
}
