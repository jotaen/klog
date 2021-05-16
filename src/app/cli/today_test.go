package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"klog/app/cli/lib"
	"testing"
)

func TestPrintsTodaysEvalutaion(t *testing.T) {
	state, err := NewTestingContext()._SetNow(1999, 3, 14, 19, 9)._SetRecords(`
1999-03-12
	5m

1999-03-13
	12h

1999-03-14
	1h

1999-03-14
	3h
	13:15 - 15:00
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
	state, err := NewTestingContext()._SetNow(1999, 3, 14, 19, 12)._SetRecords(`
1999-03-12 (3h10m!)
	6h50m

1999-03-14 (6h!)
	14:38 - 18:13
`)._Run((&Today{DiffArgs: lib.DiffArgs{Diff: true}}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
             Total    Should     Diff
Today        3h35m       6h!   -2h25m
Other        6h50m    3h10m!   +3h40m
          ===========================
All        +10h25m    9h10m!   +1h15m
`, state.printBuffer)
}

func TestPrintsEvaluationWithNow(t *testing.T) {
	state, err := NewTestingContext()._SetNow(1999, 3, 14, 18, 13)._SetRecords(`
1999-03-12
	6h50m

1999-03-14
	14:38 - ??
`)._Run((&Today{NowArgs: lib.NowArgs{Now: true}}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
             Total
Today        3h35m
Other        6h50m
          ========
All        +10h25m
`, state.printBuffer)
}

func TestPrintsEvaluationWithDiffAndNow(t *testing.T) {
	state, err := NewTestingContext()._SetNow(1999, 3, 14, 18, 13)._SetRecords(`
1999-03-12 (3h10m!)
	6h50m

1999-03-14 (6h!)
	14:38 - ?
`)._Run((&Today{DiffArgs: lib.DiffArgs{Diff: true}, NowArgs: lib.NowArgs{Now: true}}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
             Total    Should     Diff   End-Time
Today        3h35m       6h!   -2h25m      20:38
Other        6h50m    3h10m!   +3h40m
          ===========================
All        +10h25m    9h10m!   +1h15m      16:58
`, state.printBuffer)
}

func TestPrintsNAWhenNoCurrentRecord(t *testing.T) {
	state, err := NewTestingContext()._SetNow(1999, 3, 16, 18, 13)._SetRecords(`
1999-03-12 (3h10m!)
	6h50m
`)._Run((&Today{DiffArgs: lib.DiffArgs{Diff: true}, NowArgs: lib.NowArgs{Now: true}}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
             Total    Should     Diff   End-Time
Today          n/a       n/a      n/a        n/a
Other        6h50m    3h10m!   +3h40m
          ===========================
All         +6h50m    3h10m!   +3h40m        n/a
`, state.printBuffer)
}
