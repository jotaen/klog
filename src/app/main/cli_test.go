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
