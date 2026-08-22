package main

import (
	"context"
	"testing"
)

func TestOpsClonePageNotAliased(t *testing.T) {
	page := OpsPage{Items: []OpsRecord{{ID: "rec-1"}, {ID: "rec-2"}}, Page: 1, PageSize: 25, Total: 2, HasNext: false}
	clone := opsClonePage(page)
	clone.Items[0].ID = "mutated"
	if page.Items[0].ID != "rec-1" {
		t.Fatal("clone mutated source page items")
	}
}

func TestOpsQueryDefaultsCapsPageSize(t *testing.T) {
	q := opsQueryDefaults(OpsQuery{Page: 1, PageSize: 999})
	if q.PageSize != 200 {
		t.Fatalf("page size=%d want 200", q.PageSize)
	}
}

func TestOpsPageOutOfRangeEmpty(t *testing.T) {
	svc := newOpsService([]OpsRecord{seedOpsRecord("rec-1")})
	page, err := svc.Search(context.Background(), OpsQuery{Page: 999, PageSize: 25})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 0 {
		t.Fatalf("expected empty page, got %d items", len(page.Items))
	}
}

func TestOpsPageZeroDefaults(t *testing.T) {
	svc := newOpsService([]OpsRecord{seedOpsRecord("rec-1")})
	page, err := svc.Search(context.Background(), OpsQuery{Page: 0, PageSize: 25})
	if err != nil {
		t.Fatal(err)
	}
	if page.Page != 1 {
		t.Fatalf("page=%d want 1", page.Page)
	}
	if len(page.Items) != 1 {
		t.Fatalf("items=%d want 1", len(page.Items))
	}
}
