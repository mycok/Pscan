package scan_test

import (
	"net"
	"strconv"
	"testing"

	"github.com/mycok/Pscan/scan"
)

func TestStateCustomString(t *testing.T) {
	ps := scan.PortState{}

	if ps.Open.String() != "closed" {
		t.Errorf("Expected: %q, but got: %q instead\n", "closed", ps.Open.String())
	}

	ps.Open = true

	if ps.Open.String() != "open" {
		t.Errorf("Expected: %q, but got: %q instead\n", "open", ps.Open.String())
	}
}

func TestRunWithHostFound(t *testing.T) {
	testCases := []struct {
		name          string
		expectedState string
	}{
		{
			name:          "OpenPort",
			expectedState: "open",
		},
		{
			name:          "ClosedPort",
			expectedState: "closed",
		},
	}

	host := "localhost"

	hl := &scan.HostList{}

	hl.Add(host)

	ports := []int{}

	for _, tc := range testCases {
		listener, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}

		defer listener.Close()

		_, portStr, err := net.SplitHostPort(listener.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		ports = append(ports, port)

		if tc.name == "ClosedPort" {
			listener.Close()
		}
	}

	res := scan.Run(hl, ports)

	if len(res) != 1 {
		t.Errorf("Expected a single result but got: %q instead\n", len(res))
	}

	if res[0].Host != host {
		t.Errorf("Expected host as: %s but got: %s instead\n", host, res[0].Host)
	}

	if !res[0].Found {
		t.Errorf("Expected host: %s : TO BE found", host)
	}

	if len(res[0].PortStates) != 2 {
		t.Errorf("Expected 2 port states, but got: %q instead\n", len(res[0].PortStates))
	}

	for i, tc := range testCases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("Expected port: %d, but got: %d instead\n", ports[i], res[0].PortStates[i].Port)
		}

		if res[0].PortStates[i].Open.String() != tc.expectedState {
			t.Errorf("Expected port %d to be %s\n", ports[i], tc.expectedState)
		}
	}
}

func TestHostNotFound(t *testing.T) {
	host := "389.389.389.389"

	hl := &scan.HostList{}
	hl.Add(host)

	res := scan.Run(hl, []int{})

	if len(res) != 1 {
		t.Errorf("Expected a single result but got: %q instead\n", len(res))
	}

	if res[0].Host != host {
		t.Errorf("Expected host as: %s but got: %s instead\n", host, res[0].Host)
	}

	if res[0].Found {
		t.Errorf("Expected host: %s NOT to be found", host)
	}

	if len(res[0].PortStates) != 0 {
		t.Errorf("Expected 0 port states, but got: %q instead\n", len(res[0].PortStates))
	}
}
