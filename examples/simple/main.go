package main

import (
	"log"

	"assets"

	"github.com/beati/sdl"
)

func main() {
	err = sdl.Run(run)
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	err := sdl.Init(sdl.InitVideo | sdl.InitGameController)
	if err != nil {
		return err
	}
	defer sdl.Quit()

	windowW := 800
	windowH := 600
	window, err := sdl.CreateWindow("SDL simple example", sdl.WindowPosUndefined,
		sdl.WindowPosUndefined, windowW, windowH, 0)
	if err != nil {
		return err
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RendererPresentVSync)
	//renderer, err := sdl.CreateRenderer(window, -1, 0)
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

	x := 100.0
	y := 100.0

	for {
		sdl.HandleEvents()
		if !sdl.Running {
			break
		}

		dt := sdl.ClockElapsed()
		//fmt.Println(dt)

		const speed = 400.0
		if sdl.KeyHeld(sdl.Key_UP) || sdl.ButtonHeld(sdl.Button_UP) {
			y -= speed * dt.Seconds()
		} else if sdl.KeyHeld(sdl.Key_DOWN) || sdl.ButtonHeld(sdl.Button_DOWN) {
			y += speed * dt.Seconds()
		}
		if sdl.KeyHeld(sdl.Key_LEFT) || sdl.ButtonHeld(sdl.Button_LEFT) {
			x -= speed * dt.Seconds()
		} else if sdl.KeyHeld(sdl.Key_RIGHT) || sdl.ButtonHeld(sdl.Button_RIGHT) {
			x += speed * dt.Seconds()
		}

		if y < 0.0 {
			y = 0.0
		}
		if y > float64(windowH-assets.Img3.H) {
			y = float64(windowH - assets.Img3.H)
		}
		if x < 0.0 {
			x = 0.0
		}
		if x > float64(windowW-assets.Img3.W) {
			x = float64(windowW - assets.Img3.W)
		}

		err = renderer.Clear()
		if err != nil {
			return err
		}

		src := sdl.Rect{0, 0, assets.Img3.W, assets.Img3.H}
		dst := sdl.Rect{int(x), int(y), assets.Img3.W, assets.Img3.H}
		renderer.CopyEx(t, &src, &dst, 0, nil, sdl.FlipNone)

		renderer.Present()
	}

	return nil
}
