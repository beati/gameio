package sdl

/*
#cgo LDFLAGS: -lSDL2
#include <SDL2/SDL.h>
#include <stdlib.h>

SDL_Event event = {0};

SDL_Surface *
SDL_LoadBMPWrapper(const char *file) {
	return SDL_LoadBMP(file);
}

int filter(void *userData, SDL_Event *event) {
	if (event->type == SDL_QUIT) {
		return 1;
	}
	return 0;
}

void SetEventFilterWrapper(void) {
	SDL_SetEventFilter(filter, NULL);
}

int QuitRequestedWrapper(void) {
	return SDL_QuitRequested() == SDL_TRUE;
}
*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

func init() {
	runtime.LockOSThread()
}

var mainThreadFunc = make(chan func())

func mainThreadCall(f func()) {
	done := make(chan bool, 1)
	mainThreadFunc <- func() {
		f()
		done <- true
	}
	<-done
}

func Run(run func() error) error {
	done := make(chan error, 1)
	go func() {
		done <- run()
	}()

	var err error

mainThreadCallLoop:
	for {
		select {
		case f := <-mainThreadFunc:
			f()
		case err = <-done:
			break mainThreadCallLoop
		}
	}

	return err
}

func getError() error {
	var err *C.char
	mainThreadCall(func() {
		err = C.SDL_GetError()
	})
	return errors.New(C.GoString(err))
}

const (
	InitTimer          = C.SDL_INIT_TIMER
	InitAudio          = C.SDL_INIT_AUDIO
	InitVideo          = C.SDL_INIT_VIDEO
	InitJoystick       = C.SDL_INIT_JOYSTICK
	InitHaptic         = C.SDL_INIT_HAPTIC
	InitGameController = C.SDL_INIT_GAMECONTROLLER
	InitEvents         = C.SDL_INIT_EVENTS
	InitEverything     = C.SDL_INIT_EVERYTHING
	InitNoParachute    = C.SDL_INIT_NOPARACHUTE
)

func Init(flags uint32) error {
	var err C.int
	mainThreadCall(func() {
		err = C.SDL_Init(C.Uint32(flags))
	})
	if err != 0 {
		return getError()
	}

	mainThreadCall(func() {
		C.SetEventFilterWrapper()
	})

	return nil
}

func Quit() {
	mainThreadCall(func() {
		C.SDL_Quit()
	})
}

type Window C.SDL_Window

const (
	WindowPosCentered  = C.SDL_WINDOWPOS_CENTERED
	WindowPosUndefined = C.SDL_WINDOWPOS_UNDEFINED
)

const (
	WindowFullscreen        = C.SDL_WINDOW_FULLSCREEN
	WindowFullscreenDesktop = C.SDL_WINDOW_FULLSCREEN_DESKTOP
	WindowOpenGL            = C.SDL_WINDOW_OPENGL
	WindowHidden            = C.SDL_WINDOW_HIDDEN
	WindowBorderless        = C.SDL_WINDOW_BORDERLESS
	WindowResizable         = C.SDL_WINDOW_RESIZABLE
	WindowMinimized         = C.SDL_WINDOW_MINIMIZED
	WindowMaximized         = C.SDL_WINDOW_MAXIMIZED
	WindowInputGrabbed      = C.SDL_WINDOW_INPUT_GRABBED
	WindowAllowHighDPI      = C.SDL_WINDOW_ALLOW_HIGHDPI
)

func CreateWindow(title string, x, y, w, h int, flags uint32) (*Window, error) {
	t := C.CString(title)
	defer C.free(unsafe.Pointer(t))

	var win *C.SDL_Window
	mainThreadCall(func() {
		win = C.SDL_CreateWindow(t, C.int(x), C.int(y), C.int(w),
			C.int(h), C.Uint32(flags))
	})
	if win == nil {
		return nil, getError()
	}
	return (*Window)(win), nil
}

func (w *Window) Destroy() {
	mainThreadCall(func() {
		C.SDL_DestroyWindow((*C.SDL_Window)(w))
	})
}

type Renderer C.SDL_Renderer

const (
	RendererSoftware      = C.SDL_RENDERER_SOFTWARE
	RendererAccelerated   = C.SDL_RENDERER_ACCELERATED
	RendererPresentVSync  = C.SDL_RENDERER_PRESENTVSYNC
	RendererTargetTexture = C.SDL_RENDERER_TARGETTEXTURE
)

func CreateRenderer(w *Window, index int, flags uint32) (*Renderer, error) {
	var r *C.SDL_Renderer
	mainThreadCall(func() {
		r = C.SDL_CreateRenderer((*C.SDL_Window)(w), C.int(index),
			C.Uint32(flags))
	})
	if r == nil {
		return nil, getError()
	}
	return (*Renderer)(r), nil
}

func (r *Renderer) Destroy() {
	mainThreadCall(func() {
		C.SDL_DestroyRenderer((*C.SDL_Renderer)(r))
	})
}

func (r *Renderer) Clear() error {
	var err C.int
	mainThreadCall(func() {
		err = C.SDL_RenderClear((*C.SDL_Renderer)(r))
	})
	if err != 0 {
		return getError()
	}
	return nil
}

type Point struct {
	X int
	Y int
}

type Rect struct {
	X int
	Y int
	W int
	H int
}

const (
	FlipNone       = C.SDL_FLIP_NONE
	FlipHorizontal = C.SDL_FLIP_HORIZONTAL
	FlipVertical   = C.SDL_FLIP_VERTICAL
)

func (r *Renderer) CopyEx(t *Texture, src, dst *Rect, angle float64,
	center *Point, flip int) error {
	var csrc *C.SDL_Rect
	var cdst *C.SDL_Rect
	var ccenter *C.SDL_Point
	if src != nil {
		csrc = &C.SDL_Rect{C.int(src.X), C.int(src.Y), C.int(src.W),
			C.int(src.H)}
	}
	if dst != nil {
		cdst = &C.SDL_Rect{C.int(dst.X), C.int(dst.Y), C.int(dst.W),
			C.int(dst.H)}
	}
	if center != nil {
		ccenter = &C.SDL_Point{C.int(center.X), C.int(center.Y)}
	}
	var err C.int
	mainThreadCall(func() {
		err = C.SDL_RenderCopyEx((*C.SDL_Renderer)(r),
			(*C.SDL_Texture)(t), csrc, cdst, C.double(angle),
			ccenter, C.SDL_RendererFlip(flip))
	})
	if err != 0 {
		return getError()
	}
	return nil
}

func (r *Renderer) Present() {
	mainThreadCall(func() {
		C.SDL_RenderPresent((*C.SDL_Renderer)(r))
	})
}

type Texture C.SDL_Texture

func LoadBMP(r *Renderer, file string) (*Texture, error) {
	f := C.CString(file)
	defer C.free(unsafe.Pointer(f))

	var s *C.SDL_Surface
	mainThreadCall(func() {
		s = C.SDL_LoadBMPWrapper(f)
	})
	if s == nil {
		return nil, getError()
	}
	defer mainThreadCall(func() {
		C.SDL_FreeSurface(s)
	})

	var err C.int
	mainThreadCall(func() {
		color := C.SDL_MapRGB(s.format, 255, 0, 255)
		err = C.SDL_SetColorKey(s, C.SDL_TRUE, color)
	})
	if err != 0 {
		return nil, getError()
	}

	var t *C.SDL_Texture
	mainThreadCall(func() {
		t = C.SDL_CreateTextureFromSurface((*C.SDL_Renderer)(r), s)
	})
	if t == nil {
		return nil, getError()
	}

	return (*Texture)(t), nil
}

func (t *Texture) Destroy() {
	mainThreadCall(func() {
		C.SDL_DestroyTexture((*C.SDL_Texture)(t))
	})
}

var Running = true

func PumpEvents() {
	mainThreadCall(func() {
		C.SDL_PumpEvents()
	})

	Running = C.QuitRequestedWrapper() == 0
}

func HandleEvents() {
	event := &C.event
	for {
		var noEvent bool
		mainThreadCall(func() {
			noEvent = int(C.SDL_PollEvent(event)) == 0
		})
		if noEvent {
			break
		}

		etype := *(*C.Uint32)(unsafe.Pointer(event))
		switch etype {
		case C.SDL_QUIT:
			Running = false
		case C.SDL_KEYDOWN:
			fallthrough
		case C.SDL_KEYUP:
			handleKeyboardEvent(event)
		case C.SDL_MOUSEBUTTONDOWN:
			fallthrough
		case C.SDL_MOUSEBUTTONUP:
			handleMouseButtonEvent(event)
		}
	}
}

func isPressed(state C.Uint8) bool {
	return state == C.SDL_PRESSED
}

type MouseState struct {
	X     int
	Y     int
	Left  bool
	Right bool
}

var Mouse MouseState

func handleMouseButtonEvent(event *C.SDL_Event) {
	button := (*C.SDL_MouseButtonEvent)(unsafe.Pointer(event))
	switch button.button {
	case C.SDL_BUTTON_LEFT:
		Mouse.Left = isPressed(button.state)
	case C.SDL_BUTTON_RIGHT:
		Mouse.Right = isPressed(button.state)
	}
	Mouse.X = int(button.x)
	Mouse.Y = int(button.y)
}

type KeyboardState struct {
	Up    bool
	Down  bool
	Left  bool
	Right bool
	X     bool
	Z     bool
}

var Keyboard KeyboardState

func handleKeyboardEvent(event *C.SDL_Event) {
	key := (*C.SDL_KeyboardEvent)(unsafe.Pointer(&event))
	switch key.keysym.sym {
	case C.SDLK_UP:
		Keyboard.Up = isPressed(key.state)
	case C.SDLK_DOWN:
		Keyboard.Down = isPressed(key.state)
	case C.SDLK_LEFT:
		Keyboard.Left = isPressed(key.state)
	case C.SDLK_RIGHT:
		Keyboard.Right = isPressed(key.state)
	case C.SDLK_x:
		Keyboard.X = isPressed(key.state)
	case C.SDLK_z:
		Keyboard.Z = isPressed(key.state)
	}
}
