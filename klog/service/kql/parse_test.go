package kql

import (
	"testing"

	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			require.ErrorIs(t, err, ErrCannotMixAndOr)
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
	p, err := Parse("2020-03-06...2020-04-22")
	require.Nil(t, err)
	assert.Equal(t,
		IsInDateRange{klog.Ɀ_Date_(2020, 3, 6), klog.Ɀ_Date_(2020, 4, 22)},
		p)
}

func TestOpenDateRangeSince(t *testing.T) {
	p, err := Parse("2020-03-01...")
	require.Nil(t, err)
	assert.Equal(t,
		IsInDateRange{klog.Ɀ_Date_(2020, 3, 1), nil},
		p)
}

func TestOpenDateRangeUntil(t *testing.T) {
	p, err := Parse("...2020-03-01")
	require.Nil(t, err)
	assert.Equal(t,
		IsInDateRange{nil, klog.Ɀ_Date_(2020, 3, 1)},
		p)
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

func TestBracketMismatch(t *testing.T) {
	for _, tt := range []string{
		"(2020-01",
		"((2020-01",
		"(2020-01-01))",
		"2020-01-01)",
		"(2020-01-01 && (2020-02-02 || 2020-03-03",
		"(2020-01-01 && (2020-02-02))) || 2020-03-03",
	} {
		t.Run(tt, func(t *testing.T) {
			p, err := Parse(tt)
			require.ErrorIs(t, err, ErrUnbalancedBrackets)
			require.Nil(t, p)
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
			require.ErrorIs(t, err, errOperatorOperand)
			require.Nil(t, p)
		})
	}
}
