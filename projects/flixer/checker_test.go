package main

import "testing"

func TestFindJustWatchCandidateMatchesTitleLine(t *testing.T) {
	href, ok := findJustWatchCandidate("the matrix", []justWatchCandidate{
		{Title: "The Matrix Reloaded\n2003", Href: "https://www.justwatch.com/us/movie/the-matrix-reloaded"},
		{Title: "The Matrix\n1999", Href: "https://www.justwatch.com/us/movie/the-matrix"},
	})

	if !ok {
		t.Fatal("findJustWatchCandidate() did not find exact title")
	}
	if href != "https://www.justwatch.com/us/movie/the-matrix" {
		t.Fatalf("findJustWatchCandidate() = %q", href)
	}
}

func TestFindJustWatchCandidateRejectsSimilarTitle(t *testing.T) {
	_, ok := findJustWatchCandidate("the matrix", []justWatchCandidate{
		{Title: "The Matrix Reloaded\n2003", Href: "https://www.justwatch.com/us/movie/the-matrix-reloaded"},
	})

	if ok {
		t.Fatal("findJustWatchCandidate() found similar title as exact match")
	}
}

func TestJustWatchHasNetflixOffer(t *testing.T) {
	if !justWatchHasNetflixOffer("Currently you are able to watch streaming on Netflix") {
		t.Fatal("justWatchHasNetflixOffer() did not detect Netflix")
	}
}

func TestJustWatchSearchURL(t *testing.T) {
	got := JustWatchSearchURL("The Matrix (upcoming)")
	want := "https://www.justwatch.com/us/search?q=the+matrix"

	if got != want {
		t.Fatalf("JustWatchSearchURL() = %q, want %q", got, want)
	}
}
