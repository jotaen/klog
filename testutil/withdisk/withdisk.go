package withdisk

import (
	"os"
)

func WithDisk(fn func(string)) {
	path := "./tmp/test"
	os.RemoveAll(path)
	os.MkdirAll(path, os.ModePerm)
	fn(path)
	os.RemoveAll(path)
	os.Remove("./tmp")
}
