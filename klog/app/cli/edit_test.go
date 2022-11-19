package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/app/cli/lib/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestEditWithAutoEditor(t *testing.T) {
	spy := newCommandSpy(nil)
	state, err := NewTestingContext()._SetEditors([]command.Command{
		command.New("editor", []string{"--file"}),
	}, "")._SetExecute(spy.Execute)._Run((&Edit{
		OutputFileArgs: lib.OutputFileArgs{File: "/tmp/file.klg"},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, 1, spy.Count)
	assert.Equal(t, "editor --file /tmp/file.klg", spy.LastCmd.ToString())
	// Hint was printed:
	assert.Equal(t, hint, strings.Trim(state.printBuffer, "\n"))
}

func TestFirstAutoEditorSucceeds(t *testing.T) {
	spy := newCommandSpy(func(c command.Command) app.Error {
		if c.Bin == "editor2" {
			return nil
		}
		return app.NewError("Error", "Error", nil)
	})
	_, err := NewTestingContext()._SetEditors([]command.Command{
		command.New("editor1", []string{"--file"}),
		command.New("editor2", []string{"--file"}),
		command.New("editor3", []string{"--file"}),
	}, "")._SetExecute(spy.Execute)._Run((&Edit{
		OutputFileArgs: lib.OutputFileArgs{File: "/tmp/file.klg"},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, 2, spy.Count)
	assert.Equal(t, "editor2 --file /tmp/file.klg", spy.LastCmd.ToString())
}

func TestEditWithExplicitEditor(t *testing.T) {
	spy := newCommandSpy(nil)
	state, err := NewTestingContext()._SetEditors([]command.Command{
		command.New("editor", []string{"--file"}),
	}, "myedit")._SetExecute(spy.Execute)._Run((&Edit{
		OutputFileArgs: lib.OutputFileArgs{File: "/tmp/file.klg"},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, 1, spy.Count)
	assert.Equal(t, "myedit /tmp/file.klg", spy.LastCmd.ToString())
	// No hint was printed:
	assert.Equal(t, "", state.printBuffer)
}

func TestFailsIfExplicitEditorIncorrect(t *testing.T) {
	spy := newCommandSpy(func(c command.Command) app.Error {
		return app.NewError("Error", "Error", nil)
	})
	_, err := NewTestingContext()._SetEditors([]command.Command{
		command.New("editor", []string{"--file"}),
	}, "myedit")._SetExecute(spy.Execute)._Run((&Edit{
		OutputFileArgs: lib.OutputFileArgs{File: "/tmp/file.klg"},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, "Cannot open preferred editor", err.Error())
}

func TestFailsIfAutoEditorsUnsuccessful(t *testing.T) {
	spy := newCommandSpy(func(c command.Command) app.Error {
		return app.NewError("Error", "Error", nil)
	})
	_, err := NewTestingContext()._SetEditors([]command.Command{
		command.New("editor1", []string{"--file"}),
		command.New("editor2", nil),
	}, "")._SetExecute(spy.Execute)._Run((&Edit{
		OutputFileArgs: lib.OutputFileArgs{File: "/tmp/file.klg"},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, 2, spy.Count)
	assert.Equal(t, "Cannot open any editor", err.Error())
}

func TestEditFailsWithoutFile(t *testing.T) {
	spy := newCommandSpy(nil)
	_, err := NewTestingContext()._SetEditors([]command.Command{
		command.New("editor", []string{"--file"}),
	}, "")._SetExecute(spy.Execute)._Run((&Edit{}).Run)
	require.Error(t, err)
	assert.Equal(t, 0, spy.Count)
}
