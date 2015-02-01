package sdl

/*
#cgo CFLAGS: -O3
#include <mach/mach_time.h>
#include <stdint.h>

static mach_timebase_info_data_t timebaseInfo;
static uint64_t base;

void
clockInit(void) {
	mach_timebase_info(&timebaseInfo);
	base = mach_absolute_time();
}

uint64_t
clockElapsed() {
	uint64_t new = mach_absolute_time();

	uint64_t elapsed = ((new - base) * timebaseInfo.numer) / timebaseInfo.denom;
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

	return time.Duration(C.clockElapsed())
}
