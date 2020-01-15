package util // import jdel.org/gosspks/util

import (
	"time"

	log "github.com/sirupsen/logrus"
)

var loggerTime = log.WithFields(log.Fields{
	"module": WhereAmI(),
})

// Duration wraps time.Duration to provide
// Round method
type Duration time.Duration

// Round Rounds d to the closest r unit
func (d Duration) Round(r time.Duration) time.Duration {
	duration := time.Duration(d)
	if r <= 0 {
		return duration
	}
	neg := duration < 0
	if neg {
		duration = -duration
	}
	if m := duration % r; m+m < r {
		duration = duration - m
	} else {
		duration = duration + r - m
	}
	if neg {
		return -duration
	}
	return duration
}
