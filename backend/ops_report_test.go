package main

import (
	"context"
	"testing"
)

func TestReportRowsDistinctKeys(t *testing.T) {
	svc := newOpsService([]OpsRecord{
		{ID: "rec-1", Subject: "cycle", Owner: "lab", Status: OpsStatusActive, Priority: OpsPriorityHigh, Labels: map[string]string{"site": "s1"}},
		{ID: "rec-2", Subject: "cycle", Owner: "lab", Status: OpsStatusActive, Priority: OpsPriorityHigh, Labels: map[string]string{"site": "s2"}},
	})
	rows, err := svc.Report(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, row := range rows {
		key := row.Site + "|" + string(row.Status)
		if seen[key] {
			t.Fatalf("duplicate report row for %s", key)
		}
		seen[key] = true
	}
}

func TestReportRowsCountedOnce(t *testing.T) {
	svc := newOpsService([]OpsRecord{
		{ID: "rec-1", Subject: "cycle", Owner: "lab", Status: OpsStatusActive, Priority: OpsPriorityHigh, Labels: map[string]string{"site": "s1"}},
		{ID: "rec-2", Subject: "cycle", Owner: "lab", Status: OpsStatusActive, Priority: OpsPriorityHigh, Labels: map[string]string{"site": "s2"}},
	})
	rows, err := svc.Report(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows=%d want 2", len(rows))
	}
	for _, row := range rows {
		if row.Count != 1 {
			t.Fatalf("count=%d want 1 for %s", row.Count, row.Site)
		}
	}
}
