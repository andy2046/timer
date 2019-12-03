# timer

[![GoDoc](https://godoc.org/github.com/andy2046/timer?status.svg)](https://godoc.org/github.com/andy2046/timer)

timer with less race condition.

## Install

```
go get github.com/andy2046/timer
```

## Example

```go
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
```
