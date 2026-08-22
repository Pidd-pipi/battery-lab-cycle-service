package main

import (
	"context"
	"sync"
	"testing"
)

func TestBatchTransitionCompletesAllJobs(t *testing.T) {
	svc := newOpsService([]OpsRecord{seedOpsRecord("rec-1"), seedOpsRecord("rec-2"), seedOpsRecord("rec-3")})
	outcomes, err := svc.BatchTransition(context.Background(), []string{"rec-1", "rec-2", "rec-3"}, OpsStatusClosed, "operator")
	if err != nil {
		t.Fatal(err)
	}
	if len(outcomes) != 3 {
		t.Fatalf("outcomes=%d want 3", len(outcomes))
	}
	seen := map[string]bool{}
	for _, o := range outcomes {
		seen[o.ID] = true
		if !o.OK {
			t.Fatalf("job %s failed: %v", o.ID, o.Err)
		}
	}
	for _, id := range []string{"rec-1", "rec-2", "rec-3"} {
		if !seen[id] {
			t.Fatalf("missing outcome for %s", id)
		}
	}
}

func TestBatchTransitionReportsErrors(t *testing.T) {
	svc := newOpsService([]OpsRecord{seedOpsRecord("rec-1")})
	outcomes, err := svc.BatchTransition(context.Background(), []string{"rec-1", "rec-missing"}, OpsStatusClosed, "operator")
	if err != nil {
		t.Fatal(err)
	}
	if len(outcomes) != 2 {
		t.Fatalf("outcomes=%d want 2", len(outcomes))
	}
	byID := map[string]BatchOutcome{}
	for _, o := range outcomes {
		byID[o.ID] = o
	}
	if byID["rec-missing"].OK {
		t.Fatal("missing record should not be OK")
	}
	if byID["rec-missing"].Err == nil {
		t.Fatal("missing record should carry an error")
	}
	if !byID["rec-1"].OK {
		t.Fatal("rec-1 should succeed")
	}
}

func TestBatchTransitionHonorsCanceledContext(t *testing.T) {
	svc := newOpsService([]OpsRecord{seedOpsRecord("rec-1")})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	outcomes, err := svc.BatchTransition(ctx, []string{"rec-1"}, OpsStatusClosed, "operator")
	if err != nil {
		t.Fatal(err)
	}
	if len(outcomes) != 0 {
		t.Fatalf("canceled context should skip all jobs, got %d outcomes", len(outcomes))
	}
}

func TestBatchTransitionWorkerCount(t *testing.T) {
	svc := newOpsService([]OpsRecord{seedOpsRecord("rec-1"), seedOpsRecord("rec-2"), seedOpsRecord("rec-3"), seedOpsRecord("rec-4")})
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			outcomes, err := svc.BatchTransition(context.Background(), []string{"rec-1", "rec-2", "rec-3", "rec-4"}, OpsStatusClosed, "operator")
			if err != nil {
				t.Error(err)
			}
			if len(outcomes) != 4 {
				t.Errorf("outcomes=%d want 4", len(outcomes))
			}
		}()
	}
	close(start)
	wg.Wait()
}
