package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/mycok/Pscan/scan"
)

// TODO: add an integration test for a compiled tool.
func TestHostActions(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}

	testCases := []struct {
		name       string
		args       []string
		expected   string
		initList   bool
		actionFunc func(io.Writer, string, []string) error
	}{
		{
			name:       "AddAction",
			args:       hosts,
			expected:   "Added host: host1\nAdded host: host2\nAdded host: host3\n",
			initList:   false,
			actionFunc: addAction,
		},
		{
			name:       "ListAction",
			args:       hosts,
			expected:   "host1\nhost2\nhost3\n",
			initList:   true,
			actionFunc: listAction,
		},
		{
			name:       "DeleteAction",
			args:       []string{"host1", "host2"},
			expected:   "Deleted host: host1\nDeleted host: host2\n",
			initList:   true,
			actionFunc: deleteAction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer

			fileName := setup(t, tc.args, tc.initList)
			defer t.Cleanup(func() {
				os.Remove(fileName)
			})

			if err := tc.actionFunc(&outputBuf, fileName, tc.args); err != nil {
				t.Fatalf("Expected no error but got: %q instead\n", err)
			}

			if tc.expected != outputBuf.String() {
				t.Errorf("Expected output: %q, but got: %q instead", tc.expected, outputBuf.String())
			}
		})
	}
}

func TestToolIntegration(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}

	fileName := setup(t, hosts, false)
	defer t.Cleanup(func() {
		os.Remove(fileName)
	})

	var outputBuf bytes.Buffer

	hostToDelete := "host2"

	hostsAfterDel := []string{"host1", "host3"}

	// Define the combined final output from all the executed operations.
	expectedOutput := ""

	for _, h := range hosts {
		expectedOutput += fmt.Sprintf("Added host: %s\n", h)
	}

	expectedOutput += strings.Join(hosts, "\n")
	expectedOutput += fmt.Sprintln()
	expectedOutput += fmt.Sprintf("Deleted host: %s\n", hostToDelete)
	expectedOutput += strings.Join(hostsAfterDel, "\n")
	expectedOutput += fmt.Sprintln()

	for _, h := range hostsAfterDel {
		expectedOutput += fmt.Sprintf("%s: Host not found", h)
		expectedOutput += fmt.Sprintln()
	}

	expectedOutput += fmt.Sprintln()

	// Add hosts to the list.
	if err := addAction(&outputBuf, fileName, hosts); err != nil {
		t.Fatalf("Expected no error but got: %q instead\n", err)
	}

	// List all hosts.
	if err := listAction(&outputBuf, fileName, nil); err != nil {
		t.Fatalf("Expected no error but got: %q instead\n", err)
	}

	// Delete a host.
	if err := deleteAction(&outputBuf, fileName, []string{hostToDelete}); err != nil {
		t.Fatalf("Expected no error but got: %q instead\n", err)
	}

	// List remaining hosts after a deletion operation.
	if err := listAction(&outputBuf, fileName, nil); err != nil {
		t.Fatalf("Expected no error but got: %q instead\n", err)
	}

	// Perform a port scan on a list of hosts.
	if err := scanAction(&outputBuf, fileName, nil); err != nil {
		t.Fatalf("Expected no error but got: %q instead\n", err)
	}

	if expectedOutput != outputBuf.String() {
		t.Errorf("Expected output: %q, but got: %q instead\n", expectedOutput, outputBuf.String())
	}

}

func TestScanAction(t *testing.T) {
	hosts := []string{"localhost", "389.389.389.389"}

	fileName := setup(t, hosts, true)
	defer t.Cleanup(func() {
		os.Remove(fileName)
	})

	ports := []int{}

	// Initialize two ports, 1 open and 1 closed
	for i := 0; i < 2; i++ {
		listener, err := net.Listen("tcp", net.JoinHostPort("localhost", "0"))
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

		if i == 1 {
			listener.Close()
		}
	}

	// Define the expected output.
	expectedOutput := fmt.Sprintln("localhost:")
	expectedOutput += fmt.Sprintf("\t%d: open\n", ports[0])
	expectedOutput += fmt.Sprintf("\t%d: closed\n", ports[1])
	expectedOutput += fmt.Sprintln()
	expectedOutput += fmt.Sprintln("389.389.389.389: Host not found")
	expectedOutput += fmt.Sprintln()

	var outputBuf bytes.Buffer

	if err := scanAction(&outputBuf, fileName, ports); err != nil {
		t.Fatalf("Expected nil error, but got: %q\n", err)
	}

	if expectedOutput != outputBuf.String() {
		t.Errorf("Expected output %q, but got: %q instead", expectedOutput, outputBuf.String())
	}
}

func setup(t *testing.T, hosts []string, initList bool) string {
	// create a temp file.
	file, err := os.CreateTemp("", "Pscan")
	if err != nil {
		t.Fatal(err)
	}

	file.Close()

	// Initialize the list if necessary.
	if initList {
		hl := &scan.HostList{}

		for _, h := range hosts {
			if err := hl.Add(h); err != nil {
				t.Fatal(err)
			}
		}

		if err := hl.Save(file.Name()); err != nil {
			t.Fatal(err)
		}
	}

	return file.Name()
}
