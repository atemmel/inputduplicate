package main

import(
	"bytes"
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

// these are for windows
const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
	ModWin
)

type Hotkey struct {
	Id int			// unique identifier
	Modifiers int	// mask of modifiers
	KeyCode int		// key code, 'A', etc
}

type HWND uintptr

type POINT struct {
	X int32
	Y int32
}

type RECT struct {
	Left int32
	Top int32
	Right int32
	Bottom int32
}

type MSG struct {
	Hwnd HWND
	Uint uintptr
	WParam int16
	LParam int64
	DWord int32
	Point POINT
}

// end of windows

func (h *Hotkey) String() string {
	mod := &bytes.Buffer{}
	if h.Modifiers & ModAlt != 0 {
		mod.WriteString("Alt+")
	}
	if h.Modifiers & ModCtrl != 0 {
		mod.WriteString("Ctrl+")
	}
	if h.Modifiers & ModShift != 0 {
		mod.WriteString("Shift+")
	}
	if h.Modifiers & ModWin != 0 {
		mod.WriteString("Win+")
	}
	return fmt.Sprintf("Hotkey[Id: %d, %s%c]", h.Id, mod, h.KeyCode)
}

func main() {
	//TODO: look into replacing all MustXY functions with "proper" error reporting 
	//TODO: look into platform specific code (win32 vs x)

	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	// has better response than RegisterHotKey, but is blocking
	//getmsgw := user32.MustFindProc("GetMessageW")
	peekmsg := user32.MustFindProc("PeekMessageW")
	reghotkey := user32.MustFindProc("RegisterHotKey")
	getcursorpos := user32.MustFindProc("GetCursorPos")
	getdesktopwindow := user32.MustFindProc("GetDesktopWindow")
	getwindowrect := user32.MustFindProc("GetWindowRect")

	res, _, _ := getdesktopwindow.Call()
	hwnd := HWND(res)
	if hwnd == 0 {
		fmt.Println("HWND was 0 :(")
		return
	}

	rect := &RECT{}
	getwindowrect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(rect)))

	w, h := rect.Right, rect.Bottom
	fmt.Printf("Display is: %dx%d\n", w, h)

	/*
	// used to block all input devices
	blockinput := user32.MustFindProc("BlockInput")
	blockinput.Call(true)
	*/

	keys := map[int16]*Hotkey{
		1: {1, ModAlt + ModCtrl,  'O'},
		2: {2, ModAlt + ModShift, 'M'},
		3: {3, ModAlt + ModCtrl,  'X'},
	}

	for _, v := range keys {
		r1, _, err := reghotkey.Call(
			0, uintptr(v.Id), uintptr(v.Modifiers), uintptr(v.KeyCode))
		if r1 == 1 {
			fmt.Println("Registered", v)
		} else {
			fmt.Println("Failed to register", v, ", error:", err)
		}
	}

	for {
		var msg = &MSG{}
		var pt = &POINT{}

		peekmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)
		getcursorpos.Call(uintptr(unsafe.Pointer(pt)))

		//fmt.Println("Mouse coords:", pt)

		if pt.X == 0 {
			fmt.Println("Balls are touching (LEFT)")
		} else if pt.X == w - 1 {
			fmt.Println("Balls are touching (RIGHT)")
		}

		if pt.Y == 0 {
			fmt.Println("Balls are touching (TOP)")
		} else if pt.Y == h - 1 {
			fmt.Println("Balls are touching (BOTTOM)")
		}

		// id is in WPARAM field
		if id := msg.WParam; id != 0 {
			fmt.Println("Hotkey pressed:", keys[id])
			if id == 3 { // ctrl+alt+x
				fmt.Println("CTRL+ALT+X pressed, exiting...")
				return
			}
		}

		time.Sleep(time.Millisecond * 10)
	}
}
