package timewheel

import (
	"github.com/qsock/qf/concurrent"
	"time"
)

// start starts the ticker using a standalone goroutine.
func (w *wheel) start() {
	go func() {
		ticker := time.NewTicker(time.Duration(w.intervalMs) * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				switch w.timer.status.Get() {
				case StatusRuning:
					w.proceed()

				case StatusStoped:
				case StatusClosed:
					ticker.Stop()
					return
				}

			}
		}
	}()
}

// proceed checks and rolls on the job.
// If a timing job is time for running, it runs in an asynchronous goroutine,
// or else it removes from current slot and re-installs the job to another wheel and slot
// according to its leftover interval in milliseconds.
func (w *wheel) proceed() {
	n := w.ticks.IncrementAndGet()
	l := w.slots[int(n%w.number)]
	length := l.Len()
	if length > 0 {
		go func(l *concurrent.List, nowTicks int64) {
			entry := (*Entry)(nil)
			nowMs := time.Now().UnixNano() / 1e6
			for i := length; i > 0; i-- {
				if v := l.PopFront(); v == nil {
					break
				} else {
					entry = v.(*Entry)
				}
				// Checks whether the time for running.
				runnable, addable := entry.check(nowTicks, nowMs)
				if runnable {
					// Just run it in another goroutine.
					go func(entry *Entry) {
						defer func() {
							if err := recover(); err != nil {
								if err != ErrClosed {
									panic(err)
								} else {
									entry.Close()
								}
							}
							if entry.Status() == StatusRuning {
								entry.SetStatus(StatusReady)
							}
						}()
						entry.job()
					}(entry)
				}
				// If rolls on the job.
				if addable {
					//If STATUS_RESET , reset to runnable state.
					if entry.Status() == StatusReset {
						entry.SetStatus(StatusReady)
					}
					entry.wheel.timer.doAddEntryByParent(entry.rawIntervalMs, entry)
				}
			}
		}(l, n)
	}
}
