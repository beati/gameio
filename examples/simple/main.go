package main

import (
	"fmt"
	"log"

	"github.com/beati/sdl"
)

func main() {
	err := sdl.Run(run)
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	err := sdl.Init(sdl.InitVideo)
	if err != nil {
		return err
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("SDL simple example", sdl.WindowPosUndefined,
		sdl.WindowPosUndefined, 800, 600, 0)
	if err != nil {
		return err
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RendererPresentVSync)
	if err != nil {
		return err
	}
	defer renderer.Destroy()

	for {
		sdl.HandleEvents()
		if !sdl.Running {
			break
		}

		if sdl.KeyPressed(sdl.Key_A) {
			fmt.Println("A pressed")
		}
		if sdl.KeyHeld(sdl.Key_B) {
			fmt.Println("B held")
		}
		if sdl.KeyReleased(sdl.Key_C) {
			fmt.Println("C released")
		}

		err = renderer.Clear()
		if err != nil {
			return err
		}

		renderer.Present()
	}

	return nil
}
