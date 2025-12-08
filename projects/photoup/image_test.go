package main

import (
	"testing"
)

func TestResize(t *testing.T) {
	w, h := rescaleDims(2000, 1500, 1000)
	if w != 1000 || h != 750 {
		t.Errorf("w=%d h=%d", w, h)
	}

	w, h = rescaleDims(1500, 2000, 1000)
	if w != 750 || h != 1000 {
		t.Errorf("w=%d h=%d", w, h)
	}

	w, h = rescaleDims(1000, 1000, 1000)
	if w != 1000 || h != 1000 {
		t.Errorf("w=%d h=%d", w, h)
	}

	w, h = rescaleDims(500, 500, 1000)
	if w != 500 || h != 500 {
		t.Errorf("w=%d h=%d", w, h)
	}
}
