package helper

import (
	"os"
	"os/signal"
	"syscall"
	gotime "time"

	"github.com/jotaen/klog/klog/app"
)

const (
	escSeqClearScreen = "\033[2J"
	escSeqCursorHome  = "\033[H"
	escSeqClearToEnd  = "\033[J"

	// Synchronized output (mode 2026): the terminal holds rendering between
	// begin and end, so each frame is displayed at once without flicker.
	// Terminals that don’t support it ignore these sequences, so they are still
	// prone to flickering.
	escSeqSyncBegin = "\033[?2026h"
	escSeqSyncEnd   = "\033[?2026l"
)

// WithRepeat repetitively invokes the callback at the desired rate.
// Before each invocation, it clears the terminal screen. Each frame is wrapped
// in synchronized output, so the callback’s output is displayed at once (thus avoiding
// flickering issues, provided that the terminal supports synchronized output).
func WithRepeat(print func(string), interval gotime.Duration, fn func(int64) app.Error) app.Error {
	// Handle ^C gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// The signal may arrive mid-frame, so we must make sure the terminal
		// isn’t left holding back its display after we exit.
		print(escSeqSyncEnd)
		os.Exit(0)
	}()

	// Call handler function repetitively
	print(escSeqClearScreen)
	ticker := gotime.NewTicker(interval)
	defer ticker.Stop()
	secondsCounter := int64(0) // Choose large type because of overflow
	for ; true; <-ticker.C {
		secondsCounter += 1
		print(escSeqSyncBegin + escSeqCursorHome + escSeqClearToEnd)
		err := fn(secondsCounter)
		print(escSeqSyncEnd)
		if err != nil {
			return err
		}
	}
	return nil
}
