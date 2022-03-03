package lib

import (
	"github.com/jotaen/klog/src/app"
	"os"
	"os/signal"
	"syscall"
	gotime "time"
)

// WithRepeat repetitively invokes the callback at a rate of 1/s.
// It always clears the terminal screen.
func WithRepeat(ctx app.Context, fn func(int64) error) error {
	// Handle ^C gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	// Call handler function repetitively
	ctx.Print("\033[2J") // Initial screen clearing
	ticker := gotime.NewTicker(1 * gotime.Second)
	defer ticker.Stop()
	secondsCounter := int64(0) // Choose large type because of overflow
	for ; true; <-ticker.C {
		secondsCounter += 1
		ctx.Print("\033[H\033[J") // Cursor reset
		err := fn(secondsCounter)
		if err != nil {
			return err
		}
	}
	return nil
}
