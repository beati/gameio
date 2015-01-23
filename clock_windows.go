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
	if (new >= base) {
		elapsed = (new - base) * 1000000000;
	} else {
		elapsed = 0;
		elapsed += 0xFFFFFFFFFFFFFFFF - base;
		elapsed += new;
	}

	elapsed /= freq;
	base = new;
	return elapsed;
}
*/
import "C"

import "time"

func ClockInit() {
	mainThreadCall(func() {
		C.clockInit()
	})
}

func ClockElapsed() time.Duration {
	var elapsed time.Duration
	mainThreadCall(func() {
		elapsed = time.Duration(C.clockElapsed())
	})
	return elapsed
}
