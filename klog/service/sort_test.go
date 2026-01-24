package service

import (
	"testing"

	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
)

func sampleRecordsForQuerying() []klog.Record {
	return []klog.Record{
		func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(1999, 12, 30))
			r.SetSummary(klog.Ɀ_RecordSummary_("Hello World", "#foo"))
			return r
		}(), func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(1999, 12, 31))
			r.AddDuration(klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("#bar"))
			return r
		}(), func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(2000, 1, 1))
			r.SetSummary(klog.Ɀ_RecordSummary_("#foo"))
			r.AddDuration(klog.NewDuration(0, 15), nil)
			r.AddDuration(klog.NewDuration(6, 0), klog.Ɀ_EntrySummary_("#bar"))
			r.AddDuration(klog.NewDuration(0, -30), nil)
			return r
		}(), func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(2000, 1, 2))
			r.SetSummary(klog.Ɀ_RecordSummary_("#foo"))
			r.AddDuration(klog.NewDuration(7, 0), nil)
			return r
		}(), func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(2000, 1, 3))
			r.SetSummary(klog.Ɀ_RecordSummary_("#foo=a"))
			r.AddDuration(klog.NewDuration(4, 0), klog.Ɀ_EntrySummary_("test", "foo #bar=1"))
			r.AddDuration(klog.NewDuration(4, 0), klog.Ɀ_EntrySummary_("#bar=2"))
			r.Start(klog.NewOpenRange(klog.Ɀ_Time_(12, 00)), nil)
			return r
		}(),
	}
}

func TestQueryWithSorting(t *testing.T) {
	ss := sampleRecordsForQuerying()
	for _, x := range []struct{ rs []klog.Record }{
		{ss},
		{[]klog.Record{ss[3], ss[1], ss[2], ss[0], ss[4]}},
		{[]klog.Record{ss[1], ss[4], ss[0], ss[3], ss[2]}},
	} {
		ascending := Sort(x.rs, true)
		assert.Equal(t, []klog.Record{ss[0], ss[1], ss[2], ss[3], ss[4]}, ascending)

		descending := Sort(x.rs, false)
		assert.Equal(t, []klog.Record{ss[4], ss[3], ss[2], ss[1], ss[0]}, descending)
	}
}
