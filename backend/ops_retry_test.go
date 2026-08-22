package main

import (
	"context"
	"errors"
	"testing"
)

func TestRunWithRetryStopsOnDomainError(t *testing.T) {
	attempts := 0
	err := RunWithRetry(context.Background(), RetryPolicy{MaxAttempts: 3}, func() error {
		attempts++
		return ErrOpsConflict
	})
	if err == nil || !errors.Is(err, ErrOpsConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
	if attempts != 1 {
		t.Fatalf("attempts=%d want 1", attempts)
	}
}

func TestRunWithRetryWrappedNotFoundStops(t *testing.T) {
	attempts := 0
	err := RunWithRetry(context.Background(), RetryPolicy{MaxAttempts: 3}, func() error {
		attempts++
		return wrapOps("get", "store.get", ErrOpsNotFound)
	})
	if err == nil || !errors.Is(err, ErrOpsNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
	if attempts != 1 {
		t.Fatalf("attempts=%d want 1", attempts)
	}
}

func TestRunWithRetrySucceedsOnRetry(t *testing.T) {
	attempts := 0
	err := RunWithRetry(context.Background(), RetryPolicy{MaxAttempts: 3}, func() error {
		attempts++
		if attempts < 2 {
			return errors.New("transient boom")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("attempts=%d want 2", attempts)
	}
}
