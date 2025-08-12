package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

const CacheFileName = ".photogrep-cache"

type cacheEntry struct {
	name string
}

func cachePath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	p := path.Join(u.HomeDir, CacheFileName)

	return p, nil
}

func cacheAsMap() (map[string]bool, error) {
	data := make(map[string]bool)

	cpath, err := cachePath()
	if err != nil {
		return data, err
	}

	raw, err := ioutil.ReadFile(cpath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return data, nil // ok if no cache exists
		}
		return data, err
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func cacheSave(data map[string]bool) error {
	cpath, err := cachePath()
	if err != nil {
		return err
	}

	text, err := json.Marshal(&data)

	err = ioutil.WriteFile(cpath, text, 0644)
	if err != nil {
		return err
	}

	return nil
}
