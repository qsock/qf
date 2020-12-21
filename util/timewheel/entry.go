package timewheel

import (
	"github.com/qsock/qf/concurrent"
	"time"
)

// Entry is the timing job entry to wheel.
type Entry struct {
	wheel         *wheel              // Belonged wheel.
	job           JobFunc             // The job function.
	singleton     *concurrent.Boolean // Singleton mode.
	status        *concurrent.Int64   // Job status.
	times         *concurrent.Int64   // Limit running times.
	create        int64               // Timer ticks when the job installed.
	interval      int64               // The interval ticks of the job.
	createMs      int64               // The timestamp in milliseconds when job installed.
	intervalMs    int64               // The interval milliseconds of the job.
	rawIntervalMs int64               // Raw input interval in milliseconds.
}

// JobFunc is the job function.
type JobFunc = func()

// addEntry adds a timing job to the wheel.
func (w *wheel) addEntry(interval time.Duration, job JobFunc, singleton bool, times int64, status int64) *Entry {
	if times <= 0 {
		times = defaultTimes
	}
	var (
		ms  = interval.Nanoseconds() / 1e6
		num = ms / w.intervalMs
	)
	if num == 0 {
		// If the given interval is lesser than the one of the wheel,
		// then sets it to one tick, which means it will be run in one interval.
		num = 1
	}
	nowMs := time.Now().UnixNano() / 1e6
	ticks := w.ticks.Get()
	entry := &Entry{
		wheel:         w,
		job:           job,
		times:         concurrent.NewInt64(times),
		status:        concurrent.NewInt64(status),
		create:        ticks,
		interval:      num,
		singleton:     concurrent.NewBoolean(singleton),
		createMs:      nowMs,
		intervalMs:    ms,
		rawIntervalMs: ms,
	}
	// Install the job to the list of the slot.
	w.slots[(ticks+num)%w.number].PushBack(entry)
	return entry
}

// addEntryByParent adds a timing job with parent entry.
func (w *wheel) addEntryByParent(interval int64, parent *Entry) *Entry {
	num := interval / w.intervalMs
	if num == 0 {
		num = 1
	}
	nowMs := time.Now().UnixNano() / 1e6
	ticks := w.ticks.Get()
	entry := &Entry{
		wheel:         w,
		job:           parent.job,
		times:         parent.times,
		status:        parent.status,
		create:        ticks,
		interval:      num,
		singleton:     parent.singleton,
		createMs:      nowMs,
		intervalMs:    interval,
		rawIntervalMs: parent.rawIntervalMs,
	}
	w.slots[(ticks+num)%w.number].PushBack(entry)
	return entry
}

// Status returns the status of the job.
func (entry *Entry) Status() int64 {
	return entry.status.Get()
}

// SetStatus custom sets the status for the job.
func (entry *Entry) SetStatus(status int64) int64 {
	return entry.status.GetAndSet(status)
}

// Start starts the job.
func (entry *Entry) Start() {
	entry.status.Set(StatusReady)
}

// Stop stops the job.
func (entry *Entry) Stop() {
	entry.status.Set(StatusStoped)
}

//Reset reset the job.
func (entry *Entry) Reset() {
	entry.status.Set(StatusReset)
}

// Close closes the job, and then it will be removed from the timer.
func (entry *Entry) Close() {
	entry.status.Set(StatusClosed)
}

// IsSingleton checks and returns whether the job in singleton mode.
func (entry *Entry) IsSingleton() bool {
	return entry.singleton.Get()
}

// SetSingleton sets the job singleton mode.
func (entry *Entry) SetSingleton(enabled bool) {
	entry.singleton.Set(enabled)
}

// SetTimes sets the limit running times for the job.
func (entry *Entry) SetTimes(times int64) {
	entry.times.Set(times)
}

// Run runs the job.
func (entry *Entry) Run() {
	entry.job()
}

// check checks if the job should be run in given ticks and timestamp milliseconds.
func (entry *Entry) check(nowTicks int64, nowMs int64) (runnable, addable bool) {
	switch entry.status.Get() {
	case StatusStoped:
		return false, true
	case StatusClosed:
		return false, false
	case StatusReset:
		return false, true
	}
	// Firstly checks using the ticks, this may be low precision as one tick is a little bit long.
	if diff := nowTicks - entry.create; diff > 0 && diff%entry.interval == 0 {
		// If not the lowest level wheel.
		if entry.wheel.level > 0 {
			diffMs := nowMs - entry.createMs
			switch {
			// Add it to the next slot, which means it will run on next interval.
			case diffMs < entry.wheel.timer.intervalMs:
				entry.wheel.slots[(nowTicks+entry.interval)%entry.wheel.number].PushBack(entry)
				return false, false

			// Normal rolls on the job.
			case diffMs >= entry.wheel.timer.intervalMs:
				// Calculate the leftover milliseconds,
				// if it is greater than the minimum interval, then re-install it.
				if leftMs := entry.intervalMs - diffMs; leftMs > entry.wheel.timer.intervalMs {
					// Re-calculate and re-installs the job proper slot.
					entry.wheel.timer.doAddEntryByParent(leftMs, entry)
					return false, false
				}
			}
		}
		// Singleton mode check.
		if entry.IsSingleton() {
			// Note that it is atomic operation to ensure concurrent safety.
			if entry.status.GetAndSet(StatusRuning) == StatusRuning {
				return false, true
			}
		}
		// Limit running times.
		times := entry.times.IncrementAndGet()
		if times <= 0 {
			// Note that it is atomic operation to ensure concurrent safety.
			if entry.status.GetAndSet(StatusClosed) == StatusClosed || times < 0 {
				return false, false
			}
		}
		// This means it does not limit the running times.
		// I know it's ugly, but it is surely high performance for running times limit.
		if times < 2000000000 && times > 1000000000 {
			entry.times.Set(defaultTimes)
		}
		return true, true
	}
	return false, true
}
