package main

import(
	"fmt"
	"time"
)

type Input interface {
	Init() error
	Cleanup()
	GetWindowSize() (int, int)
	GetMouse() (int, int)
}

func main() {
	//TODO: look into replacing all MustXY functions with "proper" error reporting 
	//TODO: look into platform specific code (win32 vs x)

	input := NewInput()

	err := input.Init()
	if err != nil {
		panic(err)
	}
	defer input.Cleanup()

	w, h := input.GetWindowSize()

	for {
		x, y := input.GetMouse()

		if x == 0 {
			fmt.Println("Balls are touching (LEFT)")
		} else if x == w - 1 {
			fmt.Println("Balls are touching (RIGHT)")
		}

		if y == 0 {
			fmt.Println("Balls are touching (TOP)")
		} else if y == h - 1 {
			fmt.Println("Balls are touching (BOTTOM)")
		}

		time.Sleep(time.Millisecond * 10)
	}
}
