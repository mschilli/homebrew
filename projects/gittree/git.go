// ///////////////////////////////////////
// git.go - git repo utility functions
// Mike Schilli, 2025 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"fmt"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func gitTopDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func gitDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func gitCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func gitHashFileLocal() (string, error) {
	dir, err := gitDir()
	if err != nil {
		return "", err
	}
	branch, err := gitCurrentBranch()
	if err != nil {
		return "", err
	}

	return path.Join(dir, "refs", "heads", branch), nil
}

func gitHashFileRemote() (string, error) {
	dir, err := gitDir()
	if err != nil {
		return "", err
	}
	return path.Join(dir, "refs", "remotes", "origin"), nil
}

type FileStatus struct {
	File   string
	Status string
}

func gitStatus() ([]FileStatus, error) {
	cmd := exec.Command("git", "status", "--porcelain", ".")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var statuses []FileStatus
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	blankRe := regexp.MustCompile(`\s+`)

	for _, line := range lines {
		parts := blankRe.Split(strings.TrimSpace(line), 2)
		if len(parts) != 2 {
			continue
		}
		statuses = append(statuses, FileStatus{
			Status: parts[0],
			File:   parts[1],
		})
	}

	return statuses, nil
}

func gitCommitCount(head string) (int64, error) {
	cmd := exec.Command("git", "rev-list", "--count", head)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	numstr := strings.TrimSpace(string(out))
	v, err := strconv.ParseInt(numstr, 10, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func gitPushStatus() (string, error) {
	counts := []int64{}
	branch, err := gitCurrentBranch()
	if err != nil {
		return "", err
	}
	heads := []string{"HEAD", "origin/" + branch}

	for _, head := range heads {
		count, err := gitCommitCount(head)
		if err != nil {
			continue
		}

		counts = append(counts, count)
	}

	result := "Up-to-date with remote."

	direction := "ahead"
	diff := counts[0] - counts[1]
	absdiff := diff
	if diff < 0 {
		direction = "behind"
		absdiff = -diff
	}

	plural := ""
	if absdiff > 1 {
		plural = "s"
	}

	if diff != 0 {
		result = fmt.Sprintf("[red]Local branch %s by %d commit%s.", direction, absdiff, plural)
	}

	return result, nil
}
