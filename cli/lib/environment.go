package lib

import (
	"klog/store"
)

type Environment struct {
	WorkDir string
	Store store.Store
}
