package gameio

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

SDL_GameController *controller = NULL;
Uint8 buttonsPrev[SDL_CONTROLLER_BUTTON_MAX];
Uint8 buttonsCurr[SDL_CONTROLLER_BUTTON_MAX];

void
handleEvents(void) {
	memcpy(keyboardPrev, keyboardCurr, numScancodes);
	memcpy(buttonsPrev, buttonsCurr, SDL_CONTROLLER_BUTTON_MAX);

	if (controller != NULL && !SDL_GameControllerGetAttached(controller)) {
		SDL_GameControllerClose(controller);
		controller = NULL;
	}

	if (controller == NULL) {
		int i = 0;
		for (i = 0; i < SDL_NumJoysticks(); i++) {
			if (SDL_IsGameController(i)) {
				controller = SDL_GameControllerOpen(i);
				if (controller != NULL) {
					break;
				}
			}
		}
	}

	SDL_Event event;
	while (SDL_PollEvent(&event)) {
		switch (event.type) {
		case SDL_QUIT:
			running = 0;
			break;
		}
	}

	if (controller != NULL) {
		int i = 0;
		for (i = 0; i < SDL_CONTROLLER_BUTTON_MAX; i++) {
			buttonsCurr[i] = SDL_GameControllerGetButton(controller, i);
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

void
quit(void) {
	if (controller != NULL) {
		SDL_GameControllerClose(controller);
	}
	SDL_Quit();
}

SDL_Rect srcRect;
SDL_Rect dstRect;
SDL_Point center;
int
RenderCopyEx(SDL_Renderer *r, SDL_Texture *t, double angle, SDL_RendererFlip flip) {
	return SDL_RenderCopyEx(r, t, &srcRect, &dstRect, angle, &center, flip);
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

func getError() error {
	err := C.SDL_GetError()
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

func Init(flags int) error {
	err := C.init(C.Uint32(flags))
	if err != 0 {
		return getError()
	}

	initKeyboardStateSlices()
	initButtonsStateSlices()

	return nil
}

func Quit() {
	C.quit()
}

var Running = true

func HandleEvents() {
	C.handleEvents()

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

var buttonsPrev []uint8
var buttonsCurr []uint8

const (
	Button_A             = C.SDL_CONTROLLER_BUTTON_A
	Button_B             = C.SDL_CONTROLLER_BUTTON_B
	Button_X             = C.SDL_CONTROLLER_BUTTON_X
	Button_Y             = C.SDL_CONTROLLER_BUTTON_Y
	Button_BACK          = C.SDL_CONTROLLER_BUTTON_BACK
	Button_GUIDE         = C.SDL_CONTROLLER_BUTTON_GUIDE
	Button_START         = C.SDL_CONTROLLER_BUTTON_START
	Button_LEFTSTICK     = C.SDL_CONTROLLER_BUTTON_LEFTSTICK
	Button_RIGHTSTICK    = C.SDL_CONTROLLER_BUTTON_RIGHTSTICK
	Button_LEFTSHOULDER  = C.SDL_CONTROLLER_BUTTON_LEFTSHOULDER
	Button_RIGHTSHOULDER = C.SDL_CONTROLLER_BUTTON_RIGHTSHOULDER
	Button_UP            = C.SDL_CONTROLLER_BUTTON_DPAD_UP
	Button_DOWN          = C.SDL_CONTROLLER_BUTTON_DPAD_DOWN
	Button_LEFT          = C.SDL_CONTROLLER_BUTTON_DPAD_LEFT
	Button_RIGHT         = C.SDL_CONTROLLER_BUTTON_DPAD_RIGHT
)

func ButtonHeld(button int) bool {
	return buttonsCurr[button] == 1
}

func ButtonPressed(button int) bool {
	return (buttonsCurr[button] &^ buttonsPrev[button]) == 1
}

func ButtonReleased(button int) bool {
	return (buttonsPrev[button] &^ buttonsCurr[button]) == 1
}

func initButtonsStateSlices() {
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&C.buttonsPrev)),
		Len:  C.SDL_CONTROLLER_BUTTON_MAX,
		Cap:  C.SDL_CONTROLLER_BUTTON_MAX,
	}
	buttonsPrev = *(*[]uint8)(unsafe.Pointer(&hdr))

	hdr = reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&C.buttonsCurr)),
		Len:  C.SDL_CONTROLLER_BUTTON_MAX,
		Cap:  C.SDL_CONTROLLER_BUTTON_MAX,
	}
	buttonsCurr = *(*[]uint8)(unsafe.Pointer(&hdr))
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

func CreateWindow(title string, x, y, w, h int, flags int) (*Window, error) {
	t := C.CString(title)
	defer C.free(unsafe.Pointer(t))

	win := C.SDL_CreateWindow(t, C.int(x), C.int(y), C.int(w), C.int(h),
		C.Uint32(flags))
	if win == nil {
		return nil, getError()
	}
	return (*Window)(win), nil
}

func (w *Window) Destroy() {
	C.SDL_DestroyWindow((*C.SDL_Window)(w))
}

type Renderer C.SDL_Renderer

const (
	RendererSoftware      = C.SDL_RENDERER_SOFTWARE
	RendererAccelerated   = C.SDL_RENDERER_ACCELERATED
	RendererPresentVSync  = C.SDL_RENDERER_PRESENTVSYNC
	RendererTargetTexture = C.SDL_RENDERER_TARGETTEXTURE
)

func CreateRenderer(w *Window, index int, flags int) (*Renderer, error) {
	r := C.SDL_CreateRenderer((*C.SDL_Window)(w), C.int(index), C.Uint32(flags))
	if r == nil {
		return nil, getError()
	}
	return (*Renderer)(r), nil
}

func (r *Renderer) Destroy() {
	C.SDL_DestroyRenderer((*C.SDL_Renderer)(r))
}

func (r *Renderer) SetDrawColor(red, green, blue, alpha uint8) error {
	err := C.SDL_SetRenderDrawColor((*C.SDL_Renderer)(r),
		C.Uint8(red), C.Uint8(green), C.Uint8(blue), C.Uint8(alpha))
	if err != 0 {
		return getError()
	}
	return nil
}

func (r *Renderer) SetLogicalSize(w, h int) error {
	err := C.SDL_RenderSetLogicalSize((*C.SDL_Renderer)(r), C.int(w), C.int(h))
	if err != 0 {
		return getError()
	}
	return nil
}

func (r *Renderer) Clear() error {
	err := C.SDL_RenderClear((*C.SDL_Renderer)(r))
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
	if src != nil {
		C.srcRect.x = C.int(src.X)
		C.srcRect.y = C.int(src.Y)
		C.srcRect.w = C.int(src.W)
		C.srcRect.h = C.int(src.H)
	}
	if dst != nil {
		C.dstRect.x = C.int(dst.X)
		C.dstRect.y = C.int(dst.Y)
		C.dstRect.w = C.int(dst.W)
		C.dstRect.h = C.int(dst.H)
	}
	if center != nil {
		C.center.x = C.int(center.X)
		C.center.y = C.int(center.Y)
	}
	err := C.RenderCopyEx((*C.SDL_Renderer)(r), (*C.SDL_Texture)(t), C.double(angle),
		C.SDL_RendererFlip(flip))
	if err != 0 {
		return getError()
	}
	return nil
}

func (r *Renderer) Present() {
	C.SDL_RenderPresent((*C.SDL_Renderer)(r))
}

type Texture C.SDL_Texture

func CreateTexture(r *Renderer, w, h int, pixels []uint32) (*Texture, error) {
	t := C.CreateTexture((*C.SDL_Renderer)(r), C.int(w), C.int(h),
		(*C.Uint32)(&pixels[0]))
	if t == nil {
		return nil, getError()
	}

	return (*Texture)(t), nil
}

func (t *Texture) Destroy() {
	C.SDL_DestroyTexture((*C.SDL_Texture)(t))
}
