package parser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseEmptyDocument(t *testing.T) {
	text := ``
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Nil(t, rs)
}
