package gameio

/*
#cgo CFLAGS: -O3
#include <stdint.h>
#include <Windows.h>

int clockInit = 0;
static int64_t freq;

int64_t
newClock(void) {
	if (!clockInit) {
		QueryPerformanceFrequency((LARGE_INTEGER *)(&freq));
		clockInit = 1;
	}
	int64_t base;
	QueryPerformanceCounter((LARGE_INTEGER *)(&base));
	return base;
}

int64_t
elapsed(int64_t *base) {
	int64_t new;
	QueryPerformanceCounter((LARGE_INTEGER *)(&new));

	int64_t elapsed;
	elapsed = ((new - *base) * 1000000000) / freq;
	*base = new;
	return elapsed;
}
*/
import "C"

import "time"

type Clock C.int64_t

func NewClock() *Clock {
	c := C.newClock()
	return (*Clock)(&c)
}

func (c *Clock) Elapsed() time.Duration {
	if *c == 0 {
		panic("clock not initialized")
	}
	return time.Duration(C.elapsed((*C.int64_t)(c)))
}
