/* 	Package scan provides types and functions to perform TCP port
scans on a list of hosts.
*/
package scan

import (
	"fmt"
	"net"
	"time"
)

// Result represents the scan result for a single host.
type Result struct {
	Host       string
	Found      bool
	PortStates []PortState
}

// PortState represents the state of a single TCP port.
type PortState struct {
	Port int
	Open state
}

type state bool

func (s state) String() string {
	if s {
		return "open"
	}

	return "closed"
}

func scanPort(host string, port int) PortState {
	p := PortState{Port: port}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return p
	}

	p.Open = true

	conn.Close()

	return p
}

// Run performs a port scan on the hosts list.
func Run(hl *HostList, ports []int) []Result {
	results := make([]Result, 0, len(hl.Hosts))

	for _, h := range hl.Hosts {
		r := Result{Host: h}

		if _, err := net.LookupHost(h); err != nil {
			results = append(results, r)

			continue
		}

		r.Found = true

		for _, p := range ports {
			r.PortStates = append(r.PortStates, scanPort(h, p))
		}

		results = append(results, r)
	}

	return results
}
