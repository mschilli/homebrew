package main

import (
	"os"
	"testing"
)

func TestRepo(t *testing.T) {
	gmf := NewGitmeta()
	f, err := os.Open("test/repos.gmf")
	if err != nil {
		panic(err)
	}
	gmf.AddGMF(f)

	clones := gmf.AllCloneables()

	wanted := 1
	if len(clones) != wanted {
		t.Log("Got ", len(clones), "but wanted", wanted)
		t.Fail()
	}

	urlWanted := "mschilli@somehost:git/bondy.git"
	if clones[0].URL != urlWanted {
		t.Log("Got URL", clones[0].URL, "but wanted", urlWanted)
		t.Fail()
	}

	dirWanted := "bondy"
	if clones[0].Dir != dirWanted {
		t.Log("Got dir", clones[0].Dir, "but wanted", dirWanted)
		t.Fail()
	}
}
