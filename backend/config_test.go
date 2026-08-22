package main

import "testing"

func TestLoadConfigHonorsPortEnv(t *testing.T) {
	t.Setenv("PORT", "18080")
	if got := LoadConfig().Port; got != "18080" {
		t.Fatalf("port=%q want 18080", got)
	}
}
