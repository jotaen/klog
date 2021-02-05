package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReportOfEmptyInput(t *testing.T) {
	out, err := RunWithContext(`

`, (&Report{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "", out)
}
