package main

import (
	"os"
	"os/exec"
	"path"
)

func cloneOrUpdate(c Cloneable, gitPath string) error {
	fullPath := path.Join(gitPath, c.Dir)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", c.URL, fullPath)

		_, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
	}

	err = os.Chdir(fullPath)
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "pull")
	_, err = cmd.CombinedOutput()
	return err
}
