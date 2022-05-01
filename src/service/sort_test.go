package service

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortRecordsByDate(t *testing.T) {
	ss := sampleRecordsForQuerying()
	for _, x := range []struct{ rs []Record }{
		{ss},
		{[]Record{ss[3], ss[1], ss[2], ss[0], ss[4]}},
		{[]Record{ss[1], ss[4], ss[0], ss[3], ss[2]}},
	} {
		ascending := Sort(x.rs, true)
		assert.Equal(t, []Record{ss[0], ss[1], ss[2], ss[3], ss[4]}, ascending)

		descending := Sort(x.rs, false)
		assert.Equal(t, []Record{ss[4], ss[3], ss[2], ss[1], ss[0]}, descending)
	}
}
