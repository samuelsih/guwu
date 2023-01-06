package limiter

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSimpleLimiter(t *testing.T) {
	var mu sync.Mutex
	var counter int

	var limiter = New(5)

	for i := 0; i < 5; i++ {
		limiter.Go(func() {
			mu.Lock()
			defer mu.Unlock()

			counter++
		})
	}

	limiter.Wait()

	if counter != 5 {
		t.Fatalf("TestSimpleLimiter - expected 5, got %v", counter)
	}
}

func TestLimiterCtx(t *testing.T) {
	var mu sync.Mutex

	var limiter = New(5)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var counter int

	for i := 0; i < 5; i++ {
		limiter.GoWithCtx(ctx, func() {
			time.Sleep(3 * time.Second)

			mu.Lock()
			defer mu.Unlock()

			counter++
		})
	}

	if counter != 0 {
		t.Fatalf("TestLimiterCtx - expected 0, got %v", counter)
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
