// ///////////////////////////////////////
// tree.go - tview terminal ui functions
// Mike Schilli, 2025 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"strings"

	"github.com/rivo/tview"
)

var statusColor = map[string]string{
	"??": "[orange]",
	"M":  "[red]",
	"MM": "[red]",
}

type Cmd struct {
	fs      []FileStatus
	pstatus string
}

func ui(topDir string) (*tview.Application, chan Cmd) {
	app := tview.NewApplication()

	root := mktree(topDir, []FileStatus{})
	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	pstatus := tview.NewTextView().SetDynamicColors(true)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tree, 0, 1, true).
		AddItem(pstatus, 1, 0, false)

	cmds := make(chan Cmd)

	go func() {
		for {
			cmd := <-cmds
			pstatus.SetText(cmd.pstatus)
			newroot := mktree(topDir, cmd.fs)
			tree.SetRoot(newroot).SetCurrentNode(newroot)
			app.QueueUpdateDraw(func() {})
		}
	}()

	app.SetRoot(layout, true)

	return app, cmds
}

func mktree(title string, entries []FileStatus) *tview.TreeNode {
	root := tview.NewTreeNode(title)
	nodeMap := map[string]*tview.TreeNode{"": root}

	for _, entry := range entries {
		parts := strings.Split(entry.File, "/")
		path := ""
		parent := root

		for i, part := range parts {
			if i > 0 {
				path += "/"
			}
			path += part

			if _, exists := nodeMap[path]; !exists {
				color := ""
				if i == len(parts)-1 {
					color = statusColor[entry.Status]
				}
				node := tview.NewTreeNode(color + part)
				nodeMap[path] = node

				parent.AddChild(node)
			}

			parent = nodeMap[path]
		}
	}

	return root
}
