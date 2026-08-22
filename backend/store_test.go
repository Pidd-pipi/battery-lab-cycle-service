package main

import (
	"sync"
	"testing"
)

func TestCellStoreParallelListRead(t *testing.T) {
	store := NewCellStore()
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 150; j++ {
				_ = store.List()
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 150; j++ {
				_, _ = store.UpdateStatus("cell-a14", "paused")
				_, _ = store.UpdateStatus("cell-a14", "cycling")
			}
		}()
	}
	close(start)
	wg.Wait()
}

func TestCellStoreParallelStatusWritesStable(t *testing.T) {
	store := NewCellStore()
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 200; j++ {
				_, _ = store.UpdateStatus("cell-a14", "paused")
				_, _ = store.UpdateStatus("cell-a14", "cycling")
			}
		}()
	}
	close(start)
	wg.Wait()
	for _, c := range store.List() {
		if c.ID == "cell-a14" {
			if c.Status != "paused" && c.Status != "cycling" && c.Status != "complete" {
				t.Fatalf("invalid status %q", c.Status)
			}
		}
	}
}

func TestCellCompleteIncrementsCycleOnce(t *testing.T) {
	store := NewCellStore()
	cell, err := store.UpdateStatus("cell-a14", "complete")
	if err != nil {
		t.Fatal(err)
	}
	if cell.Cycle != 143 {
		t.Fatalf("cycle=%d want 143", cell.Cycle)
	}
	cell, err = store.UpdateStatus("cell-a14", "complete")
	if err != nil {
		t.Fatal(err)
	}
	if cell.Cycle != 143 {
		t.Fatalf("repeat complete cycle=%d want 143", cell.Cycle)
	}
}
