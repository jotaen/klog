package klog

import (
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib/terminalformat"
	"io"
	"io/ioutil"
	"os"
)

type Env struct {
	files map[string]string
}

func (e *Env) run(invocation ...[]string) []string {
	// Create temp directory and change work dir to it.
	tmpDir, tErr := ioutil.TempDir("", "")
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
	outs := make([]string, len(invocation))
	for i, args := range invocation {
		r, w, _ := os.Pipe()
		os.Stdout = w

		runErr, code := Run(tmpDir, app.Meta{
			Specification: "[Specification text]",
			License:       "[License text]",
			Version:       "v0.0",
			BuildHash:     "1234567",
		}, false, args)

		_ = w.Close()
		if runErr != nil || code != 0 {
			outs[i] = runErr.Error()
			continue
		}
		out, _ := io.ReadAll(r)
		outs[i] = terminalformat.StripAllAnsiSequences(string(out))
	}

	// Clean up temp dir.
	rErr := os.RemoveAll(tmpDir)
	assertNil(rErr)
	os.Stdout = oldStdout

	return outs
}

func assertNil(e error) {
	if e != nil {
		panic(e)
	}
}
