//go:build darwin
// +build darwin

package menuet

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

#ifndef __ALERT_H_H__
#import "alert.h"
#endif

void showAlert(const char *jsonString);

*/
import "C"
import (
	"encoding/json"
	"log"
	"unsafe"
)

// Alert represents an NSAlert
type Alert struct {
	MessageText     string
	InformativeText string
	Buttons         []string
	Inputs          []string
}

// AlertClicked represents a selected alert button
type AlertClicked struct {
	Button int
	Inputs []string
}

// Alert shows an alert, and returns the index of the button pressed, or -1 if none
func (a *Application) Alert(alert Alert) AlertClicked {
	if a.alertChannel != nil {
		log.Printf("Alert message already showing")
		return AlertClicked{-1, nil}
	}
	b, err := json.Marshal(alert)
	if err != nil {
		log.Printf("Marshal: %v", err)
		return AlertClicked{-1, nil}
	}
	cstr := C.CString(string(b))
	C.showAlert(cstr)
	C.free(unsafe.Pointer(cstr))
	a.alertChannel = make(chan AlertClicked)
	response := <-a.alertChannel
	a.alertChannel = nil
	return response
}

//export alertClicked
func alertClicked(button int, valuesCString *C.char) {
	valuesJSON := C.GoString(valuesCString)
	values := make([]string, 0)
	err := json.Unmarshal([]byte(valuesJSON), &values)
	if err != nil {
		log.Printf("Unmarshal: %v", err)
		return
	}
	app := App()
	if app.alertChannel == nil {
		log.Printf("Alert message double clicked?")
		return
	}
	app.alertChannel <- AlertClicked{
		Button: button,
		Inputs: values,
	}
}
