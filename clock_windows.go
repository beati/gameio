package gameio

/*
#cgo CFLAGS: -O3
#include <stdint.h>
#include <Windows.h>

static int64_t freq;

static void
initFreq(void) {
	QueryPerformanceFrequency((LARGE_INTEGER *)(&freq));
}

static int64_t
initClock(void) {
	int64_t base;
	QueryPerformanceCounter((LARGE_INTEGER *)(&base));
	return base;
}

static int64_t
elapsed(int64_t *base) {
	int64_t new;
	QueryPerformanceCounter((LARGE_INTEGER *)(&new));

	int64_t elapsed = ((new - *base) * 1000000000) / freq;
	*base = new;
	return elapsed;
}
*/
import "C"

import "time"

type Clock C.int64_t

func init() {
	C.initFreq()
}

func InitClock() Clock {
	return (Clock)(C.initClock())
}

func (c *Clock) Elapsed() time.Duration {
	return time.Duration(C.elapsed((*C.int64_t)(c)))
}
