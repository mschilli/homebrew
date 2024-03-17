// ///////////////////////////////////////
// ping.go - ping icmp widget for monitor
// Mike Schilli, 2023 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"fmt"
	"github.com/prometheus-community/pro-bing"
	"time"
)

func ping(addr ...string) chan string {
	ch := make(chan string)
	firstTime := true

	go func() {
		for {
			pinger, err := probing.NewPinger(addr[0])
			pinger.Timeout, _ = time.ParseDuration("10s")

			if err != nil {
				ch <- err.Error()
				time.Sleep(10 * time.Second)
				continue
			}

			if firstTime {
				ch <- "Pinging ..."
				firstTime = false
			}

			pinger.Count = 3
			err = pinger.Run()
			if err != nil {
				ch <- err.Error()
				time.Sleep(10 * time.Second)
				continue
			}

			stats := pinger.Statistics()
			ch <- fmt.Sprintf("%v ", stats.Rtts)
			time.Sleep(10 * time.Second)
		}
	}()

	return ch
}
