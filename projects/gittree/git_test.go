package main

import (
	"fmt"
	"testing"
)

func TestStatus(t *testing.T) {
	statuses, err := gitStatus()
	if err != nil {
		panic(err)
	}
	for _, status := range statuses {
		fmt.Printf("file: %s status: %s\n", status.File, status.Status)
	}
}

func TestPush(t *testing.T) {
	status, err := gitPushStatus()
	if err != nil {
		panic(err)
	}
	fmt.Printf("status=%s\n", status)
}
