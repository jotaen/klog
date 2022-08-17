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
