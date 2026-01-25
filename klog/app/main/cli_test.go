package klog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleInputFiles(t *testing.T) {
	(&Env{
		files: map[string]string{
			"test.klg": "2020-01-01\nSome stuff\n\t1h\n",
			"foo.klg":  "2021-03-02\n    2h #foo",
		},
	}).execute(t,
		invocation{
			args: []string{"print", "test.klg", "foo.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "2020-01-01"), out)
				assert.True(t, strings.Contains(out, "2021-03-02"), out)
			}},
		invocation{
			args: []string{"tags", "foo.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "#foo 2h"), out)
			}},
	)
}

func TestHandlesInvocationErrors(t *testing.T) {
	(&Env{
		files: map[string]string{},
	}).execute(t,
		invocation{
			args: []string{"print", "--foo"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 1, code)
				assert.True(t, strings.Contains(out, "Invocation error: unknown flag --foo"), out)
			}},
	)
}

func TestPrintAppErrors(t *testing.T) {
	(&Env{
		files: map[string]string{
			"invalid.klg": "2020-01-01asdf",
			"valid.klg":   "2020-01-01",
		},
	}).execute(t,
		invocation{
			args: []string{"print", "invalid.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 8, code)
				assert.True(t, strings.Contains(out, "[SYNTAX ERROR] in line 1 of file"), out)
				assert.True(t, strings.Contains(out, "invalid.klg"), out)
				assert.True(t, strings.Contains(out, "2020-01-01asdf"), out)
				assert.True(t, strings.Contains(out, "^^^^^^^^^^^^^^"), out)
				assert.True(t, strings.Contains(out, "Invalid date"), out)
			}},
		invocation{
			args: []string{"start", "valid.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
			}},
		invocation{
			args: []string{"start", "valid.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 8, code)
				assert.True(t, strings.Contains(out, "Error: Manipulation failed"), out)
				assert.True(t, strings.Contains(out, "There is already an open range in this record"), out)
			}},
		invocation{
			args: []string{"print", "--filter", "2020-01 #foo", "valid.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 1, code)
				assert.True(t, strings.Contains(out, "Missing operator"), out)
				assert.True(t, strings.Contains(out, "2020-01 #foo"), out)
				assert.True(t, strings.Contains(out, "————————^^^^"), out)
				assert.True(t, strings.Contains(out, "Cursor positions 8-12 in query"), out)
			}},
	)
}

func TestConfigureAndUseBookmark(t *testing.T) {
	(&Env{
		files: map[string]string{
			"test.klg": "2020-01-01\nSome stuff\n\t1h7m\n",
		},
	}).execute(t,
		invocation{
			args: []string{"bookmarks", "set", "test.klg", "tst"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "Created new bookmark"), out)
				assert.True(t, strings.Contains(out, "@tst"), out)
				assert.True(t, strings.Contains(out, "test.klg"), out)
			}},
		invocation{
			args: []string{"bookmarks", "set", "test.klg", "tst"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "Changed bookmark"), out)
				assert.True(t, strings.Contains(out, "@tst"), out)
			}},
		invocation{
			args: []string{"bookmarks", "list"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "@tst"), out)
			}},
		invocation{
			args: []string{"total", "@tst"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "1h7m"), out)
			}},
	)
}

func TestCreateBookmarkTargetFileOnDemand(t *testing.T) {
	(&Env{
		files: map[string]string{},
	}).execute(t,
		invocation{
			args: []string{"bookmarks", "set", "--create", "test.klg", "tst"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "Created new bookmark and created target file:"), out)
				assert.True(t, strings.Contains(out, "@tst"), out)
				assert.True(t, strings.Contains(out, "test.klg"), out)
			}},
		invocation{
			args: []string{"bookmarks", "set", "--create", "test.klg", "tst"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 1, code)
				assert.True(t, strings.Contains(out, "Error: Cannot create file"), out)
				assert.True(t, strings.Contains(out, "There is already a file at that location"), out)
			}},
	)
}

func TestWriteToFile(t *testing.T) {
	(&Env{
		files: map[string]string{
			"test.klg": "2020-01-01\nSome stuff\n\t1h\n",
		},
	}).execute(t,
		invocation{
			args: []string{"track", "--date", "2020-01-01", "30m", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
			}},
		invocation{
			args: []string{"total", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "1h30m"), out)
				assert.True(t, strings.Contains(out, "1 record"), out)
			}},
	)
}

func TestDecodesDate(t *testing.T) {
	(&Env{
		files: map[string]string{
			"test.klg": "2020-01-01\nSome stuff\n\t1h7m\n",
		},
	}).execute(t,
		invocation{
			args: []string{"total", "--date", "2020-1-1", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 1, code)
				assert.True(t, strings.Contains(out, "`2020-1-1` is not a valid date"), out)
			}},
		invocation{
			args: []string{"total", "--date", "2020-01-01", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "1h7m"), out)
			}},
	)
}

func TestDecodesTime(t *testing.T) {
	(&Env{
		files: map[string]string{
			"test.klg": "2020-01-01\n\t9:00-?\n",
		},
	}).execute(t,
		invocation{
			args: []string{"stop", "--date", "2020-01-01", "--time", "1:0", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 1, code)
				assert.True(t, strings.Contains(out, "`1:0` is not a valid time"), out)
			}},
		invocation{
			args: []string{"stop", "--date", "2020-01-01", "--time", "10:00", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "9:00-10:00"), out)
			}},
		invocation{
			args: []string{"total", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "1h"), out)
			}},
	)
}

func TestDecodesShouldTotal(t *testing.T) {
	(&Env{
		files: map[string]string{
			"test.klg": "",
		},
	}).execute(t,
		invocation{
			args: []string{"create", "--date", "2020-01-01", "--should", "asdf", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 1, code)
				assert.True(t, strings.Contains(out, "`asdf` is not a valid should total"), out)
			}},
		invocation{
			args: []string{"create", "--date", "2020-01-01", "--should", "5h1m!", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "5h1m!"), out)
			}},
		invocation{
			args: []string{"total", "--diff", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "5h1m!"), out)
			}},
	)
}

func TestDecodesPeriod(t *testing.T) {
	(&Env{
		files: map[string]string{
			"test.klg": "2000-01-05\n\t1h\n\n2000-05-24\n\t1h\n",
		},
	}).execute(t,
		invocation{
			args: []string{"total", "--period", "2000", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "2h"), out)
			}},
		invocation{
			args: []string{"total", "--period", "2000-01", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "1h"), out)
			}},
		invocation{
			args: []string{"total", "--period", "2000-Q1", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "1h"), out)
			}},
		invocation{
			args: []string{"total", "--period", "2000-W21", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "1h"), out)
			}},
		invocation{
			args: []string{"total", "--period", "foo", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 1, code)
				assert.True(t, strings.Contains(out, "`foo` is not a valid period"), out)
			}},
	)
}

func TestDecodesRounding(t *testing.T) {
	(&Env{
		files: map[string]string{
			"test.klg": "2020-01-01",
		},
	}).execute(t,
		invocation{
			args: []string{"start", "--round", "asdf", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 1, code)
				assert.True(t, strings.Contains(out, "`asdf` is not a valid rounding value"), out)
			}},
		invocation{
			args: []string{"start", "--round", "30m", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "- ?"), out)
			}},
	)
}

func TestDecodesTags(t *testing.T) {
	(&Env{
		files: map[string]string{
			"test.klg": "2020-01-01\n#foo\n\t2h\n\n2020-01-02\n\t1h #bar=1",
		},
	}).execute(t,
		invocation{
			args: []string{"print", "--tag", "asdf=asdf=asdf", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 1, code)
				assert.True(t, strings.Contains(out, "`asdf=asdf=asdf` is not a valid tag"), out)
			}},
		invocation{
			args: []string{"print", "--tag", "foo&bar", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 1, code)
				assert.True(t, strings.Contains(out, "`foo&bar` is not a valid tag"), out)
			}},
		invocation{
			args: []string{"print", "--tag", "foo", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "#foo"), out)
			}},
		invocation{
			args: []string{"print", "--tag", "bar=1", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "#bar=1"), out)
			}},
		invocation{
			args: []string{"print", "--tag", "#bar='1'", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "#bar=1"), out)
			}},
	)
}

func TestDecodesRecordSummary(t *testing.T) {
	(&Env{
		files: map[string]string{
			"test.klg": "2020-01-01\nTest.",
		},
	}).execute(t,
		invocation{
			args: []string{"create", "--summary", "Foo", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "Foo"), out)
			}},
		invocation{
			args: []string{"create", "--summary", "Foo\nBar", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "Foo\nBar"), out)
			}},
		invocation{
			args: []string{"create", "--summary", "Foo\n\nBar", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 1, code)
				assert.True(t, strings.Contains(out, "A record summary cannot contain blank lines"), out)
			}},
		invocation{
			args: []string{"create", "--summary", "Foo\n Bar", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 1, code)
				assert.True(t, strings.Contains(out, "A record summary cannot contain blank lines"), out)
			}},
	)
}

func TestDecodesEntryType(t *testing.T) {
	(&Env{
		files: map[string]string{
			"test.klg": "2020-01-01\n\t1h\n\t9:00-12:00",
		},
	}).execute(t,
		invocation{
			args: []string{"total", "--entry-type", "duration", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "1h"), out)
			}},
		invocation{
			args: []string{"total", "--entry-type", "DURATION", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "1h"), out)
			}},
		invocation{
			args: []string{"total", "--entry-type", "duration-positive", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "1h"), out)
			}},
		invocation{
			args: []string{"total", "--entry-type", "duration-negative", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "0m"), out)
			}},
		invocation{
			args: []string{"total", "--entry-type", "open_range", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "0m"), out)
			}},
		invocation{
			args: []string{"total", "--entry-type", "open-range", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "0m"), out)
			}},
		invocation{
			args: []string{"total", "--entry-type", "range", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 0, code)
				assert.True(t, strings.Contains(out, "3h"), out)
			}},
		invocation{
			args: []string{"total", "--entry-type", "asdf", "test.klg"},
			test: func(t *testing.T, code int, out string) {
				assert.Equal(t, 1, code)
				assert.True(t, strings.Contains(out, "is not a valid entry type"), out)
			}},
	)
}
