package main

import "testing"

func TestChangeStatusReturnedCycleMatchesStored(t *testing.T) {
	svc := NewCycleService(NewCellStore())
	cell, err := svc.ChangeStatus("cell-a14", "complete")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range svc.Cells() {
		if c.ID == "cell-a14" && c.Cycle != cell.Cycle {
			t.Fatalf("stored cycle=%d returned=%d", c.Cycle, cell.Cycle)
		}
	}
}
