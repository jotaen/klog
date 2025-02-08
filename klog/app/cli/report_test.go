package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app/cli/util"
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
		FilterArgs: util.FilterArgs{Date: klog.Ɀ_Date_(2022, 10, 31)},
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
`)._SetNow(2021, 3, 4, 0, 0)._Run((&Report{WarnArgs: util.WarnArgs{NoWarn: true}}).Run)
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
`)._Run((&Report{DiffArgs: util.DiffArgs{Diff: true}}).Run)
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
`)._Run((&Report{DiffArgs: util.DiffArgs{Diff: true}, DecimalArgs: util.DecimalArgs{Decimal: true}}).Run)
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
`)._Run((&Report{AggregateBy: "week", DiffArgs: util.DiffArgs{Diff: true}, Fill: true}).Run)
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
`)._Run((&Report{AggregateBy: "quarter", DiffArgs: util.DiffArgs{Diff: true}, Fill: true}).Run)
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
`)._Run((&Report{AggregateBy: "month", DiffArgs: util.DiffArgs{Diff: true}, Fill: true}).Run)
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
`)._Run((&Report{AggregateBy: "year", DiffArgs: util.DiffArgs{Diff: true}, Fill: true}).Run)
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

func TestReportWithChart(t *testing.T) {
	t.Run("Daily (default) aggregation", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
2025-01-11
	-2h

2025-01-11
	0m

2025-01-13
	1m

2025-01-14
	5h

2025-01-16
	5h1m

2025-01-17
	5h15m

2025-01-18
	5h30m

2025-01-20
	5h51m

2025-01-22
	6h

2025-01-25
	9h
`)._Run((&Report{Chart: true}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
                       Total                                      
2025 Jan    Sat 11.      -2h                                      
            Mon 13.       1m  ▇                                   
            Tue 14.       5h  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇                
            Thu 16.     5h1m  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇               
            Fri 17.    5h15m  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇               
            Sat 18.    5h30m  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇              
            Mon 20.    5h51m  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇            
            Wed 22.       6h  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇            
            Sat 25.       9h  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇
                    ========                                      
                      39h38m                                      
`, state.printBuffer)
	})

	t.Run("Weekly aggregation", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
2025-01-01
	40h

2025-01-08
	48h45m

2025-01-15
	31h15m
`)._Run((&Report{Chart: true, AggregateBy: "w", WarnArgs: util.WarnArgs{NoWarn: true}}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
                 Total                                                   
2025  Week  1      40h  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇         
      Week  2   48h45m  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇
      Week  3   31h15m  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇                 
              ========                                                   
                  120h                                                   
`, state.printBuffer)
	})

	t.Run("Monthly aggregation", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
2025-01-01
	173h

2025-02-01
	208h30m

2025-03-01
	126h15m
`)._Run((&Report{Chart: true, AggregateBy: "m", WarnArgs: util.WarnArgs{NoWarn: true}}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
            Total                                                       
2025 Jan     173h  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇         
     Feb  208h30m  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇
     Mar  126h15m  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇                     
         ========                                                       
          507h45m                                                       
`, state.printBuffer)
	})

	t.Run("Quarterly aggregation", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
2025-01-01
	316h

2025-04-01
	392h30m

2025-07-01
	237h45m
`)._Run((&Report{Chart: true, AggregateBy: "q", WarnArgs: util.WarnArgs{NoWarn: true}}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
           Total                                                    
2025 Q1     316h  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇          
     Q2  392h30m  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇
     Q3  237h45m  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇                    
        ========                                                    
         946h15m                                                    
`, state.printBuffer)
	})

	t.Run("Yearly aggregation", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
2025-01-01
	1735h

2026-01-01
	2154h45m

2027-01-01
	1189h15m
`)._Run((&Report{Chart: true, AggregateBy: "y", WarnArgs: util.WarnArgs{NoWarn: true}}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
        Total                                         
2025    1735h  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇        
2026 2154h45m  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇
2027 1189h15m  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇                 
     ========                                         
        5079h                                         
`, state.printBuffer)
	})
}

func TestReportWithScaledChart(t *testing.T) {
	t.Run("Custom resolution (large)", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
2025-01-14
	12h

2025-01-16
	18h37m
`)._Run((&Report{Chart: true, ChartResolution: 120}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
                       Total            
2025 Jan    Tue 14.      12h  ▇▇▇▇▇▇    
            Thu 16.   18h37m  ▇▇▇▇▇▇▇▇▇▇
                    ========            
                      30h37m            
`, state.printBuffer)
	})

	t.Run("Custom resolution (small)", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
2025-01-14
	1h30m

2025-01-16
	45m
`)._Run((&Report{Chart: true, ChartResolution: 5}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
                       Total                    
2025 Jan    Tue 14.    1h30m  ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇
            Thu 16.      45m  ▇▇▇▇▇▇▇▇▇         
                    ========                    
                       2h15m                    
`, state.printBuffer)
	})

	t.Run("Setting resolution implies --chart", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
2025-01-14
	12h

2025-01-16
	18h37m
`)._Run((&Report{ChartResolution: 120}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
                       Total            
2025 Jan    Tue 14.      12h  ▇▇▇▇▇▇    
            Thu 16.   18h37m  ▇▇▇▇▇▇▇▇▇▇
                    ========            
                      30h37m            
`, state.printBuffer)
	})

	t.Run("With --diff flag", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
2025-01-14 (2h!)
	1h30m

2025-01-16 (1h!)
	45m
`)._Run((&Report{Chart: true, DiffArgs: util.DiffArgs{Diff: true}}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
                       Total    Should     Diff        
2025 Jan    Tue 14.    1h30m       2h!     -30m  ▇▇▇▇▇▇
            Thu 16.      45m       1h!     -15m  ▇▇▇   
                    ======== ========= ========        
                       2h15m       3h!     -45m        
`, state.printBuffer)
	})

	t.Run("Invalid resolution", func(t *testing.T) {
		_, err := NewTestingContext()._SetRecords(`
2025-01-14
	12h

2025-01-16
	18h37m
`)._Run((&Report{ChartResolution: -10}).Run)
		require.Error(t, err)
		assert.Equal(t, "Invalid resolution", err.Error())
	})
}
