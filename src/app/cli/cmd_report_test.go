package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReportOfEmptyInput(t *testing.T) {
	out, err := RunWithContext(``, (&Report{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "", out)
}

func TestReportOfRecords(t *testing.T) {
	/*
		Aspects tested:
		- Multiple records per date unified into one item
		- Sorting by date
		- Not repeating year or month label
		- Weekdays
	*/
	out, err := RunWithContext(`
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
`, (&Report{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
                      Total
2019 Dec    Su  1.       0m
2020 Dec    We 30.    9h47m
2021 Jan    Su 17.     333h
            Tu 19.       5m
     Mar    Tu  2.    -8h2m
            We  3.       1h
                   ========
                   +335h50m
`, out)
}

func TestReportConsecutive(t *testing.T) {
	out, err := RunWithContext(`
2020-09-29
	1h

2020-10-04
	3h

2020-10-02
`, (&Report{Fill: true}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
                      Total
2020 Sep    Tu 29.       1h
            We 30.  
     Oct    Th  1.  
            Fr  2.       0m
            Sa  3.  
            Su  4.       3h
                   ========
                        +4h
`, out)
}

func TestReportWithDiff(t *testing.T) {
	out, err := RunWithContext(`
2018-07-07 (8h!)
	8h

2018-07-08 (5h30m!)
	2h

2018-07-09 (2h!)
	5h20m

2018-07-09 (19m!)
`, (&Report{DiffArg: DiffArg{Diff: true}}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
                      Total    Should     Diff
2018 Jul    Sa  7.       8h       8h!       0m
            Su  8.       2h    5h30m!   -3h30m
            Mo  9.    5h20m    2h19m!    +3h1m
                   ===========================
                    +15h20m   15h49m!     -29m
`, out)
}
