package cli

import (
	"testing"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app/cli/args"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-01
	4h33m

1920-02-02
	9:00-12:00
`)._SetNow(1920, 2, 3, 15, 24)._Run((&Create{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-01
	4h33m

1920-02-02
	9:00-12:00

1920-02-03
`, state.writtenFileContents)
}

func TestCreateWithValues(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1975-12-31
	1h
	2h


1976-01-01
	1h
`)._Run((&Create{
		AtDateArgs:  args.AtDateArgs{Date: klog.Ɀ_Date_(1976, 1, 2)},
		ShouldTotal: klog.NewShouldTotal(5, 55),
		Summary:     klog.Ɀ_RecordSummary_("This is a", "new record!"),
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1975-12-31
	1h
	2h


1976-01-01
	1h

1976-01-02 (5h55m!)
This is a
new record!
`, state.writtenFileContents)
}

func TestCreateWithShouldTotalAlias(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(``).
		_SetNow(1976, 1, 1, 2, 2)._Run((&Create{
		ShouldTotalAlias: klog.NewShouldTotal(0, 30),
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, "1976-01-01 (30m!)\n", state.writtenFileContents)
}

func TestCreateWithStyle(t *testing.T) {
	t.Run("For empty file and no preferences, use recommended default.", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(``).
			_SetNow(1920, 2, 3, 15, 24)._Run((&Create{}).Run)
		require.Nil(t, err)
		assert.Equal(t, "1920-02-03\n", state.writtenFileContents)
	})

	t.Run("Without any preference, detect from file.", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1920/02/01
	4h33m

1920/02/02
	9:00-12:00
`)._SetNow(1920, 2, 3, 15, 24)._Run((&Create{}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920/02/01
	4h33m

1920/02/02
	9:00-12:00

1920/02/03
`, state.writtenFileContents)
	})

	t.Run("Use preference from config file, if given.", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1920/02/01
	4h33m

1920/02/02
	9:00-12:00
`)._SetFileConfig(`
date_format = YYYY-MM-DD
`)._SetNow(1920, 2, 3, 15, 24)._Run((&Create{}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920/02/01
	4h33m

1920/02/02
	9:00-12:00

1920-02-03
`, state.writtenFileContents)
	})

	t.Run("If explicit flag was provided, that takes ultimate precedence.", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1920/02/01
	4h33m

1920/02/02
	9:00-12:00
`)._SetFileConfig(`
date_format = YYYY/MM/DD
`)._SetNow(1920, 2, 3, 15, 24)._Run((&Create{
			AtDateArgs: args.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 3)},
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920/02/01
	4h33m

1920/02/02
	9:00-12:00

1920-02-03
`, state.writtenFileContents)
	})
}

func TestCreateWithShouldTotalConfig(t *testing.T) {
	t.Run("With should-total from config file", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1920-02-01
	4h33m

1920-02-02
	9:00-12:00
`)._SetFileConfig(`
default_should_total = 30m!
`)._SetNow(1920, 2, 3, 15, 24)._Run((&Create{}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920-02-01
	4h33m

1920-02-02
	9:00-12:00

1920-02-03 (30m!)
`, state.writtenFileContents)
	})

	t.Run("--should-total flag trumps should-total from config file", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1920-02-01
	4h33m

1920-02-02
	9:00-12:00
`)._SetFileConfig(`
default_should_total = 30m!
`)._SetNow(1920, 2, 3, 15, 24)._Run((&Create{
			ShouldTotal: klog.NewShouldTotal(5, 55),
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920-02-01
	4h33m

1920-02-02
	9:00-12:00

1920-02-03 (5h55m!)
`, state.writtenFileContents)
	})
}
