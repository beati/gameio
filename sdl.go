package sdl

/*
#cgo LDFLAGS: -lSDL2
#cgo CFLAGS: -O3
#include <SDL2/SDL.h>
#include <stdlib.h>
#include <string.h>

int running = 1;

Uint8 keyboardPrev[SDL_NUM_SCANCODES];
int numScancodes;
const Uint8 *keyboardCurr;

void
handleEvents(void) {
	memcpy(keyboardPrev, keyboardCurr, numScancodes);

	SDL_Event event;
	while (SDL_PollEvent(&event)) {
		switch (event.type) {
		case SDL_QUIT:
			running = 0;
			break;
		}
	}
}

int
eventFilter(void *userData, SDL_Event *event) {
	(void)userData;

	switch (event->type) {
	case SDL_QUIT:
		return 1;
	}
	return 0;
}

int
init(Uint32 flags) {
	int err = SDL_Init(flags);
	if (err != 0) {
		return err;
	}
	SDL_SetEventFilter(eventFilter, NULL);
	keyboardCurr = SDL_GetKeyboardState(&numScancodes);
	return 0;
}

SDL_Texture *
CreateTexture(SDL_Renderer *r, int w, int h, Uint32 *pixels) {
	SDL_Texture *t = SDL_CreateTexture(r, SDL_PIXELFORMAT_RGBA8888,
		SDL_TEXTUREACCESS_STATIC, w, h);
	if (t == NULL) {
		goto fail1;
	}

	int err = 0;
	err = SDL_SetTextureBlendMode(t, SDL_BLENDMODE_BLEND);
	if (err != 0) {
		goto fail2;
	}

	err = SDL_UpdateTexture(t, NULL, pixels, w * sizeof(pixels[0]));
	if (err != 0) {
		goto fail2;
	}

	if (0) {
fail2:
		SDL_DestroyTexture(t);
fail1:
		t = NULL;
	}

	return t;
}
*/
import "C"

import (
	"errors"
	"reflect"
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
		err = C.init(C.Uint32(flags))
	})
	if err != 0 {
		return getError()
	}

	initKeyboardStateSlices()

	return nil
}

func Quit() {
	mainThreadCall(func() {
		C.SDL_Quit()
	})
}

var Running = true

func HandleEvents() {
	mainThreadCall(func() {
		C.handleEvents()
	})

	Running = C.running == 1
}

var keyboardPrev []uint8
var keyboardCurr []uint8

func KeyHeld(scancode int) bool {
	return keyboardCurr[scancode] == 1
}

func KeyPressed(scancode int) bool {
	return (keyboardCurr[scancode] &^ keyboardPrev[scancode]) == 1
}

func KeyReleased(scancode int) bool {
	return (keyboardPrev[scancode] &^ keyboardCurr[scancode]) == 1
}

func initKeyboardStateSlices() {
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&C.keyboardPrev)),
		Len:  C.SDL_NUM_SCANCODES,
		Cap:  C.SDL_NUM_SCANCODES,
	}
	keyboardPrev = *(*[]uint8)(unsafe.Pointer(&hdr))

	hdr = reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(C.keyboardCurr)),
		Len:  C.SDL_NUM_SCANCODES,
		Cap:  C.SDL_NUM_SCANCODES,
	}
	keyboardCurr = *(*[]uint8)(unsafe.Pointer(&hdr))
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

func CreateTexture(r *Renderer, w, h int, pixels []uint32) (*Texture, error) {
	var t *C.SDL_Texture
	mainThreadCall(func() {
		t = C.CreateTexture((*C.SDL_Renderer)(r), C.int(w), C.int(h),
			(*C.Uint32)(&pixels[0]))
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
