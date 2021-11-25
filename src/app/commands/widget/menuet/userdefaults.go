//go:build darwin

package menuet

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

#ifndef __USERDEFAULTS_H_H__
#import "userdefaults.h"
#endif

*/
import "C"
import (
	"encoding/json"
	"sync"
	"unsafe"
)

// UserDefaults represents values stored in NSUserDefaults
type UserDefaults struct {
	strings map[string]string
	ints    map[string]int
	bools   map[string]bool
}

var defaultsInstance *UserDefaults
var defaultsOnce sync.Once

// Defaults returns the userDefaults singleton
func Defaults() *UserDefaults {
	defaultsOnce.Do(func() {
		defaultsInstance = &UserDefaults{
			strings: make(map[string]string),
			ints:    make(map[string]int),
			bools:   make(map[string]bool),
		}
	})
	return defaultsInstance
}

// SetString sets a string default
func (u *UserDefaults) SetString(key, value string) {
	ckey := C.CString(string(key))
	cvalue := C.CString(string(value))
	C.setString(ckey, cvalue)
	C.free(unsafe.Pointer(ckey))
	C.free(unsafe.Pointer(cvalue))
	u.strings[key] = value
}

// String gets a string default, "" if not set
func (u *UserDefaults) String(key string) string {
	val, ok := u.strings[key]
	if ok {
		return val
	}
	ckey := C.CString(string(key))
	cvalue := C.getString(ckey)
	value := C.GoString(cvalue)
	C.free(unsafe.Pointer(ckey))
	return value
}

// SetInteger sets an integer default
func (u *UserDefaults) SetInteger(key string, value int) {
	ckey := C.CString(string(key))
	C.setInteger(ckey, C.long(value))
	C.free(unsafe.Pointer(ckey))
	u.ints[key] = value
}

// Integer gets an integer default, 0 if not set
func (u *UserDefaults) Integer(key string) int {
	val, ok := u.ints[key]
	if ok {
		return val
	}
	ckey := C.CString(string(key))
	value := C.getInteger(ckey)
	C.free(unsafe.Pointer(ckey))
	return int(value)
}

// SetBoolean sets a boolean default
func (u *UserDefaults) SetBoolean(key string, value bool) {
	ckey := C.CString(string(key))
	C.setBoolean(ckey, C._Bool(value))
	C.free(unsafe.Pointer(ckey))
	u.bools[key] = value
}

// Boolean gets a boolean default, 0 if not set
func (u *UserDefaults) Boolean(key string) bool {
	val, ok := u.bools[key]
	if ok {
		return val
	}
	ckey := C.CString(string(key))
	value := C.getBoolean(ckey)
	C.free(unsafe.Pointer(ckey))
	return bool(value)
}

// Marshal marshals an object into JSON and stores it in user defaults, see json.Marshal docs
func (u *UserDefaults) Marshal(key string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	u.SetString(key, string(b))
	return nil
}

// Unmarshal unmarshals an object from JSON that was stored in user defaults, see json.Unmarshal docs
func (u *UserDefaults) Unmarshal(key string, v interface{}) error {
	b := u.String(key)
	return json.Unmarshal([]byte(b), v)
}
