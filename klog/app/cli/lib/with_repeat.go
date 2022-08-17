package lib

import (
	"os"
	"os/signal"
	"syscall"
	gotime "time"
)

// WithRepeat repetitively invokes the callback at the desired rate.
// It always clears the terminal screen.
func WithRepeat(print func(string), interval gotime.Duration, fn func(int64) error) error {
	// Handle ^C gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	// Call handler function repetitively
	print("\033[2J") // Initial screen clearing
	ticker := gotime.NewTicker(interval)
	defer ticker.Stop()
	secondsCounter := int64(0) // Choose large type because of overflow
	for ; true; <-ticker.C {
		secondsCounter += 1
		print("\033[H\033[J") // Cursor reset
		err := fn(secondsCounter)
		if err != nil {
			return err
		}
	}
	return nil
}
