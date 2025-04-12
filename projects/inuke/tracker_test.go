// ///////////////////////////////////////
// tracker_test.go - Test tracker
// Mike Schilli, 2025 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"testing"
)

func TestTracker(t *testing.T) {
	tr := NewTracker(3)
	tr.MoveLeft()

	if tr.Current() != 3 {
		t.Log("Expected current at", 3, " but got ", tr.Current())
		t.Fail()
	}

	tr.MoveRight()
	if tr.Current() != 1 {
		t.Log("Expected current at", 1, " but got ", tr.Current())
		t.Fail()
	}

	tr.MoveLeft()
	tr.MoveRight()
	tr.RemoveLeft()

	if tr.Total() != 2 {
		t.Log("Expected total at", 2, " but got ", tr.Total())
		t.Fail()
	}
	if tr.Current() != 1 {
		t.Log("Expected total at", 1, " but got ", tr.Current())
		t.Fail()
	}
}
