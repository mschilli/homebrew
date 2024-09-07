package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type GMFSSH struct{}

func NewGMFSSH() GMFSSH {
	return GMFSSH{}
}

func (g GMFSSH) Applicable(e GitMetaEntry) bool {
	return e.Type == "ssh-user-repos"
}

func (g GMFSSH) Expand(e GitMetaEntry) ([]Cloneable, error) {
	clones := []Cloneable{}

	res, err := sshGitDirs(e.Host, e.Path)
	if err != nil {
		return clones, err
	}

	for _, dir := range res {
		ldir := strings.TrimSuffix(dir, ".git")
		clones = append(clones, Cloneable{
			URL: fmt.Sprintf("%s:%s/%s", e.Host, e.Path, dir),
			Dir: ldir,
		})
	}

	return clones, nil
}

func sshGitDirs(host string, path string) ([]string, error) {
	results := []string{}

	cmd := exec.Command("ssh", host, "ls", path+"/*/config")

	output, err := cmd.Output()
	if err != nil {
		return results, err
	}

	reader := bytes.NewReader(output)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "/")
		if len(parts) > 2 {
			results = append(results, parts[len(parts)-2])
		}
	}

	return results, err
}
