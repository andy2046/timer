package timer_test

import (
	"fmt"
	"github.com/andy2046/timer"
	"time"
)

func ExampleTimer() {
	interval := 5 * time.Second
	quit := make(chan struct{}, 1)

	timer := timer.New(interval)
	defer timer.Stop()

	go func() {
		time.Sleep(3 * interval)
		quit <- struct{}{}
	}()

	for {
		timer.Reset(interval)
		select {
		case <-quit:
			fmt.Println("exit")
			return
		case <-timer.C:
			fmt.Println("timer", time.Now())
			continue
		}
	}
}
