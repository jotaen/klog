package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReportOfEmptyInput(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(``)._Run((&Report{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "", state.printBuffer)
}

func TestReportOfEmptyFilteredData(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
2022-10-30
	8h
`)._Run((&Report{
		FilterArgs: lib.FilterArgs{Date: klog.Ɀ_Date_(2022, 10, 31)},
		Fill:       true,
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, "", state.printBuffer)
}

func TestDayReportOfRecords(t *testing.T) {
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
                     335h50m
`, state.printBuffer)
}

func TestDayReportConsecutive(t *testing.T) {
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
                          4h
`, state.printBuffer)
}

func TestDayReportWithDiff(t *testing.T) {
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
                    ======== ========= ========
                      15h20m   15h49m!     -29m
`, state.printBuffer)
}

func TestDayReportWithDecimal(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
2018-07-07 (8h!)
	8h

2018-07-08 (5h30m!)
	2h

2018-07-09 (2h!)
	5h20m

2018-07-09 (19m!)
`)._Run((&Report{DiffArgs: lib.DiffArgs{Diff: true}, DecimalArgs: lib.DecimalArgs{Decimal: true}}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
                       Total    Should     Diff
2018 Jul    Sat  7.      480       480        0
            Sun  8.      120       330     -210
            Mon  9.      320       139      181
                    ======== ========= ========
                         920       949      -29
`, state.printBuffer)
}

func TestWeekReport(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
2016-01-03
  1h

2016-01-04
  1h

2016-12-31
  1h

2017-01-01
  1h
`)._Run((&Report{AggregateBy: "week"}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
                 Total
2015  Week 53       1h
2016  Week  1       1h
      Week 52       2h
              ========
                    4h
`, state.printBuffer)
}

func TestWeekReportWithFillAndDiff(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
2018-12-09 (8h!)
	8h

2018-12-26 (1h30m!)
	2h

2018-12-31 (30m!)
	15m

2019-01-02 (2h!)
	3h

2019-01-08 (19m!)
`)._Run((&Report{AggregateBy: "week", DiffArgs: lib.DiffArgs{Diff: true}, Fill: true}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
                 Total    Should     Diff
2018  Week 49       8h       8h!       0m
      Week 50                            
      Week 51                            
      Week 52       2h    1h30m!     +30m
2019  Week  1    3h15m    2h30m!     +45m
      Week  2       0m      19m!     -19m
              ======== ========= ========
                13h15m   12h19m!     +56m
`, state.printBuffer)
}

func TestQuarterReport(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
2018-02-02 (8h!)
	8h

2018-04-10 (5h30m!)
	2h

2018-05-23 (2h!)
	5h20m

2019-01-01 (19m!)
`)._Run((&Report{AggregateBy: "quarter", DiffArgs: lib.DiffArgs{Diff: true}, Fill: true}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
           Total    Should     Diff
2018 Q1       8h       8h!       0m
     Q2    7h20m    7h30m!     -10m
     Q3                            
     Q4                            
2019 Q1       0m      19m!     -19m
        ======== ========= ========
          15h20m   15h49m!     -29m
`, state.printBuffer)
}

func TestMonthReport(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
2018-02-02 (8h!)
	8h

2018-04-10 (5h30m!)
	2h

2018-05-23 (2h!)
	5h20m

2019-01-01 (19m!)
`)._Run((&Report{AggregateBy: "month", DiffArgs: lib.DiffArgs{Diff: true}, Fill: true}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
            Total    Should     Diff
2018 Feb       8h       8h!       0m
     Mar                            
     Apr       2h    5h30m!   -3h30m
     May    5h20m       2h!   +3h20m
     Jun                            
     Jul                            
     Aug                            
     Sep                            
     Oct                            
     Nov                            
     Dec                            
2019 Jan       0m      19m!     -19m
         ======== ========= ========
           15h20m   15h49m!     -29m
`, state.printBuffer)
}

func TestYearReport(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
2016-02-02 (8h!)
	8h

2018-04-10 (5h30m!)
	2h

2018-05-23 (2h!)
	5h20m

2019-01-01 (19m!)
`)._Run((&Report{AggregateBy: "year", DiffArgs: lib.DiffArgs{Diff: true}, Fill: true}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
        Total    Should     Diff
2016       8h       8h!       0m
2017                            
2018    7h20m    7h30m!     -10m
2019       0m      19m!     -19m
     ======== ========= ========
       15h20m   15h49m!     -29m
`, state.printBuffer)
}
