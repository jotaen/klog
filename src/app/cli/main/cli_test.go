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

func TestBookmarkFile(t *testing.T) {
	klog := &Env{
		files: map[string]string{
			"test.klg": "2020-01-01\nSome stuff\n\t1h7m\n",
		},
	}
	out := klog.run(
		[]string{"bookmarks", "set", "test.klg", "tst"},
		[]string{"bookmarks", "list"},
		[]string{"total", "@tst"},
	)
	// Out 1 like: `@tst -> /tmp/.../test.klg`
	assert.True(t, strings.Contains(out[1], "@tst"), out)
	assert.True(t, strings.Contains(out[1], "test.klg"), out)
	// Out 2 like: `Total: 1h7m`
	assert.True(t, strings.Contains(out[2], "1h7m"), out)
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
	assert.True(t, strings.Contains(out[1], "9:00 - 10:00"), out)
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
		[]string{"total", "--period", "foo", "test.klg"},
	)
	assert.True(t, strings.Contains(out[0], "2h"), out)
	assert.True(t, strings.Contains(out[1], "1h"), out)
	assert.True(t, strings.Contains(out[2], "`foo` is not a valid period"), out)
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
