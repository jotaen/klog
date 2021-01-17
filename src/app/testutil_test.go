package app

import (
	"os"
)

func WithContext(fn func(Context)) {
	path := "./tmp.klg"
	_ = os.Remove(path)
	file, err := os.Create(path)
	if err != nil {
		panic("Could not create context")
	}
	defer file.Close()
	ctx := Context{}
	fn(ctx)
	_ = os.Remove(path)
}
