package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"klog/app/cli/lib"
	"testing"
)

func TestSkipsWhenThereAreNoRecords(t *testing.T) {
	out, err := NewTestingContext()._SetRecords(``)._Run((&Now{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "\nNo record found for today\n", out)
}

func TestSkipsWhenThereAreNoRecentRecords(t *testing.T) {
	out, err := NewTestingContext()._SetNow(1999, 3, 14, 0, 0)._SetRecords(`
1999-03-12
	4h
`)._Run((&Now{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "\nNo record found for today\n", out)
}

func TestPrintsTodaysEvalutaion(t *testing.T) {
	out, err := NewTestingContext()._SetNow(1999, 3, 14, 15, 0)._SetRecords(`
1999-03-13
	12h5m

1999-03-14
	1h

1999-03-14
	3h
	13:15 - ?
`)._Run((&Now{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
            Today    Overall
Total       5h45m     17h50m
`, out)
}

func TestPrintsEvalutaionWithDiff(t *testing.T) {
	out, err := NewTestingContext()._SetNow(1999, 3, 14, 3, 13)._SetRecords(`
1999-03-12 (3h10m!)
	2h50m

1999-03-13 (6h!)
	23:38 - ?
`)._Run((&Now{DiffArgs: lib.DiffArgs{Diff: true}}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
        Yesterday    Overall
Total       3h35m      6h25m
Should        6h!     9h10m!
Diff       -2h25m     -2h45m
E.T.A.       5:38       5:58
`, out)
}
