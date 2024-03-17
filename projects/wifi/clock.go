// ///////////////////////////////////////
// clock.go - timer widget for monitor
// Mike Schilli, 2023 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"time"
)

func clock(arg ...string) chan string {
	ch := make(chan string)
	start := time.Now()

	go func() {
		for {
			z := time.Unix(0, 0).UTC()
			ch <- z.Add(time.Since(start)).Format("15:04:05")
			time.Sleep(1 * time.Second)
		}
	}()

	return ch
}
