// +build windows

package main

import(
	"bytes"
	"errors"
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

func NewInput() *InputWindows {
	return &InputWindows{}
}

type InputWindows struct {
	getCursorPos *syscall.Proc
	getDesktopWindow *syscall.Proc
	getWindowRect *syscall.Proc
	hwnd HWND
}

func (i *InputWindows) Init() error {
	user32, err := syscall.LoadDLL("user32")
	if err != nil {
		return err
	}
	defer user32.Release()
	i.getCursorPos, err = user32.FindProc("GetCursorPos")
	if err != nil {
		return err
	}
	i.getDesktopWindow, err = user32.FindProc("GetDesktopWindow")
	if err != nil {
		return err
	}
	i.getWindowRect, err = user32.FindProc("GetWindowRect")
	if err != nil {
		return err
	}

	res, _, _ := i.getDesktopWindow.Call()
	i.hwnd = HWND(res)
	if i.hwnd == 0 {
		return errors.New("Recieved value of '0' when requesting HWND")
	}

	return nil
}

func (i *InputWindows) Cleanup() {

}

func (i *InputWindows) GetWindowSize() (int, int) {
	rect := &RECT{}
	i.getWindowRect.Call(uintptr(i.hwnd), uintptr(unsafe.Pointer(rect)))

	w, h := rect.Right, rect.Bottom
	return int(w), int(h)
}

func (i *InputWindows) GetMouse() (int, int) {
	pt := &POINT{}
	i.getCursorPos.Call(uintptr(unsafe.Pointer(pt)))
	return int(pt.X), int(pt.Y)
}

func thisCouldProbablyBeAComment() {
	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	// has better response than RegisterHotKey, but is blocking
	//getmsgw := user32.MustFindProc("GetMessageW")
	peekmsg := user32.MustFindProc("PeekMessageW")
	reghotkey := user32.MustFindProc("RegisterHotKey")

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

		peekmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

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

// even more legacy junk
func doMainLoop() {
	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

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

	for {
		var pt = &POINT{}

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

		time.Sleep(time.Millisecond * 10)
	}
}

// end of windows
