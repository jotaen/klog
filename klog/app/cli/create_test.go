package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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
		AtDateArgs:  lib.AtDateArgs{Date: klog.Ɀ_Date_(1976, 1, 2)},
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
	// For empty file and no preferences, use recommended default.
	{
		state, err := NewTestingContext()._SetRecords(``).
			_SetNow(1920, 2, 3, 15, 24)._Run((&Create{}).Run)
		require.Nil(t, err)
		assert.Equal(t, "1920-02-03\n", state.writtenFileContents)
	}

	// Without any preference, detect from file.
	{
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
	}

	// Use preference from config file, if given.
	{
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
	}

	// If explicit flag was provided, that takes ultimate precedence.
	{
		state, err := NewTestingContext()._SetRecords(`
1920/02/01
	4h33m

1920/02/02
	9:00-12:00
`)._SetFileConfig(`
date_format = YYYY/MM/DD
`)._SetNow(1920, 2, 3, 15, 24)._Run((&Create{
			AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 3)},
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920/02/01
	4h33m

1920/02/02
	9:00-12:00

1920-02-03
`, state.writtenFileContents)
	}
}

func TestCreateWithShouldTotalConfig(t *testing.T) {
	// With should-total from config file
	{
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
	}

	// --should-total flag trumps should-total from config file
	{
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
	}
}
