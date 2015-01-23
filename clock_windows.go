package sdl

/*
#cgo CFLAGS: -O3
#include <Windows.h>

LONGLONG freq;
LONGLONG base;

void
clockInit(void) {
	QueryPerformanceFrequency((LARGE_INTEGER *)(&freq));
	QueryPerformanceCounter((LARGE_INTEGER *)(&base));
}

LONGLONG
clockElapsed() {
	LONGLONG new;
	QueryPerformanceCounter((LARGE_INTEGER *)(&new));

	LONGLONG elapsed;
	elapsed = (new - base) * 1000000000;
	elapsed /= freq;
	base = new;
	return elapsed;
}
*/
import "C"

import "time"

var initialized bool

func ClockInit() {
	mainThreadCall(func() {
		C.clockInit()
	})
	initialized = true
}

func ClockElapsed() time.Duration {
	if !initialized {
		panic("clock not initialized")
	}

	var elapsed time.Duration
	mainThreadCall(func() {
		elapsed = time.Duration(C.clockElapsed())
	})
	return elapsed
}
