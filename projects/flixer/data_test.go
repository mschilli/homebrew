package main

import "testing"

func TestVisiblePickIndexesHidesUpcomingByDefault(t *testing.T) {
	picks := []Pick{
		{Title: "Ready"},
		{Title: "Later (upcoming)"},
		{Title: "Also Ready"},
	}

	got := VisiblePickIndexes(picks, false)
	want := []int{0, 2}

	if len(got) != len(want) {
		t.Fatalf("VisiblePickIndexes() = %v, want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("VisiblePickIndexes() = %v, want %v", got, want)
		}
	}
}

func TestVisiblePickIndexesIncludesUpcomingWhenRequested(t *testing.T) {
	picks := []Pick{
		{Title: "Ready"},
		{Title: "Later (upcoming)"},
	}

	got := VisiblePickIndexes(picks, true)
	want := []int{0, 1}

	if len(got) != len(want) {
		t.Fatalf("VisiblePickIndexes() = %v, want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("VisiblePickIndexes() = %v, want %v", got, want)
		}
	}
}

func TestUpcomingPickIndexes(t *testing.T) {
	picks := []Pick{
		{Title: "Ready"},
		{Title: "Later (upcoming)"},
		{Title: "Also Later (upcoming)"},
	}

	got := UpcomingPickIndexes(picks)
	want := []int{1, 2}

	if len(got) != len(want) {
		t.Fatalf("UpcomingPickIndexes() = %v, want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("UpcomingPickIndexes() = %v, want %v", got, want)
		}
	}
}
