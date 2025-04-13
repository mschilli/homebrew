// ///////////////////////////////////////
// inuke.go - GUI to sort out bad images
// Mike Schilli, 2021 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"container/list"
	"os"

	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/hashicorp/golang-lru"
	"path"
	"path/filepath"
)

var Cache *lru.Cache

const Version = "1.61"

func main() {
	version := flag.Bool("version", false, "print version info")
	offset := flag.Int("offset", 0, "start with image at idx")

	flag.Parse()

	if *version {
		fmt.Printf("%s %s\n", path.Base(os.Args[0]), Version)
		return
	}

	win := app.New().NewWindow(versionInfo())

	var err error

	undo := undoStackNew()
	undo.Clear()

	Cache, err = lru.New(128)
	panicOnErr(err)

	cwd, err := os.Getwd()
	panicOnErr(err)

	dir, err := storage.ListerForURI(
		storage.NewFileURI(cwd))
	panicOnErr(err)

	files, err := dir.List()
	panicOnErr(err)

	images := list.New()

	for _, file := range files {
		if isImage(file) {
			images.PushBack(file)
		}
	}

	whacker := whackerNew(images)
	whacker.WhackerStart() // start whacker in background

	if images.Len() == 0 {
		panic("No images found.")
	}

	cur := images.Front()
	if *offset != 0 {
	    for i := 0; i < *offset; i++ {
		cur = scrollRight(images, cur)
	    }
	}

	img := canvas.NewImageFromResource(nil)
	img.SetMinSize(
		fyne.NewSize(DspWidth, DspHeight))
	lbl := widget.NewLabel(
		fmt.Sprintf("%s  [H]-Left [L]-Right [D]elete " +
		    "[U]ndo [S]tash [Q]uit", versionInfo()))
	con := container.NewVBox(img, lbl)
	win.SetContent(con)

	showURI(win, cur.Value.(fyne.URI), 1, images.Len())
	showImage(img, cur.Value.(fyne.URI))
	preloadImage(scrollRight(images,
		cur).Value.(fyne.URI))

	tr := NewTracker(images.Len())

	win.Canvas().SetOnTypedKey(
		func(ev *fyne.KeyEvent) {
			key := string(ev.Name)
			switch key {
			case "L":
				cur = scrollRight(images, cur)
				tr.MoveRight()
			case "H":
				cur = scrollLeft(images, cur)
				tr.MoveLeft()
			case "D":
				if images.Len() == 1 {
					panic("Not enough images!!")
				}
				old := cur
				cur = scrollRight(images, cur)
				tr.MoveRight()
				uri := old.Value.(fyne.URI)

				undo.Push(uri)

				toTrash(old.Value.(fyne.URI))
				images.Remove(old)
				tr.RemoveLeft()
			case "S":
				toStash(cur.Value.(fyne.URI))
			case "U":
				uri := undo.Pop()
				if uri == nil {
					return
				}
				images.InsertBefore(uri, cur)
				tr.InsertLeft()
				fromTrash(uri)
				cur = scrollLeft(images, cur)
			case "Q":
				os.Exit(0)
			}
			showURI(win, cur.Value.(fyne.URI),
			    tr.Current(), images.Len())
			showImage(img,
				cur.Value.(fyne.URI))
			preloadImage(scrollRight(images,
				cur).Value.(fyne.URI))
		})

	win.ShowAndRun()
}

func versionInfo() string {
	return fmt.Sprintf("iNuke v%s", Version)
}

func showURI(win fyne.Window, uri fyne.URI, current int, length int) {
	win.SetTitle(fmt.Sprintf("%s (%d of %d)",
	    filepath.Base(uri.String()), current, length))
}

func scrollRight(l *list.List,
	e *list.Element) *list.Element {
	e = e.Next()
	if e == nil {
		e = l.Front()
	}
	return e
}

func scrollLeft(l *list.List,
	e *list.Element) *list.Element {
	e = e.Prev()
	if e == nil {
		e = l.Back()
	}
	return e
}

type undoEle struct {
	URI fyne.URI
}

type undoStack struct {
	stack *list.List
}

func undoStackNew() *undoStack {
	l := undoStack{stack: list.New()}
	return &l
}

func (u *undoStack) Clear() {
	u.stack = u.stack.Init()
}

func (u *undoStack) Push(uri fyne.URI) {
	ue := undoEle{URI: uri}
	u.stack.PushBack(ue)
}

func (u *undoStack) Pop() fyne.URI {
	if u.stack.Len() == 0 {
		return nil
	}

	ue := u.stack.Back()
	u.stack.Remove(ue)

	uev := ue.Value.(undoEle)
	return uev.URI
}

type whackStack struct {
	images *list.List
}

func whackerNew(l *list.List) *whackStack {
	w := whackStack{images: l}
	return &w
}

func (w *whackStack) WhackerStart() {
	// Try to preload all images one by one, and whack those that can't be loaded,
	// so inuke won't even display them later
	if w.images.Len() == 0 {
	    return
	}
	go func() {
		e := w.images.Front()
		for {
			_, err := loadImage(e.Value.(fyne.URI))
			old := e
			e = e.Next()
			if err != nil {
				fmt.Printf("Whacking %s\n", old.Value.(fyne.URI))
				w.images.Remove(old)
			}
			if e == nil {
				return // end of list
			}
		}
	}()
}
