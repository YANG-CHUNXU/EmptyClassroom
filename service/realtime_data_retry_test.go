package service

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestShouldRetryRealtime(t *testing.T) {
	if shouldRetryRealtime(nil) {
		t.Fatalf("nil error should not retry")
	}

	if shouldRetryRealtime(errors.New("temporary")) != true {
		t.Fatalf("temporary errors should retry")
	}

	if shouldRetryRealtime(ErrLoginRejected) {
		t.Fatalf("login rejection should not retry")
	}

	if shouldRetryRealtime(errors.Join(ErrLoginRejected, errors.New("wrapped"))) {
		t.Fatalf("wrapped login rejection should not retry")
	}

	if shouldRetryRealtime(context.Canceled) {
		t.Fatalf("context canceled should not retry")
	}

	if shouldRetryRealtime(context.DeadlineExceeded) {
		t.Fatalf("context deadline exceeded should not retry")
	}
}

func TestWaitForRealtimeRetryCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	err := waitForRealtimeRetry(ctx, time.Second)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("waitForRealtimeRetry() error = %v, want context canceled", err)
	}
	if elapsed := time.Since(start); elapsed > 100*time.Millisecond {
		t.Fatalf("waitForRealtimeRetry() took too long to return: %v", elapsed)
	}
}

func TestWaitForRealtimeRetryDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := waitForRealtimeRetry(ctx, time.Second)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("waitForRealtimeRetry() error = %v, want context deadline exceeded", err)
	}
}
