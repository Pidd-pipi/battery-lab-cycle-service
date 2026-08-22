package main

import (
	"context"
	"sync"
	"testing"
)

func TestOpsCacheGetReturnsIsolatedCopy(t *testing.T) {
	c := newOpsCache()
	c.Store("rec-1", seedOpsRecord("rec-1"))
	got, ok := c.Get(context.Background(), "rec-1")
	if !ok {
		t.Fatal("missing")
	}
	got.Labels["site"] = "hacked"
	again, _ := c.Get(context.Background(), "rec-1")
	if again.Labels["site"] != "s1" {
		t.Fatalf("cache mutated through returned record: %q", again.Labels["site"])
	}
}

func TestOpsCacheStoreKeepsCallerLabelsIsolated(t *testing.T) {
	c := newOpsCache()
	rec := seedOpsRecord("rec-1")
	c.Store("rec-1", rec)
	rec.Labels["site"] = "hacked"
	got, _ := c.Get(context.Background(), "rec-1")
	if got.Labels["site"] != "s1" {
		t.Fatalf("cache aliased caller labels: %q", got.Labels["site"])
	}
}

func TestOpsCacheConcurrentGetStoreStable(t *testing.T) {
	c := newOpsCache()
	c.Store("rec-1", seedOpsRecord("rec-1"))
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 150; j++ {
				_, _ = c.Get(context.Background(), "rec-1")
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 150; j++ {
				c.Store("rec-1", seedOpsRecord("rec-1"))
			}
		}()
	}
	close(start)
	wg.Wait()
}
