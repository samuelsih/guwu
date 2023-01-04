package limiter

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestSimpleLimiter(t *testing.T) {
	var counter atomic.Uint32

	var limiter = New(5)

	for i := 0; i < 5; i++ {
		limiter.Go(func() {
			counter.Add(1)
		})
	}

	limiter.Wait()

	if counter.Load() != 5 {
		t.Fatalf("TestSimpleLimiter - expected 5, got %v", counter.Load())
	}
}

func TestLimiterCtx(t *testing.T) {
	var limiter = New(5)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var counter atomic.Uint32

	for i := 0; i < 5; i++ {
		limiter.GoWithCtx(ctx, func() {
			time.Sleep(3 * time.Second)
			counter.Add(1)
		})
	}

	if counter.Load() != 0 {
		t.Fatalf("TestLimiterCtx - expected 0, got %v", counter.Load())
	}

	limiter.Wait()
}

func TestLimiterPanic(t *testing.T) {
	var limiter = New(1)

	limiter.Go(func() {
		panic("foo")
	})

	if limiter.PanicErr != nil {
		t.Fatalf("expected panic: foo, got %v", limiter.PanicErr)
	}
}
