package main

import "testing"

func TestOpsGuardResumeAllowed(t *testing.T) {
	g := newOpsGuard()
	rec := seedOpsRecord("rec-1")
	rec.Status = OpsStatusPaused
	if !g.ResumeAllowed(rec) {
		t.Fatal("paused record should be resumable")
	}
}

func TestOpsGuardCloseAllowed(t *testing.T) {
	g := newOpsGuard()
	rec := seedOpsRecord("rec-1")
	rec.Status = OpsStatusActive
	if !g.CloseAllowed(rec) {
		t.Fatal("active record should be closable")
	}
}

func TestOpsGuardAllowsLegalTransition(t *testing.T) {
	g := newOpsGuard()
	if !g.CanTransition(OpsStatusActive, OpsStatusClosed) {
		t.Fatal("active to closed should be legal")
	}
}
