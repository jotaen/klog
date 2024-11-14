package klog

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestHandleInputFiles(t *testing.T) {
	out := (&Env{
		files: map[string]string{
			"test.klg": "2020-01-01\nSome stuff\n\t1h\n",
			"foo.klg":  "2021-03-02\n    2h #foo",
		},
	}).run(
		[]string{"print", "test.klg", "foo.klg"},
		[]string{"tags", "foo.klg"},
	)
	// Out 0 like: `2020-01-01\nSome stuff\n    1h\n\n2021-03-02\n    2h`
	assert.True(t, strings.Contains(out[0], "2020-01-01"), out)
	assert.True(t, strings.Contains(out[0], "2021-03-02"), out)
	// Out 1 like: `#foo 2h`
	assert.True(t, strings.Contains(out[1], "#foo 2h"), out)
}

func TestHandlesInvocationErrors(t *testing.T) {
	out := (&Env{
		files: map[string]string{},
	}).run(
		[]string{"print", "--foo"},
	)
	// Invocation error: unknown flag --foo
	assert.True(t, strings.Contains(out[0], "Invocation error: unknown flag --foo"), out)
}

func TestPrintAppErrors(t *testing.T) {
	out := (&Env{
		files: map[string]string{
			"invalid.klg": "2020-01-01asdf",
			"valid.klg":   "2020-01-01",
		},
	}).run(
		[]string{"print", "invalid.klg"},
		[]string{"start", "valid.klg"},
		[]string{"start", "valid.klg"},
	)
	// Out 0 should contain pretty-printed parsing errors.
	assert.True(t, strings.Contains(out[0], "[SYNTAX ERROR] in line 1 of file"), out)
	assert.True(t, strings.Contains(out[0], "invalid.klg"), out)
	assert.True(t, strings.Contains(out[0], "2020-01-01asdf"), out)
	assert.True(t, strings.Contains(out[0], "^^^^^^^^^^^^^^"), out)
	assert.True(t, strings.Contains(out[0], "Invalid date"), out)
	// Out 1 should go through without errors.
	// Out 2 should then display logical error, since there is an open-range already.
	assert.True(t, strings.Contains(out[2], "Error: Manipulation failed"), out)
	assert.True(t, strings.Contains(out[2], "There is already an open range in this record"), out)
}

func TestConfigureAndUseBookmark(t *testing.T) {
	klog := &Env{
		files: map[string]string{
			"test.klg": "2020-01-01\nSome stuff\n\t1h7m\n",
		},
	}
	out := klog.run(
		[]string{"bookmarks", "set", "test.klg", "tst"},
		[]string{"bookmarks", "set", "test.klg", "tst"},
		[]string{"bookmarks", "list"},
		[]string{"total", "@tst"},
	)
	// Out 0 like: `Created new bookmark`, `@tst -> /tmp/.../test.klg`
	assert.True(t, strings.Contains(out[0], "Created new bookmark"), out)
	assert.True(t, strings.Contains(out[0], "@tst"), out)
	assert.True(t, strings.Contains(out[0], "test.klg"), out)
	// Out 1 like: `Changed bookmark`, `@tst -> /tmp/.../test.klg`
	assert.True(t, strings.Contains(out[1], "Changed bookmark"), out)
	assert.True(t, strings.Contains(out[1], "@tst"), out)
	// Out 2 like: `@tst -> /tmp/.../test.klg`
	assert.True(t, strings.Contains(out[2], "@tst"), out)
	// Out 3 like: `Total: 1h7m`
	assert.True(t, strings.Contains(out[3], "1h7m"), out)
}

func TestCreateBookmarkTargetFileOnDemand(t *testing.T) {
	klog := &Env{}
	out := klog.run(
		[]string{"bookmarks", "set", "--create", "test.klg", "tst"},
		[]string{"bookmarks", "set", "--create", "test.klg", "tst"},
	)
	// Out 0 like: `Created new bookmark`, `@tst -> /tmp/.../test.klg`
	assert.True(t, strings.Contains(out[0], "Created new bookmark and created target file:"), out)
	assert.True(t, strings.Contains(out[0], "@tst"), out)
	assert.True(t, strings.Contains(out[0], "test.klg"), out)
	// Out 1 like: `Error: Cannot create file`, `There is already a file at that location`
	assert.True(t, strings.Contains(out[1], "Error: Cannot create file"), out)
	assert.True(t, strings.Contains(out[1], "There is already a file at that location"), out)
}

func TestWriteToFile(t *testing.T) {
	klog := &Env{
		files: map[string]string{
			"test.klg": "2020-01-01\nSome stuff\n\t1h\n",
		},
	}
	out := klog.run(
		[]string{"track", "--date", "2020-01-01", "30m", "test.klg"},
		[]string{"total", "test.klg"},
	)
	// Out 1 like: `Total 1h30m (In 1 record)`
	assert.True(t, strings.Contains(out[1], "1h30m"), out)
	assert.True(t, strings.Contains(out[1], "1 record"), out)
}

func TestDecodesDate(t *testing.T) {
	klog := &Env{
		files: map[string]string{
			"test.klg": "2020-01-01\nSome stuff\n\t1h7m\n",
		},
	}
	out := klog.run(
		[]string{"total", "--date", "2020-1-1", "test.klg"},
		[]string{"total", "--date", "2020-01-01", "test.klg"},
	)
	assert.True(t, strings.Contains(out[0], "`2020-1-1` is not a valid date"), out)
	assert.True(t, strings.Contains(out[1], "1h7m"), out)
}

func TestDecodesTime(t *testing.T) {
	klog := &Env{
		files: map[string]string{
			"test.klg": "2020-01-01\n\t9:00-?\n",
		},
	}
	out := klog.run(
		[]string{"stop", "--date", "2020-01-01", "--time", "1:0", "test.klg"},
		[]string{"stop", "--date", "2020-01-01", "--time", "10:00", "test.klg"},
		[]string{"total", "test.klg"},
	)
	assert.True(t, strings.Contains(out[0], "`1:0` is not a valid time"), out)
	assert.True(t, strings.Contains(out[1], "9:00-10:00"), out)
	assert.True(t, strings.Contains(out[2], "1h"), out)
}

func TestDecodesShouldTotal(t *testing.T) {
	klog := &Env{
		files: map[string]string{
			"test.klg": "",
		},
	}
	out := klog.run(
		[]string{"create", "--date", "2020-01-01", "--should", "asdf", "test.klg"},
		[]string{"create", "--date", "2020-01-01", "--should", "5h1m!", "test.klg"},
		[]string{"total", "--diff", "test.klg"},
	)
	assert.True(t, strings.Contains(out[0], "`asdf` is not a valid should total"), out)
	assert.True(t, strings.Contains(out[1], "5h1m!"), out)
	assert.True(t, strings.Contains(out[2], "5h1m!"), out)
}

func TestDecodesPeriod(t *testing.T) {
	klog := &Env{
		files: map[string]string{
			"test.klg": "2000-01-05\n\t1h\n\n2000-05-24\n\t1h\n",
		},
	}
	out := klog.run(
		[]string{"total", "--period", "2000", "test.klg"},
		[]string{"total", "--period", "2000-01", "test.klg"},
		[]string{"total", "--period", "2000-Q1", "test.klg"},
		[]string{"total", "--period", "2000-W21", "test.klg"},
		[]string{"total", "--period", "foo", "test.klg"},
	)
	assert.True(t, strings.Contains(out[0], "2h"), out)
	assert.True(t, strings.Contains(out[1], "1h"), out)
	assert.True(t, strings.Contains(out[2], "1h"), out)
	assert.True(t, strings.Contains(out[3], "1h"), out)
	assert.True(t, strings.Contains(out[4], "`foo` is not a valid period"), out)
}

func TestDecodesRounding(t *testing.T) {
	klog := &Env{
		files: map[string]string{
			"test.klg": "2020-01-01",
		},
	}
	out := klog.run(
		[]string{"start", "--round", "asdf", "test.klg"},
		[]string{"start", "--round", "30m", "test.klg"},
	)
	assert.True(t, strings.Contains(out[0], "`asdf` is not a valid rounding value"), out)
	assert.True(t, strings.Contains(out[1], "- ?"), out)
}

func TestDecodesTags(t *testing.T) {
	klog := &Env{
		files: map[string]string{
			"test.klg": "2020-01-01\n#foo\n\n2020-01-02\n\t1h #bar=1",
		},
	}
	out := klog.run(
		[]string{"print", "--tag", "asdf=asdf=asdf", "test.klg"},
		[]string{"print", "--tag", "foo&bar", "test.klg"},
		[]string{"print", "--tag", "foo", "test.klg"},
		[]string{"print", "--tag", "bar=1", "test.klg"},
		[]string{"print", "--tag", "#bar='1'", "test.klg"},
	)
	assert.True(t, strings.Contains(out[0], "`asdf=asdf=asdf` is not a valid tag"), out)
	assert.True(t, strings.Contains(out[1], "`foo&bar` is not a valid tag"), out)
	assert.True(t, strings.Contains(out[2], "#foo"), out)
	assert.True(t, strings.Contains(out[3], "#bar=1"), out)
	assert.True(t, strings.Contains(out[4], "#bar=1"), out)
}

func TestDecodesRecordSummary(t *testing.T) {
	klog := &Env{
		files: map[string]string{
			"test.klg": "2020-01-01\nTest.",
		},
	}
	out := klog.run(
		[]string{"create", "--summary", "Foo", "test.klg"},
		[]string{"create", "--summary", "Foo\nBar", "test.klg"},
		[]string{"create", "--summary", "Foo\n\nBar", "test.klg"},
		[]string{"create", "--summary", "Foo\n Bar", "test.klg"},
	)
	assert.True(t, strings.Contains(out[0], "Foo"), out)
	assert.True(t, strings.Contains(out[1], "Foo\nBar"), out)
	assert.True(t, strings.Contains(out[2], "A record summary cannot contain blank lines"), out)
	assert.True(t, strings.Contains(out[3], "A record summary cannot contain blank lines"), out)
}

func TestDecodesEntryType(t *testing.T) {
	klog := &Env{
		files: map[string]string{
			"test.klg": "2020-01-01\n\t1h\n\t9:00-12:00",
		},
	}
	out := klog.run(
		[]string{"total", "--entry-type", "duration", "test.klg"},
		[]string{"total", "--entry-type", "DURATION", "test.klg"},
		[]string{"total", "--entry-type", "duration-positive", "test.klg"},
		[]string{"total", "--entry-type", "duration-negative", "test.klg"},
		[]string{"total", "--entry-type", "open_range", "test.klg"},
		[]string{"total", "--entry-type", "open-range", "test.klg"},
		[]string{"total", "--entry-type", "range", "test.klg"},
		[]string{"total", "--entry-type", "asdf", "test.klg"},
	)
	assert.True(t, strings.Contains(out[0], "1h"), out)
	assert.True(t, strings.Contains(out[1], "1h"), out)
	assert.True(t, strings.Contains(out[2], "1h"), out)
	assert.True(t, strings.Contains(out[3], "0m"), out)
	assert.True(t, strings.Contains(out[4], "0m"), out)
	assert.True(t, strings.Contains(out[5], "0m"), out)
	assert.True(t, strings.Contains(out[6], "3h"), out)
	assert.True(t, strings.Contains(out[7], "is not a valid entry type"), out)
}
