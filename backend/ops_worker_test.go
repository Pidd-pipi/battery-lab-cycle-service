package main

import (
	"context"
	"testing"
)

func TestOpsBatchWorkerSendsOutcomePerJob(t *testing.T) {
	svc := newOpsService([]OpsRecord{seedOpsRecord("rec-1"), seedOpsRecord("rec-2")})
	jobs := make(chan string, 2)
	results := make(chan BatchOutcome, 2)
	jobs <- "rec-1"
	jobs <- "rec-2"
	close(jobs)
	opsBatchWorker(context.Background(), jobs, results, svc, OpsStatusClosed, "operator")
	close(results)
	if len(results) != 2 {
		t.Fatalf("outcomes=%d want 2", len(results))
	}
	for r := range results {
		if !r.OK {
			t.Fatalf("job %s failed: %v", r.ID, r.Err)
		}
	}
}
