// /////////////////////////////////////////
// util.go - Utility Functions
// 2024, Mike Schilli, m@perlmeister.com
// /////////////////////////////////////////
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path"
)

func photoDir() (string, error) {
	randomStr := make([]byte, 32)
	if _, err := rand.Read(randomStr); err != nil {
		return "", err
	}
	hash := sha256.Sum256(randomStr)
	shaDir := hex.EncodeToString(hash[:])[0:30]

	dir := "ph"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}

		err := os.WriteFile(path.Join(dir, ".htaccess"), []byte("Options -Indexes\n"), 0644)
		if err != nil {
			return "", err
		}
	}

	path := path.Join(dir, shaDir)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return "", err
	}

	return path, nil
}
