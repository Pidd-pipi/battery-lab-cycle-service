package main

import (
	"context"
	"testing"
	"time"
)

func TestOpsAuditForHonorsCanceledContext(t *testing.T) {
	a := newOpsAudit()
	a.Add("rec-1", "created", "lab")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := a.For(ctx, "rec-1"); err == nil {
		t.Fatal("expected error for canceled context")
	}
}

func TestOpsAuditSinceHonorsCanceledContext(t *testing.T) {
	a := newOpsAudit()
	a.Add("rec-1", "created", "lab")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := a.Since(ctx, time.Now().Add(-time.Hour)); err == nil {
		t.Fatal("expected error for canceled context")
	}
}

func TestOpsAuditForKeepsHistory(t *testing.T) {
	a := newOpsAudit()
	a.Add("rec-2", "created", "lab")
	a.Add("rec-1", "created", "lab")
	a.Add("rec-1", "status_changed", "operator")
	events, err := a.For(context.Background(), "rec-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fatalf("filtered len=%d want 2", len(events))
	}
	if len(a.events) != 3 {
		t.Fatalf("history len=%d want 3", len(a.events))
	}
	if a.events[0].RecordID != "rec-2" {
		t.Fatalf("history corrupted: first record=%q", a.events[0].RecordID)
	}
}

func TestOpsAuditSinceKeepsHistory(t *testing.T) {
	a := newOpsAudit()
	a.Add("rec-2", "created", "lab")
	a.Add("rec-1", "created", "lab")
	a.Add("rec-1", "status_changed", "operator")
	a.events[0].At = time.Now().Add(-2 * time.Hour).UTC().Format(time.RFC3339Nano)
	events, err := a.Since(context.Background(), time.Now().Add(-time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fatalf("since len=%d want 2", len(events))
	}
	if len(a.events) != 3 {
		t.Fatalf("history len=%d want 3", len(a.events))
	}
	if a.events[0].RecordID != "rec-2" {
		t.Fatalf("history corrupted: first record=%q", a.events[0].RecordID)
	}
}
