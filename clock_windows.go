package sdl

/*
#cgo CFLAGS: -O3
#include <stdint.h>
#include <Windows.h>

static int64_t freq;
static int64_t base;

void
clockInit(void) {
	QueryPerformanceFrequency((LARGE_INTEGER *)(&freq));
	QueryPerformanceCounter((LARGE_INTEGER *)(&base));
}

int64_t
clockElapsed() {
	int64_t new;
	QueryPerformanceCounter((LARGE_INTEGER *)(&new));

	int64_t elapsed;
	elapsed = ((new - base) * 1000000000) / freq;
	base = new;
	return elapsed;
}
*/
import "C"

import "time"

var initialized bool

func ClockInit() {
	C.clockInit()
	initialized = true
}

func ClockElapsed() time.Duration {
	if !initialized {
		panic("clock not initialized")
	}

	var elapsed time.Duration
	elapsed = time.Duration(C.clockElapsed())
	return elapsed
}
