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
}
