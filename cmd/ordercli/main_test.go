package main

import "testing"

func TestRunHelp(t *testing.T) {
	if code := run([]string{"--help"}); code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestRunBadArgs(t *testing.T) {
	if code := run([]string{"definitely-not-a-command"}); code == 0 {
		t.Fatalf("expected non-zero exit code")
	}
}
