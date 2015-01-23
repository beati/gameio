package main

import (
	"fmt"
	"log"

	"assets"

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

	t, err := sdl.CreateTexture(renderer, assets.Img3.W, assets.Img3.H,
		assets.Img3.Pixels)
	if err != nil {
		return err
	}
	defer t.Destroy()

	sdl.ClockInit()

	for {
		sdl.HandleEvents()
		if !sdl.Running {
			break
		}

		dt := sdl.ClockElapsed()
		fmt.Println(dt)

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

		src := sdl.Rect{0, 0, assets.Img3.W, assets.Img3.H}
		dst := sdl.Rect{100, 100, assets.Img3.W, assets.Img3.H}
		renderer.CopyEx(t, &src, &dst, 0, nil, sdl.FlipNone)

		renderer.Present()
	}

	return nil
}
