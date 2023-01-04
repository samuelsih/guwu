package limiter

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

type Limiter struct {
	PanicErr error

	worker sync.WaitGroup
	limit  chan struct{}
	once   sync.Once
}

func New(amount int) *Limiter {
	if amount <= 0 {
		return &Limiter{limit: make(chan struct{}, 1)}
	}

	return &Limiter{limit: make(chan struct{}, amount)}
}

func (lim *Limiter) Go(f func()) {
	lim.limit <- struct{}{}

	lim.worker.Add(1)

	go func() {
		defer func() {
			<-lim.limit
			lim.worker.Done()
		}()

		defer lim.catchPanic()

		f()
	}()
}

func (lim *Limiter) GoWithCtx(ctx context.Context, f func()) {
	if ctx.Err() != nil {
		return
	}

	select {
	case lim.limit <- struct{}{}:

	case <-ctx.Done():
		return
	}

	lim.worker.Add(1)

	go func() {
		defer func() {
			<-lim.limit
			lim.worker.Done()
		}()

		defer lim.catchPanic()

		f()
	}()
}

func (lim *Limiter) Wait() {
	lim.worker.Wait()
}

func (lim *Limiter) catchPanic() {
	if r := recover(); r != nil {
		lim.once.Do(func() {
			lim.PanicErr = getPanicValue(r)
		})
	}
}

func getPanicValue(val any) error {
	switch v := val.(type) {
	case nil:
		return nil
	case string:
		return fmt.Errorf("panic: %v\n%s", v, getStack())
	case error:
		return fmt.Errorf("panic in limiter %w\n%s", v, getStack())
	default:
		return fmt.Errorf("unhandle panic: %+v\n%s", v, getStack())
	}
}

func getStack() []byte {
	buf := make([]byte, 64<<10)
	buf = buf[:runtime.Stack(buf, false)]
	return buf
}
