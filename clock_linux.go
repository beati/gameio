package sdl

/*
#cgo CFLAGS: -O3
#include <stdint.h>
#include <time.h>

static struct timespec base;

int
clockInit(void) {
	return clock_gettime(CLOCK_MONOTONIC_RAW, &base);
}

int64_t
clockElapsed() {
	struct timespec new;
	int err = clock_gettime(CLOCK_MONOTONIC_RAW, &new);
	if (err != 0) {
		return -1;
	}

	struct timespec elapsed;
	elapsed.tv_sec = new.tv_sec - base.tv_sec;
	elapsed.tv_nsec = new.tv_nsec - base.tv_nsec;
	base = new;
	return elapsed.tv_sec * 1000000000 + elapsed.tv_nsec;
}
*/
import "C"

import "time"

var initialized bool

func ClockInit() {
	var err C.int
	err = C.clockInit()
	if err != 0 {
		panic("clock_gettime error")
	}
	initialized = true
}

func ClockElapsed() time.Duration {
	if !initialized {
		panic("clock not initialized")
	}

	elapsed := time.Duration(C.clockElapsed())
	if elapsed < 0 {
		panic("clock_gettime error")
	}
	return elapsed
}
