package main

import "testing"

func TestOpsStateResumePausedAllowed(t *testing.T) {
	m := newOpsStateMachine()
	if err := m.Move(OpsStatusPaused, OpsStatusActive, "resume"); err != nil {
		t.Fatalf("resume from paused rejected: %v", err)
	}
}

func TestOpsStateCloseActiveAllowed(t *testing.T) {
	m := newOpsStateMachine()
	if err := m.Move(OpsStatusActive, OpsStatusClosed, "close"); err != nil {
		t.Fatalf("close from active rejected: %v", err)
	}
}

func TestOpsStatusValidIncludesPaused(t *testing.T) {
	if !opsStatusValid(OpsStatusPaused) {
		t.Fatal("paused should be a valid status")
	}
}

func TestOpsStatusTerminalExcludesPaused(t *testing.T) {
	if opsStatusTerminal(OpsStatusPaused) {
		t.Fatal("paused must not be terminal")
	}
}
