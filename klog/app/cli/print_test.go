package cli

import (
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrintOutEmptyInput(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(``)._Run((&Print{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "", state.printBuffer)
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
    9:00 - 13:00
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
		_Run((&Print{SortArgs: lib.SortArgs{Sort: "asc"}}).Run)
	assert.Equal(t, "\n2018-01-30\n\n2018-01-31\n\n2018-02-01\n\n", stateSortedAsc.printBuffer)

	stateSortedDesc, _ := NewTestingContext().
		_SetRecords(original).
		_Run((&Print{SortArgs: lib.SortArgs{Sort: "desc"}}).Run)
	assert.Equal(t, "\n2018-02-01\n\n2018-01-31\n\n2018-01-30\n\n", stateSortedDesc.printBuffer)
}
