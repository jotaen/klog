package filter

import (
	"testing"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleRecordsForQuerying() []klog.Record {
	rs, _, err := parser.NewSerialParser().Parse(`
1999-12-29
No tags here

1999-12-30
Hello World #foo #first

1999-12-31
	5h #bar       [300]

2000-01-01
#foo #third
	1:30-1:45     [15]             
	6h #bar       [360]
	-30m          [-30]

2000-01-02
#foo #fourth
	7h #xyz       [420]

2000-01-03
#foo=a #fifth
	12:00-16:00   [240]
		#bar=1
	3h #bar=2     [180]
	12:00-?       [0]
`)
	if err != nil {
		panic(err)
	}
	return rs
}

type expect struct {
	date      klog.Date
	durations []int
}

func assertResult(t *testing.T, es []expect, rs []klog.Record) {
	require.Equal(t, len(es), len(rs), "unexpected number of records")
	for i, expct := range es {
		assert.Equal(t, expct.date, rs[i].Date(), "unexpected date")
		require.Equal(t, len(expct.durations), len(rs[i].Entries()), "unexpected number of entries")
		actualDurations := make([]int, len(rs[i].Entries()))
		for j, e := range rs[i].Entries() {
			actualDurations[j] = e.Duration().InMinutes()
		}
		assert.Equal(t, expct.durations, actualDurations, "unexpected duration")
	}
}

func TestQueryWithNoClauses(t *testing.T) {
	rs, hprws := Filter(And{}, sampleRecordsForQuerying())
	assert.False(t, hprws)
	assertResult(t, []expect{
		{klog.Ɀ_Date_(1999, 12, 29), []int{}},
		{klog.Ɀ_Date_(1999, 12, 30), []int{}},
		{klog.Ɀ_Date_(1999, 12, 31), []int{300}},
		{klog.Ɀ_Date_(2000, 1, 1), []int{15, 360, -30}},
		{klog.Ɀ_Date_(2000, 1, 2), []int{420}},
		{klog.Ɀ_Date_(2000, 1, 3), []int{240, 180, 0}},
	}, rs)
}

func TestQueryWithNoMatches(t *testing.T) {
	rs, hprws := Filter(IsInDateRange{
		From: klog.Ɀ_Date_(2002, 1, 1),
		To:   klog.Ɀ_Date_(2002, 1, 1),
	}, sampleRecordsForQuerying())
	assert.False(t, hprws)
	assertResult(t, []expect{}, rs)
}

func TestQueryAgainstEmptyInput(t *testing.T) {
	rs, hprws := Filter(IsInDateRange{
		From: klog.Ɀ_Date_(2002, 1, 1),
		To:   klog.Ɀ_Date_(2002, 1, 1),
	}, nil)
	assert.False(t, hprws)
	assertResult(t, []expect{}, rs)
}

func TestQueryWithAtDate(t *testing.T) {
	rs, hprws := Filter(IsInDateRange{
		From: klog.Ɀ_Date_(2000, 1, 2),
		To:   klog.Ɀ_Date_(2000, 1, 2),
	}, sampleRecordsForQuerying())
	assert.False(t, hprws)
	assertResult(t, []expect{
		{klog.Ɀ_Date_(2000, 1, 2), []int{420}},
	}, rs)
}

func TestQueryWithAfter(t *testing.T) {
	rs, hprws := Filter(IsInDateRange{
		From: klog.Ɀ_Date_(2000, 1, 1),
		To:   nil,
	}, sampleRecordsForQuerying())
	assert.False(t, hprws)
	assertResult(t, []expect{
		{klog.Ɀ_Date_(2000, 1, 1), []int{15, 360, -30}},
		{klog.Ɀ_Date_(2000, 1, 2), []int{420}},
		{klog.Ɀ_Date_(2000, 1, 3), []int{240, 180, 0}},
	}, rs)
}

func TestQueryWithBefore(t *testing.T) {
	rs, hprws := Filter(IsInDateRange{
		From: nil,
		To:   klog.Ɀ_Date_(2000, 1, 1),
	}, sampleRecordsForQuerying())
	assert.False(t, hprws)
	assertResult(t, []expect{
		{klog.Ɀ_Date_(1999, 12, 29), []int{}},
		{klog.Ɀ_Date_(1999, 12, 30), []int{}},
		{klog.Ɀ_Date_(1999, 12, 31), []int{300}},
		{klog.Ɀ_Date_(2000, 1, 1), []int{15, 360, -30}},
	}, rs)
}

func TestQueryWithTagOnOverallSummary(t *testing.T) {
	rs, hprws := Filter(HasTag{klog.NewTagOrPanic("foo", "")}, sampleRecordsForQuerying())
	assert.False(t, hprws)
	assertResult(t, []expect{
		{klog.Ɀ_Date_(1999, 12, 30), []int{}},
		{klog.Ɀ_Date_(2000, 1, 1), []int{15, 360, -30}},
		{klog.Ɀ_Date_(2000, 1, 2), []int{420}},
		{klog.Ɀ_Date_(2000, 1, 3), []int{240, 180, 0}},
	}, rs)
}

func TestQueryWithTagOnEntries(t *testing.T) {
	rs, hprws := Filter(HasTag{klog.NewTagOrPanic("bar", "")}, sampleRecordsForQuerying())
	assert.True(t, hprws)
	assertResult(t, []expect{
		{klog.Ɀ_Date_(1999, 12, 31), []int{300}},
		{klog.Ɀ_Date_(2000, 1, 1), []int{360}},
		{klog.Ɀ_Date_(2000, 1, 3), []int{240, 180}},
	}, rs)
}

func TestQueryWithTagOnEntriesAndInSummary(t *testing.T) {
	rs, hprws := Filter(And{[]Predicate{HasTag{klog.NewTagOrPanic("foo", "")}, HasTag{klog.NewTagOrPanic("bar", "")}}}, sampleRecordsForQuerying())
	assert.True(t, hprws)
	assertResult(t, []expect{
		{klog.Ɀ_Date_(2000, 1, 1), []int{360}},
		{klog.Ɀ_Date_(2000, 1, 3), []int{240, 180}},
	}, rs)
}

func TestQueryWithTagValues(t *testing.T) {
	rs, hprws := Filter(HasTag{klog.NewTagOrPanic("foo", "a")}, sampleRecordsForQuerying())
	assert.False(t, hprws)
	assertResult(t, []expect{
		{klog.Ɀ_Date_(2000, 1, 3), []int{240, 180, 0}},
	}, rs)
}

func TestQueryWithTagValuesInEntries(t *testing.T) {
	rs, hprws := Filter(HasTag{klog.NewTagOrPanic("bar", "1")}, sampleRecordsForQuerying())
	assert.True(t, hprws)
	assertResult(t, []expect{
		{klog.Ɀ_Date_(2000, 1, 3), []int{240}},
	}, rs)
}

func TestQueryWithTagNonMatchingValues(t *testing.T) {
	rs, hprws := Filter(HasTag{klog.NewTagOrPanic("bar", "3")}, sampleRecordsForQuerying())
	assert.False(t, hprws)
	assertResult(t, []expect{}, rs)
}

func TestQueryWithEntryTypes(t *testing.T) {
	{
		rs, hprws := Filter(IsEntryType{ENTRY_TYPE_DURATION}, sampleRecordsForQuerying())
		assert.True(t, hprws)
		assertResult(t, []expect{
			{klog.Ɀ_Date_(1999, 12, 31), []int{300}},
			{klog.Ɀ_Date_(2000, 1, 1), []int{360, -30}},
			{klog.Ɀ_Date_(2000, 1, 2), []int{420}},
			{klog.Ɀ_Date_(2000, 1, 3), []int{180}},
		}, rs)
	}
	{
		rs, hprws := Filter(IsEntryType{ENTRY_TYPE_DURATION_NEGATIVE}, sampleRecordsForQuerying())
		assert.True(t, hprws)
		assertResult(t, []expect{
			{klog.Ɀ_Date_(2000, 1, 1), []int{-30}},
		}, rs)
	}
	{
		rs, hprws := Filter(IsEntryType{ENTRY_TYPE_DURATION_POSITIVE}, sampleRecordsForQuerying())
		assert.True(t, hprws)
		assertResult(t, []expect{
			{klog.Ɀ_Date_(1999, 12, 31), []int{300}},
			{klog.Ɀ_Date_(2000, 1, 1), []int{360}},
			{klog.Ɀ_Date_(2000, 1, 2), []int{420}},
			{klog.Ɀ_Date_(2000, 1, 3), []int{180}},
		}, rs)
	}
	{
		rs, hprws := Filter(IsEntryType{ENTRY_TYPE_RANGE}, sampleRecordsForQuerying())
		assert.True(t, hprws)
		assertResult(t, []expect{
			{klog.Ɀ_Date_(2000, 1, 1), []int{15}},
			{klog.Ɀ_Date_(2000, 1, 3), []int{240}},
		}, rs)
	}
	{
		rs, hprws := Filter(IsEntryType{ENTRY_TYPE_OPEN_RANGE}, sampleRecordsForQuerying())
		assert.True(t, hprws)
		assertResult(t, []expect{
			{klog.Ɀ_Date_(2000, 1, 3), []int{0}},
		}, rs)
	}
}

func TestComplexFilterQueries(t *testing.T) {
	{
		rs, hprws := Filter(Or{[]Predicate{
			IsInDateRange{From: klog.Ɀ_Date_(2000, 1, 2), To: nil},
			HasTag{klog.NewTagOrPanic("first", "")},
			And{[]Predicate{
				Not{HasTag{klog.NewTagOrPanic("something", "1")}},
				IsEntryType{ENTRY_TYPE_RANGE},
			}},
		}}, sampleRecordsForQuerying())
		assert.True(t, hprws)
		assertResult(t, []expect{
			{klog.Ɀ_Date_(1999, 12, 30), []int{}},
			{klog.Ɀ_Date_(2000, 1, 1), []int{15}},
			{klog.Ɀ_Date_(2000, 1, 2), []int{420}},
			{klog.Ɀ_Date_(2000, 1, 3), []int{240, 180, 0}},
		}, rs)
	}
	{
		rs, hprws := Filter(And{[]Predicate{
			IsInDateRange{From: klog.Ɀ_Date_(2000, 1, 1), To: klog.Ɀ_Date_(2000, 1, 3)},
			HasTag{klog.NewTagOrPanic("bar", "")},
			Not{HasTag{klog.NewTagOrPanic("third", "")}},
			IsEntryType{ENTRY_TYPE_RANGE},
		}}, sampleRecordsForQuerying())
		assert.True(t, hprws)
		assertResult(t, []expect{
			{klog.Ɀ_Date_(2000, 1, 3), []int{240}},
		}, rs)
	}
	{
		rs, hprws := Filter(Not{Or{[]Predicate{
			IsInDateRange{From: klog.Ɀ_Date_(1999, 12, 30), To: klog.Ɀ_Date_(2000, 1, 1)},
			HasTag{klog.NewTagOrPanic("xyz", "")},
			And{[]Predicate{
				IsEntryType{ENTRY_TYPE_DURATION_POSITIVE},
				HasTag{klog.NewTagOrPanic("bar", "")},
			}},
		}}}, sampleRecordsForQuerying())
		assert.True(t, hprws)
		assertResult(t, []expect{
			{klog.Ɀ_Date_(1999, 12, 29), []int{}},
			{klog.Ɀ_Date_(2000, 1, 3), []int{240, 0}},
		}, rs)
	}
}
