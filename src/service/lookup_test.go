package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "klog"
	"testing"
)

func sampleRecords() []Record {
	return []Record{
		func() Record {
			r := NewRecord(Ɀ_Date_(1999, 12, 30))
			_ = r.SetSummary("#foo")
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(1999, 12, 31))
			r.AddDuration(NewDuration(5, 0), "#bar")
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 1))
			_ = r.SetSummary("#foo")
			r.AddDuration(NewDuration(0, 15), "")
			r.AddDuration(NewDuration(6, 0), "#bar")
			r.AddDuration(NewDuration(0, -30), "")
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 2))
			_ = r.SetSummary("#foo")
			r.AddDuration(NewDuration(7, 0), "")
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 3))
			_ = r.SetSummary("#foo")
			r.AddDuration(NewDuration(4, 0), "#bar")
			r.AddDuration(NewDuration(4, 0), "#bar")
			return r
		}(),
	}
}

func TestFindFilterWithNoClauses(t *testing.T) {
	rs := FindFilter(sampleRecords(), Filter{})
	require.Len(t, rs, 5)
	assert.Equal(t, NewDuration(5+6+7+8, -30+15), Total(rs...))
}

func TestFindFilterWithAfter(t *testing.T) {
	rs := FindFilter(sampleRecords(), Filter{AfterEq: Ɀ_Date_(2000, 1, 1)})
	require.Len(t, rs, 3)
	assert.Equal(t, 1, rs[0].Date().Day())
	assert.Equal(t, 2, rs[1].Date().Day())
	assert.Equal(t, 3, rs[2].Date().Day())
}

func TestFindFilterWithBefore(t *testing.T) {
	rs := FindFilter(sampleRecords(), Filter{BeforeEq: Ɀ_Date_(2000, 1, 1)})
	require.Len(t, rs, 3)
	assert.Equal(t, 30, rs[0].Date().Day())
	assert.Equal(t, 31, rs[1].Date().Day())
	assert.Equal(t, 1, rs[2].Date().Day())
}

func TestFindFilterWithTagOnEntries(t *testing.T) {
	rs := FindFilter(sampleRecords(), Filter{Tags: []string{"bar"}})
	require.Len(t, rs, 3)
	assert.Equal(t, 31, rs[0].Date().Day())
	assert.Equal(t, 1, rs[1].Date().Day())
	assert.Equal(t, 3, rs[2].Date().Day())
	assert.Equal(t, NewDuration(5+8+6, 0), Total(rs...))
}

func TestFindFilterWithTagOnOverallSummary(t *testing.T) {
	rs := FindFilter(sampleRecords(), Filter{Tags: []string{"foo"}})
	require.Len(t, rs, 4)
	assert.Equal(t, 30, rs[0].Date().Day())
	assert.Equal(t, 1, rs[1].Date().Day())
	assert.Equal(t, 2, rs[2].Date().Day())
	assert.Equal(t, 3, rs[3].Date().Day())
	assert.Equal(t, NewDuration(6+7+8, -30+15), Total(rs...))
}

func TestFindFilterWithTagOnEntriesAndInSummary(t *testing.T) {
	rs := FindFilter(sampleRecords(), Filter{Tags: []string{"foo", "bar"}})
	require.Len(t, rs, 2)
	assert.Equal(t, 1, rs[0].Date().Day())
	assert.Equal(t, 3, rs[1].Date().Day())
	assert.Equal(t, NewDuration(8+6, 0), Total(rs...))
}
