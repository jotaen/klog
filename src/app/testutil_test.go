package app

import (
	"os"
)

func WithContext(fn func(Context)) {
	path := "./tmp.klg"
	_ = os.Remove(path)
	file, err := os.Create(path)
	if err != nil {
		panic("Could not initialise test environment")
	}
	defer file.Close()
	ctx := Context{}
	fn(ctx)
	err = os.Remove(path)
	if err != nil {
		panic("Could clean up")
	}
}
