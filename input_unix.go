// +build !windows

package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
// #include <X11/Xutil.h>
import "C"

import(
	"fmt"
)

func doMainLoop() {
	//x := -1
	//y := -1
	var event C.XEvent
	//var button int
	display := C.XOpenDisplay(nil)
    if display == nil {
        panic("Cannot connect to X server!\n")
    }
    defer C.XCloseDisplay(display)

	root := C.XDefaultRootWindow(display)
	C.XGrabPointer(display, root, C.False, C.ButtonPressMask, C.GrabModeAsync, C.GrabModeAsync, C.None, C.None, C.CurrentTime)
	for {
		C.XSelectInput(display, root, C.ButtonReleaseMask)
		for {
			C.XNextEvent(display, &event)
			switch event[0] {
				case C.ButtonPress:
					for i := range event {
						fmt.Println(i, event[i])
					}
			}

			//if x >= 0 && y >= 0 {
			//break
			//}
			//if button == C.Button1 {
			//fmt.Printf("leftclick at %d %d \n", x, y)
			//}
			//} else {
			//fmt.Printf("rightclick at %d %d \n", x, y)
			//}
		}
	}
}
