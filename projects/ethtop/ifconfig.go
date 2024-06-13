///////////////////////////////////////////
// ifconfig - Query host's network adapters
// 2018, Mike Schilli, m@perlmeister.com
///////////////////////////////////////////
package main

import (
	"fmt"
	"net"
	"sort"
	"strings"
)

func NifsAsStrings() []string {
	var list []string

	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		network := fmt.Sprintf("%10s",
			iface.Name)
		addrs, _ := iface.Addrs()
		if len(addrs) == 0 {
			continue
		}
		for _, addr := range addrs {
			split := strings.Split(addr.String(), "/")
			a := split[0]
			if net.ParseIP(a).To4() != nil {
				network += " " + a
				list = append(list, network)
			}
		}
	}
	sort.Strings(list)
	return list
}
