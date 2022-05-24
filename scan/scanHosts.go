/* 	Package scan provides types and functions to perform TCP port
scans on a list of hosts.
*/
package scan

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single TCP port.
type PortState struct {
	Port int
	Open state

}

func scanPort(host string, port int) PortState {
	p := PortState{ Port: port }

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return p
	}

	conn.Close()
	p.Open = true

	return p
}

type state bool

func (s state) String() string {
	if s {
		return "open"
	}

	return "closed"
}

// Result represents the scan result for a single host.
type Result struct {
	Host string
	Found bool
	PortStates []PortState
}

// Run performs a port scan on the hosts list.
func Run(hl *HostList, ports []int) []Result {
	res := make([]Result, 0, len(hl.Hosts))

	for _, h := range hl.Hosts {
		r := Result{Host: h}

		if _, err := net.LookupHost(h); err != nil {
			r.Found = false
			res = append(res, r)

			continue
		}

		for _, p := range ports {
			r.PortStates = append(r.PortStates, scanPort(h, p))
		}

		res = append(res, r)
	}

	return res
}