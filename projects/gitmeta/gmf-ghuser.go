package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Repository struct {
	Name    string `json:"name"`
	Fork    bool   `json:"fork"`
	SSHURL  string `json:"ssh_url"`
	HTMLURL string `json:html_url"`
}

type GMFGithubUser struct{}

func NewGMFGithubUser() GMFGithubUser {
	return GMFGithubUser{}
}

func (g GMFGithubUser) Applicable(e GitMetaEntry) bool {
	return e.Type == "github-user-repos"
}

func (g GMFGithubUser) Expand(e GitMetaEntry) ([]Cloneable, error) {
	clones := []Cloneable{}

	res, err := githubProjectsOfUser(e.User)
	if err != nil {
		return clones, err
	}

	for _, repo := range res {
		clones = append(clones, Cloneable{URL: repo.SSHURL, Dir: repo.Name})
	}

	return clones, nil
}

func githubProjectsOfUser(user string) ([]Repository, error) {
	results := []Repository{}

	url := "https://api.github.com/users/" + user + "/repos"

	resp, err := http.Get(url)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return results, fmt.Errorf("Status %s", resp.Status)
	}

	var repos []Repository

	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return results, err
	}

	for _, repo := range repos {
		if !repo.Fork {
			results = append(results, repo)
		}
	}

	return results, err
}
