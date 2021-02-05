package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTotalOfEmptyInput(t *testing.T) {
	out, err := RunWithContext(``, (&Total{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "Total: 0m\n(In 0 records)\n", out)
}

func TestTotalOfInput(t *testing.T) {
	out, err := RunWithContext(`
2018-11-08
	1h

2018-11-09
	16:00-17:00

2018-11-10
Open ranges are not considered
	16:00 - ?
`, (&Total{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "Total: 2h\n(In 3 records)\n", out)
}

func TestTotalWithDiffing(t *testing.T) {
	out, err := RunWithContext(`
2018-11-08 (8h!)
	8h30m

2018-11-09 (7h45m!)
	8:00 - 16:00
`, (&Total{Diff: true}).Run)
	require.Nil(t, err)
	assert.Equal(t, "Total: 16h30m\nShould: 15h45m!\nDiff: +45m\n(In 2 records)\n", out)
}
