package klog

import (
	"github.com/jotaen/klog/klog/app"
	tf "github.com/jotaen/klog/lib/terminalformat"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"testing"
)

type Env struct {
	files map[string]string
}

type invocation struct {
	args []string
	test func(t *testing.T, code int, out string)
}

func (e *Env) execute(t *testing.T, is ...invocation) {
	// Create temp directory and change work dir to it.
	tmpDir, tErr := os.MkdirTemp("", "")
	assertNil(tErr)
	cErr := os.Chdir(tmpDir)
	assertNil(cErr)

	// Write out all files from `Env`.
	for name, contents := range e.files {
		err := os.WriteFile(name, []byte(contents), 0644)
		assertNil(err)
	}

	// Capture “old” stdout, so that we can restore later.
	oldStdout := os.Stdout

	// Run all commands one after the other.
	for _, invoke := range is {
		r, w, _ := os.Pipe()
		os.Stdout = w

		config := app.NewDefaultConfig(tf.COLOUR_THEME_NO_COLOUR)
		code, runErr := Run(app.NewFileOrPanic(tmpDir), app.Meta{
			Specification: "[Specification text]",
			License:       "[License text]",
			Version:       "v0.0",
			SrcHash:       "abc1234",
		}, config, invoke.args)

		_ = w.Close()

		t.Run(strings.Join(invoke.args, "__"), func(t *testing.T) {
			if runErr != nil {
				require.NotEqual(t, 0, code, "App returned error, but exit code was 0")
			} else {
				out, _ := io.ReadAll(r)
				invoke.test(t, code, tf.StripAllAnsiSequences(string(out)))
			}
		})
	}

	// Clean up temp dir.
	rErr := os.RemoveAll(tmpDir)
	assertNil(rErr)
	os.Stdout = oldStdout
}

func assertNil(e error) {
	if e != nil {
		panic(e)
	}
}
