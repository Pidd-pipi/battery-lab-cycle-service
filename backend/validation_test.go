package main

import "testing"

func TestValidateCellStatusRejectsUnknown(t *testing.T) {
	if err := ValidateCellStatus("failed"); err == nil {
		t.Fatal("expected error for unknown status")
	}
	if err := ValidateCellStatus("complete"); err != nil {
		t.Fatalf("complete should be valid: %v", err)
	}
}
