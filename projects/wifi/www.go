// ///////////////////////////////////////
// www.go - http fetcher widget
// Mike Schilli, 2023 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"fmt"
	"net/http"
	"time"
)

func httpGet(arg ...string) chan string {
	ch := make(chan string)
	firstTime := true

	go func() {
		for {
			if firstTime {
				ch <- "Fetching ..."
				firstTime = false
			}

			now := time.Now()

			_, err := http.Get(arg[0])
			if err != nil {
				ch <- err.Error()
				time.Sleep(10 * time.Second)
				continue
			}

			dur := time.Since(now)
			ch <- fmt.Sprintf("%.3f OK ", dur.Seconds())
			time.Sleep(10 * time.Second)
		}
	}()

	return ch
}
