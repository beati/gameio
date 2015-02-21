package gameio

/*
#cgo CFLAGS: -O3
#include <stdint.h>
#include <time.h>

static int
initClock(struct timespec *c) {
	return clock_gettime(CLOCK_MONOTONIC_RAW, c);
}

static int64_t
elapsed(struct timespec *base) {
	struct timespec new;
	int err = clock_gettime(CLOCK_MONOTONIC_RAW, &new);
	if (err != 0) {
		return -1;
	}

	struct timespec elapsed;
	elapsed.tv_sec = new.tv_sec - base->tv_sec;
	elapsed.tv_nsec = new.tv_nsec - base->tv_nsec;
	*base = new;
	return elapsed.tv_sec * 1000000000 + elapsed.tv_nsec;
}
*/
import "C"

import "time"

type Clock C.struct_timespec

func InitClock() Clock {
	var c Clock
	err := C.initClock((*C.struct_timespec)(&c))
	if err != 0 {
		panic("clock_gettime error")
	}
	return c
}

func (c *Clock) Elapsed() time.Duration {
	elapsed := time.Duration(C.elapsed((*C.struct_timespec)(c)))
	if elapsed < 0 {
		panic("clock_gettime error")
	}
	return elapsed
}
