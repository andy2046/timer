package timer

import (
	"time"
)

// The Timer type represents a single event.
// When the Timer expires, the current time will be sent on C,
// unless the Timer was created by AfterFunc.
type Timer struct {
	C <-chan time.Time
	r *time.Timer
}

// New creates a new Timer that will send
// the current time on its channel after at least duration d.
func New(d time.Duration) *Timer {
	r := time.NewTimer(d)
	return &Timer{
		C: r.C,
		r: r,
	}
}

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer, false if the timer has already
// expired or been stopped.
// Stop does not close the channel, to prevent a read from the channel succeeding
// incorrectly.
func (t *Timer) Stop() bool {
	return t.r.Stop()
}

// Reset changes the timer to expire after duration d.
// It returns true if the timer had been active, false if the timer had
// expired or been stopped.
func (t *Timer) Reset(d time.Duration) bool {
	stopped := t.r.Stop()
	if !stopped {
		select {
		case <-t.r.C:
		default:
		}
	}
	t.r.Reset(d)
	return stopped
}

// After waits for the duration to elapse and then sends the current time
// on the returned channel.
// It is equivalent to New(d).C.
func After(d time.Duration) <-chan time.Time {
	return New(d).C
}

// AfterFunc waits for the duration to elapse and then calls f
// in its own goroutine. It returns a Timer that can
// be used to cancel the call using its Stop method.
func AfterFunc(d time.Duration, f func()) *Timer {
	r := time.AfterFunc(d, f)
	return &Timer{
		C: r.C,
		r: r,
	}
}
