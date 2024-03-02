package cli

import (
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrintOutEmptyInput(t *testing.T) {
	{
		state, err := NewTestingContext()._SetRecords(``)._Run((&Print{}).Run)
		require.Nil(t, err)
		assert.Equal(t, "", state.printBuffer)
	}
	{
		state, err := NewTestingContext()._SetRecords(``)._Run((&Print{
			WithTotals: true,
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, "", state.printBuffer)
	}
}

func TestPrintOutRecord(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
2018-01-31
Hello #world
    1h
`)._Run((&Print{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-31
Hello #world
    1h

`, state.printBuffer)
}

func TestPrintOutRecordInCanonicalFormat(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
2018-01-31
Hello #world
  09:00-13:00
  22:00  -  24:00
  60m
  2h0m
  0h
`)._Run((&Print{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-31
Hello #world
    9:00-13:00
    22:00 - 0:00>
    1h
    2h
    0m

`, state.printBuffer)
}

func TestPrintOutRecordsInChronologicalOrder(t *testing.T) {
	original := "2018-02-01\n\n2018-01-30\n\n2018-01-31"

	stateUnsorted, _ := NewTestingContext().
		_SetRecords(original).
		_Run((&Print{}).Run)
	assert.Equal(t, "\n"+original+"\n\n", stateUnsorted.printBuffer)

	stateSortedAsc, _ := NewTestingContext().
		_SetRecords(original).
		_Run((&Print{SortArgs: util.SortArgs{Sort: "asc"}}).Run)
	assert.Equal(t, "\n2018-01-30\n\n2018-01-31\n\n2018-02-01\n\n", stateSortedAsc.printBuffer)

	stateSortedDesc, _ := NewTestingContext().
		_SetRecords(original).
		_Run((&Print{SortArgs: util.SortArgs{Sort: "desc"}}).Run)
	assert.Equal(t, "\n2018-02-01\n\n2018-01-31\n\n2018-01-30\n\n", stateSortedDesc.printBuffer)
}

func TestPrintRecordsWithDurations(t *testing.T) {
	state, err := NewTestingContext()._SetNow(2018, 02, 07, 19, 00)._SetRecords(`
2018-01-31
Hello #world
Test test test
	1h

2018-02-04
	10:00 - 17:22
	+6h
	-5m

2018-02-07
	35m
		Foo
		Bar
	18:00 - ? I just
		started something
`)._Run((&Print{
		WithTotals: true,
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
     1h  |  2018-01-31
         |  Hello #world
         |  Test test test
     1h  |      1h

 13h17m  |  2018-02-04
  7h22m  |      10:00 - 17:22
     6h  |      +6h
    -5m  |      -5m

    35m  |  2018-02-07
    35m  |      35m
         |          Foo
         |          Bar
     0m  |      18:00 - ? I just
         |          started something

`, state.printBuffer)
}
