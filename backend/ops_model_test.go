package main

import (
	"context"
	"testing"
)

func TestOpsRecordCloneAllowsLabelWrite(t *testing.T) {
	rec := OpsRecord{ID: "rec-1", Subject: "cycle", Owner: "lab", Status: OpsStatusQueued, Priority: OpsPriorityNormal}
	cloned := rec.Clone()
	cloned.Labels["site"] = "s1"
	if rec.Labels["site"] != "" {
		t.Fatal("clone mutated original labels")
	}
}

func TestOpsNormalizeInitializesEmptyLabels(t *testing.T) {
	rec := normalizeOpsRecord(OpsRecord{ID: "REC-1", Subject: "  cycle  ", Owner: " lab ", Status: OpsStatusQueued, Priority: OpsPriorityNormal})
	if rec.Labels == nil {
		t.Fatal("labels not initialized")
	}
	rec.Labels["site"] = "s1"
}

func TestOpsNormalizeDefaultsRevision(t *testing.T) {
	rec := normalizeOpsRecord(OpsRecord{ID: "rec-1", Subject: "cycle", Owner: "lab", Status: OpsStatusQueued, Priority: OpsPriorityNormal})
	if rec.Revision != 1 {
		t.Fatalf("revision=%d want 1", rec.Revision)
	}
}

func TestOpsFirstUpdateExpectedRevision(t *testing.T) {
	svc := newOpsService(nil)
	created, err := svc.Create(context.Background(), OpsRecord{ID: "rec-1", Subject: "cycle", Owner: "lab", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s1"}})
	if err != nil {
		t.Fatal(err)
	}
	if created.Revision != 1 {
		t.Fatalf("created revision=%d want 1", created.Revision)
	}
	if _, err := svc.Transition(context.Background(), "rec-1", 1, OpsStatusActive, "operator"); err != nil {
		t.Fatalf("transition with expected revision failed: %v", err)
	}
}
