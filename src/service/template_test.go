package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var now = time.Date(1995, 3, 31, 13, 15, 29, 0, time.UTC)

func TestRenderTemplate(t *testing.T) {
	result, err := RenderTemplate(`
{{ TODAY }}
Foo #xyz

{{YESTERDAY}} (8h30m!)
	1h
	{{ NOW }} - ?
`, now)
	require.Nil(t, err)
	assert.Equal(t, RecordText(`
1995-03-31
Foo #xyz

1995-03-30 (8h30m!)
	1h
	13:15 - ?
`), result)
}

func TestTemplateFailsIfNoValidRecord(t *testing.T) {
	result, err := RenderTemplate(`
{{ TODAY }} foo
	This is all malformed
`, now)
	require.Error(t, err)
	assert.Equal(t, RecordText(""), result)
}

func TestAppendTemplateToRecord(t *testing.T) {
	newRecord := RecordText("1659-08-02 (8h!)\nMy template\n    1h")
	for _, x := range []struct {
		file string
		glue string
	}{
		{"", ""},
		{"   ", "\n\n"},
		{"1659-08-01\n    14:00-19:00 Appointment\n\n\n", ""},
		{"1659-08-01\n    14:00-19:00 Appointment\n\n", ""},
		{"1659-08-01\n    14:00-19:00 Appointment\n    \n", "\n"},
		{"1659-08-01\n    14:00-19:00 Appointment\n", "\n"},
		{"1659-08-01\n    14:00-19:00 Appointment", "\n\n"},
	} {
		txt := AppendableText(x.file, newRecord)
		require.Equal(t, x.glue+"1659-08-02 (8h!)\nMy template\n    1h", txt)
	}
}
