package gameio

/*
#cgo LDFLAGS: -lfmod
#cgo CFLAGS: -O3
#include <fmod/fmod.h>
#include <fmod/fmod_errors.h>
#include <stdlib.h>

FMOD_SYSTEM *sys;

FMOD_RESULT FMODErr;

int
initFMOD(void) {
	FMODErr = FMOD_System_Create(&sys);
	if (FMODErr != FMOD_OK) {
		return -1;
	}

	FMODErr = FMOD_System_Init(sys, 32, FMOD_INIT_NORMAL, NULL);
	if (FMODErr != FMOD_OK) {
		return -1;
	}

	return 0;
}

int
quitFMOD(void) {
	FMODErr = FMOD_System_Close(sys);
	if (FMODErr != FMOD_OK) {
		return -1;
	}

	FMODErr = FMOD_System_Release(sys);
	if (FMODErr != FMOD_OK) {
		return -1;
	}

	return 0;
}

FMOD_SOUND *
createSound(const char *path) {
	FMOD_SOUND *sound;
	FMODErr = FMOD_System_CreateSound(sys, path, FMOD_LOOP_OFF, NULL, &sound);
	if (FMODErr != FMOD_OK) {
		return NULL;
	}
	return sound;
}

int
playSound(FMOD_SOUND *s) {
	FMODErr = FMOD_System_PlaySound(sys, s, NULL, 0, NULL);
	if (FMODErr != FMOD_OK) {
		return -1;
	}
	return 0;
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

func getFMODError() error {
	return errors.New(C.GoString(C.FMOD_ErrorString(C.FMODErr)))
}

func InitFMOD() error {
	err := C.initFMOD()
	if err != 0 {
		return getFMODError()
	}
	return nil
}

func QuitFMOD() error {
	err := C.quitFMOD()
	if err != 0 {
		return getFMODError()
	}
	return nil
}

func UpdateFMOD() error {
	C.FMODErr = C.FMOD_System_Update(C.sys)
	if C.FMODErr != C.FMOD_OK {
		return getFMODError()
	}
	return nil
}

type Sound C.FMOD_SOUND

func CreateSound(path string) (*Sound, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	s := C.createSound(cpath)
	if s == nil {
		return nil, getFMODError()
	}
	return (*Sound)(s), nil
}

func (s *Sound) Destroy() {
	C.FMOD_Sound_Release((*C.FMOD_SOUND)(s))
}

func (s *Sound) Play() error {
	err := C.playSound((*C.FMOD_SOUND)(s))
	if err != 0 {
		return getFMODError()
	}
	return nil
}
