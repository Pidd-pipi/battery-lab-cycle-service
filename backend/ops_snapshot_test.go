package main

import (
	"context"
	"testing"
)

func TestReportRowIDsIndependent(t *testing.T) {
	svc := newOpsService([]OpsRecord{
		{ID: "rec-1", Subject: "cycle", Owner: "lab", Status: OpsStatusActive, Priority: OpsPriorityHigh, Labels: map[string]string{"site": "s1"}},
		{ID: "rec-2", Subject: "cycle", Owner: "lab", Status: OpsStatusActive, Priority: OpsPriorityHigh, Labels: map[string]string{"site": "s2"}},
	})
	rows, err := svc.Report(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		if len(row.IDs) != 1 {
			t.Fatalf("row %s/%s has %d ids want 1: %v", row.Site, row.Status, len(row.IDs), row.IDs)
		}
	}
}

func TestReportRowsIndependentAcrossCalls(t *testing.T) {
	first := buildReportRows([]OpsRecord{seedOpsRecord("rec-1")})
	if len(first) != 1 || first[0].Status != OpsStatusActive {
		t.Fatalf("first rows unexpected: %+v", first)
	}
	_ = buildReportRows([]OpsRecord{{ID: "rec-2", Subject: "cycle", Owner: "lab", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s1"}}})
	if first[0].Status != OpsStatusActive {
		t.Fatalf("previous rows mutated by later call: %+v", first)
	}
}
