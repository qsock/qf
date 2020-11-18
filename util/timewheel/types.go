package timewheel

import (
	"errors"
	"math"
	"time"
)

const (
	StatusReady          = 0                     // Job is ready for running.
	StatusRuning         = 1                     // Job is already running.
	StatusStoped         = 2                     // Job is stopped.
	StatusReset          = 3                     // Job is reset.
	StatusClosed         = -1                    // Job is closed and waiting to be deleted.
	defaultTimes         = math.MaxInt32         // Default limit running times, a big number.
	defaultSlotNumber    = 10                    // Default slot number.
	defaultWheelInterval = 10 * time.Millisecond // Default wheel interval.
	defaultWheelLevel    = 6                     // Default wheel level.
)

var (
	ErrClosed = errors.New("closed")
)
