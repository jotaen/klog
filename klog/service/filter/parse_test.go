package filter

import (
	"testing"

	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type et struct { // “Error Test”
	input string
	kind  error
	pos   int
	len   int
}

func checkError(t *testing.T, expected et, actual ParseError) {
	assert.ErrorIs(t, actual.Original(), expected.kind)
	pos, len := actual.Position()
	assert.Equal(t, expected.pos, pos)
	assert.Equal(t, expected.len, len)
}

func TestAtDate(t *testing.T) {
	p, err := Parse("2020-03-01")
	require.Nil(t, err)
	assert.Equal(t,
		IsInDateRange{klog.Ɀ_Date_(2020, 3, 1), klog.Ɀ_Date_(2020, 3, 1)},
		p)
}

func TestAndOperator(t *testing.T) {
	p, err := Parse("2020-01-01 && #hello")
	require.Nil(t, err)
	assert.Equal(t,
		And{[]Predicate{
			IsInDateRange{klog.Ɀ_Date_(2020, 1, 1), klog.Ɀ_Date_(2020, 1, 1)},
			HasTag{klog.NewTagOrPanic("hello", "")},
		}}, p)
}

func TestOrOperator(t *testing.T) {
	p, err := Parse("#foo || 1999-12-31")
	require.Nil(t, err)
	assert.Equal(t,
		Or{[]Predicate{
			HasTag{klog.NewTagOrPanic("foo", "")},
			IsInDateRange{klog.Ɀ_Date_(1999, 12, 31), klog.Ɀ_Date_(1999, 12, 31)},
		}}, p)
}

func TestCannotMixAndOrOnSameLevel(t *testing.T) {
	for _, tt := range []string{
		"#foo || 1999-12-31 && 2000-01-02",
		"#foo || 1999-12-31 || 2000-01-02 && 2020-07-24",
		"#foo && 1999-12-31 || 2000-01-02",
		"#foo && 1999-12-31 && 2000-01-02 || 2020-07-24",
		"#foo && (1999-12-31 || 2000-01-02 && 2021-05-17)",
		"#foo || (1999-12-31 && 2000-01-02 || 2021-05-17)",
	} {
		t.Run(tt, func(t *testing.T) {
			p, err := Parse(tt)
			require.ErrorIs(t, err.Original(), ErrCannotMixAndOr)
			require.Nil(t, p)
		})
	}
}

func TestNotOperator(t *testing.T) {
	p, err := Parse("!2020-01-01 && !#hello && !(2021-04-05 || #foo)")
	require.Nil(t, err)
	assert.Equal(t,
		And{[]Predicate{
			Not{
				IsInDateRange{klog.Ɀ_Date_(2020, 1, 1), klog.Ɀ_Date_(2020, 1, 1)},
			},
			Not{
				HasTag{klog.NewTagOrPanic("hello", "")},
			},
			Not{
				Or{[]Predicate{
					IsInDateRange{klog.Ɀ_Date_(2021, 4, 5), klog.Ɀ_Date_(2021, 4, 5)},
					HasTag{klog.NewTagOrPanic("foo", "")},
				}},
			},
		}}, p)
}

func TestGrouping(t *testing.T) {
	p, err := Parse("(#foo || #bar || #xyz) && 1999-12-31")
	require.Nil(t, err)
	assert.Equal(t,
		And{[]Predicate{
			Or{[]Predicate{
				HasTag{klog.NewTagOrPanic("foo", "")},
				HasTag{klog.NewTagOrPanic("bar", "")},
				HasTag{klog.NewTagOrPanic("xyz", "")},
			}},
			IsInDateRange{klog.Ɀ_Date_(1999, 12, 31), klog.Ɀ_Date_(1999, 12, 31)},
		}}, p)
}

func TestNestedGrouping(t *testing.T) {
	p, err := Parse("((#foo && (#bar || #xyz)) && 1999-12-31) || 1970-03-12")
	require.Nil(t, err)
	assert.Equal(t,
		Or{[]Predicate{
			And{[]Predicate{
				And{[]Predicate{
					HasTag{klog.NewTagOrPanic("foo", "")},
					Or{[]Predicate{
						HasTag{klog.NewTagOrPanic("bar", "")},
						HasTag{klog.NewTagOrPanic("xyz", "")},
					}},
				}},
				IsInDateRange{klog.Ɀ_Date_(1999, 12, 31), klog.Ɀ_Date_(1999, 12, 31)},
			}},
			IsInDateRange{klog.Ɀ_Date_(1970, 3, 12), klog.Ɀ_Date_(1970, 03, 12)},
		}}, p)
}

func TestClosedDateRange(t *testing.T) {
	for _, tt := range []struct {
		input string
		from  klog.Date
		to    klog.Date
	}{
		{"2020-03-06...2020-04-22", klog.Ɀ_Date_(2020, 3, 6), klog.Ɀ_Date_(2020, 4, 22)},
		{"2020-Q1...2020-Q2", klog.Ɀ_Date_(2020, 1, 1), klog.Ɀ_Date_(2020, 6, 30)},
		{"2020-04-17...2021", klog.Ɀ_Date_(2020, 4, 17), klog.Ɀ_Date_(2021, 12, 31)},
	} {
		p, err := Parse(tt.input)
		require.Nil(t, err)
		assert.Equal(t,
			IsInDateRange{tt.from, tt.to},
			p)
	}
}

func TestOpenDateRangeSince(t *testing.T) {
	for _, tt := range []struct {
		input string
		from  klog.Date
	}{
		{"2020-03-01...", klog.Ɀ_Date_(2020, 3, 1)},
		{"2020-W23...", klog.Ɀ_Date_(2020, 6, 1)},
	} {
		p, err := Parse(tt.input)
		require.Nil(t, err)
		assert.Equal(t,
			IsInDateRange{tt.from, nil},
			p)
	}
}

func TestOpenDateRangeUntil(t *testing.T) {
	for _, tt := range []struct {
		input string
		to    klog.Date
	}{
		{"...2020-03-01", klog.Ɀ_Date_(2020, 3, 1)},
		{"...2020-W23", klog.Ɀ_Date_(2020, 6, 7)},
	} {
		p, err := Parse(tt.input)
		require.Nil(t, err)
		assert.Equal(t,
			IsInDateRange{nil, tt.to},
			p)
	}
}

func TestPeriod(t *testing.T) {
	p, err := Parse("2020 || 2021-Q2 || 2022-08 || 2023-W46")
	require.Nil(t, err)
	assert.Equal(t,
		Or{[]Predicate{
			IsInDateRange{klog.Ɀ_Date_(2020, 1, 1), klog.Ɀ_Date_(2020, 12, 31)},
			IsInDateRange{klog.Ɀ_Date_(2021, 4, 1), klog.Ɀ_Date_(2021, 6, 30)},
			IsInDateRange{klog.Ɀ_Date_(2022, 8, 1), klog.Ɀ_Date_(2022, 8, 31)},
			IsInDateRange{klog.Ɀ_Date_(2023, 11, 13), klog.Ɀ_Date_(2023, 11, 19)},
		}}, p)
}

func TestTags(t *testing.T) {
	p, err := Parse("#tag || #tag-with=value || #tag-with='quoted value'")
	require.Nil(t, err)
	assert.Equal(t,
		Or{[]Predicate{
			HasTag{klog.NewTagOrPanic("tag", "")},
			HasTag{klog.NewTagOrPanic("tag-with", "value")},
			HasTag{klog.NewTagOrPanic("tag-with", "quoted value")},
		}}, p)
}

func TestEntryType(t *testing.T) {
	p, err := Parse("type:duration || type:range || type:open-range || type:duration-positive || type:duration-negative")
	require.Nil(t, err)
	assert.Equal(t,
		Or{[]Predicate{
			IsEntryType{ENTRY_TYPE_DURATION},
			IsEntryType{ENTRY_TYPE_RANGE},
			IsEntryType{ENTRY_TYPE_OPEN_RANGE},
			IsEntryType{ENTRY_TYPE_DURATION_POSITIVE},
			IsEntryType{ENTRY_TYPE_DURATION_NEGATIVE},
		}}, p)
}

func TestBracketMismatch(t *testing.T) {
	for _, tt := range []et{
		{"(2020-01", errUnbalancedBrackets, 0, 0},
		{"((2020-01", errUnbalancedBrackets, 0, 0},
		{"(2020-01-01 && (2020-02-02 || 2020-03-03", errUnbalancedBrackets, 0, 0},
		{"(2020-01-01))", errUnbalancedBrackets, 0, 13},
		{"2020-01-01)", errUnbalancedBrackets, 0, 11},
		{"(2020-01-01 && (2020-02-02))) || 2020-03-03", errUnbalancedBrackets, 0, 43},
	} {
		t.Run(tt.input, func(t *testing.T) {
			p, err := Parse(tt.input)
			require.Nil(t, p)
			checkError(t, tt, err)
		})
	}
}

func TestOperatorOperandSequence(t *testing.T) {
	for _, tt := range []string{
		// Operands: (date, date-range, period, tag)
		"2020-01-01 2020-02-02",
		"2020-01-01 (#foo && #bar)",
		"(#foo && #bar) 2020-01-01",
		"(#foo && #bar) #foo",
		"2020-01-01...2020-02-28 #foo",
		"2020-01-01... #foo",
		"...2020-01-01 #foo",
		"2020-01 2020-02",
		"2020-01-01 #foo",
		"2020-01 #foo",
		"#foo 2020-01-01",
		"#foo 2020-01",
		"#foo #foo",
		"type:duration #foo",
		"#foo type:duration",
		"2020 type:duration",
		"type:duration 2025-Q4",

		// And:
		"2020-01-01 && #tag #foo",
		"2020-01-01 && && 2020-02-02",
		"2020-01-01 && ( && #foo)",

		// Or:
		"2020-01-01 || #tag #foo",
		"2020-01-01 || || 2020-02-02",
		"2020-01-01 && ( || #foo)",

		// Not:
		"!&& #foo",
		"!|| #foo",
		"(!) #foo",
		"#foo !",
	} {
		t.Run(tt, func(t *testing.T) {
			p, err := Parse(tt)
			require.ErrorIs(t, err.Original(), errOperatorOperand)
			require.Nil(t, p)
		})
	}
}

func TestTokenizeError(t *testing.T) {
	p, err := Parse("2020-03-01(")
	require.Nil(t, p)
	require.Error(t, err)
}

func TestOperandOperatorSequenceExpectedError(t *testing.T) {
	for _, tt := range []et{
		{"2020-01-01 2020-02-02", ErrOperatorExpected, 11, 10},
		{"#foo type:open-range", ErrOperatorExpected, 5, 15},
		{"2020-01-01 && || #foo", ErrOperandExpected, 14, 2},
		{"2020-01-01 && (&& #foo)", ErrOperandExpected, 15, 2},
		{"2020-01-01 &&", ErrOperandExpected, 12, 1},
		{"2020-01-01 && ()", ErrOperandExpected, 15, 1},
	} {
		t.Run(tt.input, func(t *testing.T) {
			p, err := Parse(tt.input)
			require.Nil(t, p)
			checkError(t, tt, err)
		})
	}
}

func TestMalformedOperandError(t *testing.T) {
	for _, tt := range []et{
		{"2020-13-35", ErrIllegalTokenValue, 0, 10},
		{"2020-13-35...2020-Q7", ErrIllegalTokenValue, 0, 20},
		{"2020-01-02...2020-01-01", ErrIllegalTokenValue, 0, 23},
		{"2020-Q7", ErrIllegalTokenValue, 0, 7},
		{"type:foo", ErrIllegalTokenValue, 0, 8},
		{"foo", ErrUnrecognisedToken, 0, 1},
	} {
		t.Run(tt.input, func(t *testing.T) {
			p, err := Parse(tt.input)
			require.Nil(t, p)
			checkError(t, tt, err)
		})
	}
}
