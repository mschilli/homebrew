// ///////////////////////////////////////
// eth.go - ifconfig widget for monitor
// Mike Schilli, 2023 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"net"
	"sort"
	"strings"
	"time"
)

func nifs(arg ...string) chan string {
	ch := make(chan string)

	go func() {
		for {
			eths, err := ifconfig()

			if err != nil {
				ch <- err.Error()
				time.Sleep(10 * time.Second)
				continue
			}

			ch <- strings.Join(eths, ", ")
			time.Sleep(10 * time.Second)
		}
	}()

	return ch
}

func ifconfig() ([]string, error) {
	var list []string

	ifaces, err := net.Interfaces()
	if err != nil {
		return list, err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return list, err
		}
		if len(addrs) == 0 {
			continue
		}
		for _, addr := range addrs {
			ip := strings.Split(addr.String(), "/")[0]
			if net.ParseIP(ip).To4() != nil {
				list = append(list, iface.Name+" "+ip)
			}
		}
	}
	sort.Strings(list)
	return list, nil
}
