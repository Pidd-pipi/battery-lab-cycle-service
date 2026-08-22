package main

import "testing"

func TestOpsZeroClockSafe(t *testing.T) {
	var clock OpsClock
	now := clock.Now()
	if now.IsZero() {
		t.Fatal("zero time returned")
	}
}
