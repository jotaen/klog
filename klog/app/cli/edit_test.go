package cli

import (
	"strings"
	"testing"

	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/args"
	"github.com/jotaen/klog/klog/app/cli/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditWithAutoEditor(t *testing.T) {
	spy := newCommandSpy(nil)
	state, err := NewTestingContext()._SetEditors([]command.Command{
		command.New("editor", []string{"--file"}),
	}, "")._SetExecute(spy.Execute)._Run((&Edit{
		OutputFileArgs: args.OutputFileArgs{File: "/tmp/file.klg"},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, 1, spy.Count)
	assert.Equal(t, "editor", spy.LastCmd.Bin)
	assert.Equal(t, []string{"--file", "/tmp/file.klg"}, spy.LastCmd.Args)
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
		OutputFileArgs: args.OutputFileArgs{File: "/tmp/file.klg"},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, 2, spy.Count)
	assert.Equal(t, "editor2", spy.LastCmd.Bin)
	assert.Equal(t, []string{"--file", "/tmp/file.klg"}, spy.LastCmd.Args)
}

func TestEditWithExplicitEditor(t *testing.T) {
	spy := newCommandSpy(nil)
	state, err := NewTestingContext()._SetEditors([]command.Command{
		command.New("editor", []string{"--file"}),
	}, "myedit")._SetExecute(spy.Execute)._Run((&Edit{
		OutputFileArgs: args.OutputFileArgs{File: "/tmp/file.klg"},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, 1, spy.Count)
	assert.Equal(t, "myedit", spy.LastCmd.Bin)
	assert.Equal(t, []string{"/tmp/file.klg"}, spy.LastCmd.Args)
	// No hint was printed:
	assert.Equal(t, "", state.printBuffer)
}

func TestEditWithExplicitEditorWithSpaces(t *testing.T) {
	spy := newCommandSpy(nil)
	_, err := NewTestingContext()._SetEditors([]command.Command{
		command.New("editor", []string{"--file"}),
	}, "'C:\\Program Files\\Sublime Text'")._SetExecute(spy.Execute)._Run((&Edit{
		OutputFileArgs: args.OutputFileArgs{File: "/tmp/file.klg"},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, 1, spy.Count)
	assert.Equal(t, "C:\\Program Files\\Sublime Text", spy.LastCmd.Bin)
	assert.Equal(t, []string{"/tmp/file.klg"}, spy.LastCmd.Args)
}

func TestEditWithExplicitEditorWithAdditionalArgs(t *testing.T) {
	spy := newCommandSpy(nil)
	_, err := NewTestingContext()._SetEditors([]command.Command{
		command.New("editor", []string{"--file"}),
	}, "myedit -f")._SetExecute(spy.Execute)._Run((&Edit{
		OutputFileArgs: args.OutputFileArgs{File: "/tmp/file.klg"},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, 1, spy.Count)
	assert.Equal(t, "myedit", spy.LastCmd.Bin)
	assert.Equal(t, []string{"-f", "/tmp/file.klg"}, spy.LastCmd.Args)
}

func TestEditFailsWithExplicitEditorThatHasMalformedSyntax(t *testing.T) {
	for _, editor := range []string{
		// Unmatched single quote:
		`'myedit`,

		// Unmatched double quote
		`myedit --arg "foo`,
	} {
		spy := newCommandSpy(nil)
		_, err := NewTestingContext()._SetEditors([]command.Command{
			command.New("editor", []string{"--file"}),
		}, editor)._SetExecute(spy.Execute)._Run((&Edit{
			OutputFileArgs: args.OutputFileArgs{File: "/tmp/file.klg"},
		}).Run)
		require.Error(t, err)
		assert.Equal(t, "Invalid editor setting", err.Error())
	}
}

func TestFailsIfExplicitEditorCrashes(t *testing.T) {
	spy := newCommandSpy(func(c command.Command) app.Error {
		return app.NewError("Error", "Error", nil)
	})
	_, err := NewTestingContext()._SetEditors([]command.Command{
		command.New("editor", []string{"--file"}),
	}, "myedit")._SetExecute(spy.Execute)._Run((&Edit{
		OutputFileArgs: args.OutputFileArgs{File: "/tmp/file.klg"},
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
		OutputFileArgs: args.OutputFileArgs{File: "/tmp/file.klg"},
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
