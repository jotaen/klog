package withcontext

import (
	"klog/app"
	"os"
)

func WithService(fn func(app.Service)) {
	path := "./tmp.klg"
	_ = os.Remove(path)
	file, err := os.Create(path)
	if err != nil {
		panic("Could not create context")
	}
	defer file.Close()
	service := app.NewService(file)
	fn(service)
	_ = os.Remove(path)
}
