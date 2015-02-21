package gameio

/*
#cgo CFLAGS: -O3
#include <mach/mach_time.h>
#include <stdint.h>

mach_timebase_info_data_t timebaseInfo;

static uint64_t
initClock(void) {
	return mach_absolute_time();
}

static uint64_t
elapsed(uint64_t *base) {
	uint64_t new = mach_absolute_time();

	uint64_t elapsed = ((new - *base) * timebaseInfo.numer) / timebaseInfo.denom;
	*base = new;
	return elapsed;
}
*/
import "C"

import "time"

type Clock C.uint64_t

func init() {
	C.mach_timebase_info(&C.timebaseInfo)
}

func InitClock() Clock {
	return (Clock)(C.initClock())
}

func (c *Clock) Elapsed() time.Duration {
	return time.Duration(C.elapsed((*C.uint64_t)(c)))
}
