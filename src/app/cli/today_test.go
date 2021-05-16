package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"klog/app/cli/lib"
	"testing"
)

func TestSkipsWhenThereAreNoRecords(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(``)._Run((&Today{}).Run)
	require.EqualError(t, err, "No current record (dated either today or yesterday)")
	assert.Equal(t, "", state.printBuffer)
}

func TestSkipsWhenThereAreNoRecentRecords(t *testing.T) {
	state, err := NewTestingContext()._SetNow(1999, 3, 14, 0, 0)._SetRecords(`
1999-03-12
	4h
`)._Run((&Today{}).Run)
	require.EqualError(t, err, "No current record (dated either today or yesterday)")
	assert.Equal(t, "", state.printBuffer)
}

func TestPrintsTodaysEvalutaion(t *testing.T) {
	state, err := NewTestingContext()._SetNow(1999, 3, 14, 15, 0)._SetRecords(`
1999-03-12
	5m

1999-03-13
	12h

1999-03-14
	1h

1999-03-14
	3h
	13:15 - ?
`)._Run((&Today{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
             Total
Today        5h45m
Other        12h5m
          ========
All        +17h50m
`, state.printBuffer)
}

func TestFallsBackToYesterday(t *testing.T) {
	state, err := NewTestingContext()._SetNow(1999, 3, 14, 15, 0)._SetRecords(`
1999-03-12
	5m

1999-03-13
	12h
`)._Run((&Today{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
             Total
Yesterday      12h
Other           5m
          ========
All         +12h5m
`, state.printBuffer)
}

func TestPrintsEvaluationWithDiff(t *testing.T) {
	state, err := NewTestingContext()._SetNow(1999, 3, 14, 18, 13)._SetRecords(`
1999-03-12 (3h10m!)
	6h50m

1999-03-14 (6h!)
	14:38 - ?
`)._Run((&Today{DiffArgs: lib.DiffArgs{Diff: true}}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
             Total    Should     Diff   End-Time
Today        3h35m       6h!   -2h25m      20:38
Other        6h50m    3h10m!   +3h40m
          ===========================
All        +10h25m    9h10m!   +1h15m      16:58
`, state.printBuffer)
}
