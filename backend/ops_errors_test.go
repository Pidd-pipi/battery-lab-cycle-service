package main

import (
	"errors"
	"testing"
)

func TestOpsErrorWrapChainPreserved(t *testing.T) {
	err := wrapOps("create", "store.put", ErrOpsConflict)
	if !errors.Is(err, ErrOpsConflict) {
		t.Fatalf("wrapped error lost sentinel: %v", err)
	}
	var typed *OpsError
	if !errors.As(err, &typed) {
		t.Fatal("expected *OpsError")
	}
	if typed.Code != "create" {
		t.Fatalf("code=%q want create", typed.Code)
	}
}

func TestOpsErrorCodeClassification(t *testing.T) {
	cases := []struct {
		err  error
		want string
	}{
		{wrapOps("get", "store.get", ErrOpsNotFound), "not_found"},
		{wrapOps("update", "store.update", ErrOpsConflict), "conflict"},
		{wrapOps("create", "store.put", ErrOpsInvalid), "invalid"},
		{wrapOps("move", "state.move", ErrOpsTransition), "transition"},
		{wrapOps("create", "policy.check", ErrOpsPolicy), "policy"},
	}
	for _, tc := range cases {
		if got := opsCode(tc.err); got != tc.want {
			t.Fatalf("opsCode(%v)=%q want %q", tc.err, got, tc.want)
		}
	}
}

func TestOpsIsHelpersDetectWrappedSentinels(t *testing.T) {
	if !opsIsConflict(wrapOps("update", "store.update", ErrOpsConflict)) {
		t.Fatal("opsIsConflict missed wrapped conflict")
	}
	if !opsIsNotFound(wrapOps("get", "store.get", ErrOpsNotFound)) {
		t.Fatal("opsIsNotFound missed wrapped not found")
	}
}

func TestOpsCodeFallsBackToOperationCode(t *testing.T) {
	err := &OpsError{Code: "create", Operation: "store.put", Cause: errors.New("boom")}
	if got := opsCode(err); got != "create" {
		t.Fatalf("opsCode=%q want create", got)
	}
}
