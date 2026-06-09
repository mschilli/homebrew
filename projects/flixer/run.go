/////////////////////////////////////////
// run.go - Run external programs
// Mike Schilli, 2026 (m@perlmeister.com)
/////////////////////////////////////////
package main

import (
  "os"
  "os/exec"
)

func runVim(path string) {
  cmd := exec.Command("vim", path)

  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  err := cmd.Run()
  if err != nil {
    panic(err)
  }
}

func runFirefox(url string) {
  cmd := exec.Command("/bin/sh", "firefox", url)
  if err := cmd.Start(); err != nil {
    panic(err)
  }
}
