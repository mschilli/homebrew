// ///////////////////////////////////////
// tracker.go - Current photo tracker
// Mike Schilli, 2025 (m@perlmeister.com)
// ///////////////////////////////////////
package main

type Tracker struct {
	total      int
	currentIdx int
}

func NewTracker(total int) *Tracker {
	return &Tracker{
		total:      total,
		currentIdx: 0,
	}
}

func (t *Tracker) Current() int {
	return t.currentIdx + 1
}

func (t *Tracker) Total() int {
	return t.total
}

func (t *Tracker) InsertLeft() {
	t.total += 1
	t.currentIdx += 1
}

func (t *Tracker) MoveLeft() {
	if t.currentIdx == 0 {
		t.currentIdx = t.total - 1
	} else {
		t.currentIdx -= 1
	}
}

func (t *Tracker) MoveRight() {
	if t.currentIdx == t.total-1 {
		t.currentIdx = 0
	} else {
		t.currentIdx += 1
	}
}

func (t *Tracker) RemoveCurrent() {
	if t.total == 0 {
		panic("Can't remove element from zero length array")
	}
	t.total -= 1
	if t.currentIdx > t.total-1 {
		t.MoveLeft()
	}
	t.MoveRight()
}
