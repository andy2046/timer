package timer

import (
	"errors"
	"runtime"
	"testing"
	"time"
)

const windowsInaccuracy = 17 * time.Millisecond

func testReset(d time.Duration) error {
	t0 := New(2 * d)
	time.Sleep(d)
	if !t0.Reset(3 * d) {
		return errors.New("resetting unfired timer returned false")
	}
	time.Sleep(2 * d)
	select {
	case <-t0.C:
		return errors.New("timer fired early")
	default:
	}
	time.Sleep(2 * d)
	select {
	case <-t0.C:
	default:
		return errors.New("reset timer did not fire")
	}

	if t0.Reset(50 * time.Millisecond) {
		return errors.New("resetting expired timer returned true")
	}
	return nil
}

func TestReset(t *testing.T) {
	const unit = 25 * time.Millisecond
	tries := []time.Duration{
		1 * unit,
		3 * unit,
		7 * unit,
		15 * unit,
	}
	var err error
	for _, d := range tries {
		err = testReset(d)
		if err == nil {
			t.Logf("passed using duration %v", d)
			return
		}
	}
	t.Error(err)
}

func TestAfterStop(t *testing.T) {
	AfterFunc(100*time.Millisecond, func() {})
	t0 := New(50 * time.Millisecond)
	c1 := make(chan bool, 1)
	t1 := AfterFunc(150*time.Millisecond, func() { c1 <- true })
	c2 := After(200 * time.Millisecond)
	if !t0.Stop() {
		t.Fatalf("failed to stop event 0")
	}
	if !t1.Stop() {
		t.Fatalf("failed to stop event 1")
	}
	<-c2
	select {
	case <-t0.C:
		t.Fatalf("event 0 was not stopped")
	case <-c1:
		t.Fatalf("event 1 was not stopped")
	default:
	}
	if t1.Stop() {
		t.Fatalf("Stop returned true twice")
	}
}

func TestAfterFunc(t *testing.T) {
	i := 10
	c := make(chan bool)
	var f func()
	f = func() {
		i--
		if i >= 0 {
			AfterFunc(0, f)
			time.Sleep(1 * time.Second)
		} else {
			c <- true
		}
	}

	AfterFunc(0, f)
	<-c
}

func TestAfter(t *testing.T) {
	const delay = 100 * time.Millisecond
	start := time.Now()
	end := <-After(delay)
	delayadj := delay
	if runtime.GOOS == "windows" {
		delayadj -= windowsInaccuracy
	}
	if duration := time.Now().Sub(start); duration < delayadj {
		t.Fatalf("After(%s) slept for only %d ns", delay, duration)
	}
	if min := start.Add(delayadj); end.Before(min) {
		t.Fatalf("After(%s) expect >= %s, got %s", delay, min, end)
	}
}
