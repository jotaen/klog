package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"klog/app/cli/lib"
	"testing"
)

func TestReportOfEmptyInput(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(``)._Run((&Report{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "", state.printBuffer)
}

func TestReportOfRecords(t *testing.T) {
	/*
		Aspects tested:
		- Multiple records per date unified into one item
		- Sorting by date
		- Not repeating year or month label
		- Weekdays
	*/
	state, err := NewTestingContext()._SetRecords(`
2021-01-17
	332h

2021-01-17
	1h

2019-12-01

2021-03-03
	<23:00 - 0:00

2020-12-30
	1h
    8:00am - 04:47pm

2021-03-02
    -8h2m

2021-01-19
	5m
`)._SetNow(2021, 3, 4, 0, 0)._Run((&Report{WarnArgs: lib.WarnArgs{NoWarn: true}}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
                       Total
2019 Dec    Sun  1.       0m
2020 Dec    Wed 30.    9h47m
2021 Jan    Sun 17.     333h
            Tue 19.       5m
     Mar    Tue  2.    -8h2m
            Wed  3.       1h
                    ========
                    +335h50m
`, state.printBuffer)
}

func TestReportConsecutive(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
2020-09-29
	1h

2020-10-04
	3h

2020-10-02
`)._Run((&Report{Fill: true}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
                       Total
2020 Sep    Tue 29.       1h
            Wed 30.  
     Oct    Thu  1.  
            Fri  2.       0m
            Sat  3.  
            Sun  4.       3h
                    ========
                         +4h
`, state.printBuffer)
}

func TestReportWithDiff(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
2018-07-07 (8h!)
	8h

2018-07-08 (5h30m!)
	2h

2018-07-09 (2h!)
	5h20m

2018-07-09 (19m!)
`)._Run((&Report{DiffArgs: lib.DiffArgs{Diff: true}}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
                       Total    Should     Diff
2018 Jul    Sat  7.       8h       8h!       0m
            Sun  8.       2h    5h30m!   -3h30m
            Mon  9.    5h20m    2h19m!    +3h1m
                    ===========================
                     +15h20m   15h49m!     -29m
`, state.printBuffer)
}
