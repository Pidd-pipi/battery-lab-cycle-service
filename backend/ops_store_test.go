package main

import (
	"context"
	"sync"
	"testing"
)

func seedOpsRecord(id string) OpsRecord {
	return OpsRecord{ID: id, Subject: "battery cycle", Owner: "lab", Status: OpsStatusActive, Priority: OpsPriorityHigh, Revision: 1, Labels: map[string]string{"site": "s1"}}
}

func TestOpsStoreGetReturnsIsolatedCopy(t *testing.T) {
	s := newOpsStore([]OpsRecord{seedOpsRecord("rec-1")})
	got, err := s.Get(context.Background(), "rec-1")
	if err != nil {
		t.Fatal(err)
	}
	got.Labels["site"] = "hacked"
	again, err := s.Get(context.Background(), "rec-1")
	if err != nil {
		t.Fatal(err)
	}
	if again.Labels["site"] != "s1" {
		t.Fatalf("store mutated through returned record: site=%q", again.Labels["site"])
	}
}

func TestOpsStoreListReturnsIsolatedCopies(t *testing.T) {
	s := newOpsStore([]OpsRecord{seedOpsRecord("rec-1")})
	items, err := s.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	items[0].Labels["site"] = "hacked"
	again, err := s.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if again[0].Labels["site"] != "s1" {
		t.Fatalf("store mutated through listed record: site=%q", again[0].Labels["site"])
	}
}

func TestOpsStoreUpdateKeepsCallerLabelsIsolated(t *testing.T) {
	s := newOpsStore(nil)
	rec := seedOpsRecord("rec-1")
	if err := s.Put(context.Background(), rec); err != nil {
		t.Fatal(err)
	}
	upd := seedOpsRecord("rec-1")
	upd.Status = OpsStatusPaused
	upd.Labels["room"] = "r2"
	if err := s.Update(context.Background(), upd, 1); err != nil {
		t.Fatal(err)
	}
	upd.Labels["room"] = "hacked"
	got, err := s.Get(context.Background(), "rec-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Labels["room"] != "r2" {
		t.Fatalf("store aliased caller labels: room=%q", got.Labels["room"])
	}
}

func TestOpsStoreConcurrentCountListStable(t *testing.T) {
	s := newOpsStore([]OpsRecord{seedOpsRecord("rec-1"), seedOpsRecord("rec-2")})
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 150; j++ {
				_ = s.Count()
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 150; j++ {
				_ = s.Update(context.Background(), seedOpsRecord("rec-1"), 0)
			}
		}()
	}
	close(start)
	wg.Wait()
}
