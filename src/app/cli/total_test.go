package cli

import (
	"github.com/jotaen/klog/src/app/cli/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTotalOfEmptyInput(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(``)._Run((&Total{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "\nTotal: 0m\n(In 0 records)\n", state.printBuffer)
}

func TestTotalOfInput(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
2018-11-08
	1h

2018-11-09
	16:00-17:00

2150-11-10
Open ranges are not considered
	16:00 - ?
`)._Run((&Total{WarnArgs: lib.WarnArgs{NoWarn: true}}).Run)
	require.Nil(t, err)
	assert.Equal(t, "\nTotal: 2h\n(In 3 records)\n", state.printBuffer)
}

func TestTotalWithDiffing(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
2018-11-08 (8h!)
	8h30m

2018-11-09 (7h45m!)
	8:00 - 16:00
`)._Run((&Total{DiffArgs: lib.DiffArgs{Diff: true}}).Run)
	require.Nil(t, err)
	assert.Equal(t, "\nTotal: 16h30m\nShould: 15h45m!\nDiff: +45m\n(In 2 records)\n", state.printBuffer)
}
