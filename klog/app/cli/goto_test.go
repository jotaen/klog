package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/app/cli/lib/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGotoLocation(t *testing.T) {
	spy := newCommandSpy(nil)
	_, err := NewTestingContext()._SetFileExplorers([]command.Command{
		command.New("goto", []string{"--file"}),
	})._SetExecute(spy.Execute)._Run((&Goto{
		OutputFileArgs: lib.OutputFileArgs{File: "/tmp/file.klg"},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, 1, spy.Count)
	assert.Equal(t, "goto", spy.LastCmd.Bin)
	assert.Equal(t, []string{"--file", "/tmp"}, spy.LastCmd.Args)
}

func TestGotoLocationFirstSucceeds(t *testing.T) {
	spy := newCommandSpy(func(c command.Command) app.Error {
		if c.Bin == "goto2" {
			return nil
		}
		return app.NewError("Error", "Error", nil)
	})
	_, err := NewTestingContext()._SetFileExplorers([]command.Command{
		command.New("goto1", []string{"--file"}),
		command.New("goto2", nil),
		command.New("goto3", []string{"--file"}),
	})._SetExecute(spy.Execute)._Run((&Goto{
		OutputFileArgs: lib.OutputFileArgs{File: "/tmp/file.klg"},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, 2, spy.Count)
	assert.Equal(t, "goto2", spy.LastCmd.Bin)
	assert.Equal(t, []string{"/tmp"}, spy.LastCmd.Args)
}

func TestGotoFails(t *testing.T) {
	spy := newCommandSpy(func(_ command.Command) app.Error {
		return app.NewError("Error", "Error", nil)
	})
	_, err := NewTestingContext()._SetFileExplorers([]command.Command{
		command.New("goto", []string{"--file"}),
	})._SetExecute(spy.Execute)._Run((&Goto{
		OutputFileArgs: lib.OutputFileArgs{File: "/tmp/file.klg"},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, 1, spy.Count)
	assert.Equal(t, "Failed to open file browser", err.Error())
}

func TestGotoFailsWithoutFile(t *testing.T) {
	spy := newCommandSpy(nil)
	_, err := NewTestingContext()._SetEditors([]command.Command{
		command.New("goto", []string{"--file"}),
	}, "")._SetExecute(spy.Execute)._Run((&Goto{}).Run)
	require.Error(t, err)
	assert.Equal(t, 0, spy.Count)
}
