package main

import (
	"context"
	"testing"
)

func TestOpsTransitionHonorsCanceledContext(t *testing.T) {
	svc := newOpsService([]OpsRecord{seedOpsRecord("rec-1")})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := svc.Transition(ctx, "rec-1", 0, OpsStatusActive, "operator"); err == nil {
		t.Fatal("expected error for canceled context")
	}
}

func TestOpsSearchHonorsCanceledContext(t *testing.T) {
	svc := newOpsService([]OpsRecord{seedOpsRecord("rec-1")})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := svc.Search(ctx, OpsQuery{}); err == nil {
		t.Fatal("expected error for canceled context")
	}
}
