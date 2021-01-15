package app

import (
	"os"
)

func WithService(fn func(Service)) {
	path := "./tmp.klg"
	_ = os.Remove(path)
	file, err := os.Create(path)
	if err != nil {
		panic("Could not create context")
	}
	defer file.Close()
	service := &context{}
	fn(service)
	_ = os.Remove(path)
}
