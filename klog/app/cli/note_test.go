package cli

import (
	"testing"

	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoteAppendsInlineToOpenRange(t *testing.T) {
	t.Run("no existing summary", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00 - ?
	1h Other
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Note{
			Text: klog.Ɀ_EntrySummary_("Working on something #tag"),
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920-02-02
	9:00 - ? Working on something #tag
	1h Other
`, state.writtenFileContents)
	})

	t.Run("existing single-line summary", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-? Working
	1h Other
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Note{
			Text: klog.Ɀ_EntrySummary_("on something #tag"),
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920-02-02
	9:00-? Working on something #tag
	1h Other
`, state.writtenFileContents)
	})

	t.Run("existing multiline summary", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-? Working
		on something
	1h Other
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Note{
			Text: klog.Ɀ_EntrySummary_("#tag"),
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920-02-02
	9:00-? Working
		on something #tag
	1h Other
`, state.writtenFileContents)
	})
}

func TestNoteAddsNewLineToOpenRange(t *testing.T) {
	t.Run("no existing summary", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-?
	1h Other
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Note{
			Text:       klog.Ɀ_EntrySummary_("Did something #tag"),
			AppendMode: "newline",
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920-02-02
	9:00-?
		Did something #tag
	1h Other
`, state.writtenFileContents)
	})

	t.Run("existing single-line summary", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-? Working
	1h Other
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Note{
			Text:       klog.Ɀ_EntrySummary_("Did something #tag"),
			AppendMode: "newline",
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920-02-02
	9:00-? Working
		Did something #tag
	1h Other
`, state.writtenFileContents)
	})

	t.Run("existing multline summary", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-? Working
		More Working
	1h Other
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Note{
			Text:       klog.Ɀ_EntrySummary_("Did something #tag"),
			AppendMode: "newline",
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920-02-02
	9:00-? Working
		More Working
		Did something #tag
	1h Other
`, state.writtenFileContents)
	})
}

func TestNoteTargetsNthEntry(t *testing.T) {
	records := `
1920-02-02
	8:00-9:00 First
	9:00-? Second
	1h Third
`
	t.Run("Counting from the start", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(records)._SetNow(1920, 2, 2, 15, 24)._Run((&Note{
			Text:       klog.Ɀ_EntrySummary_("Note"),
			Entry:      1,
			AppendMode: "newline",
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920-02-02
	8:00-9:00 First
		Note
	9:00-? Second
	1h Third
`, state.writtenFileContents)
	})

	t.Run("Counting from the end", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(records)._SetNow(1920, 2, 2, 15, 24)._Run((&Note{
			Text:  klog.Ɀ_EntrySummary_("and more"),
			Entry: -1,
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920-02-02
	8:00-9:00 First
	9:00-? Second
	1h Third and more
`, state.writtenFileContents)
	})

	t.Run("Out of range", func(t *testing.T) {
		for _, entry := range []int{4, -4} {
			_, err := NewTestingContext()._SetRecords(records)._SetNow(1920, 2, 2, 15, 24)._Run((&Note{
				Text:  klog.Ɀ_EntrySummary_("Note"),
				Entry: entry,
			}).Run)
			require.Error(t, err)
		}
	})
}

func TestNoteFallsBackToYesterday(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	22:00-?
`)._SetNow(1920, 2, 3, 1, 10)._Run((&Note{
		Text: klog.Ɀ_EntrySummary_("Still working"),
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	22:00-? Still working
`, state.writtenFileContents)
}

func TestNoteFailsWithoutOpenRange(t *testing.T) {
	_, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-10:00
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Note{
		Text: klog.Ɀ_EntrySummary_("Text"),
	}).Run)
	require.Error(t, err)
}
