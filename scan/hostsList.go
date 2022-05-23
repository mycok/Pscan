// Package scan provides types and functions to perform TCP port
// scans on a list of hosts.
package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

var (
	ErrExists    = errors.New("host already in the list")
	ErrNotExists = errors.New("host not in the list")
)

// HostList represents a list of hosts upon which to perform a port scan.
type HostList struct {
	Hosts []string
}

// search searches for a host from the hosts list.
func (hl *HostList) search(host string) (bool, int) {
	sort.Strings(hl.Hosts)

	i := sort.SearchStrings(hl.Hosts, host)
	if i < len(hl.Hosts) && hl.Hosts[i] == host {
		return true, i
	}

	return false, -1
}

// Add appends a host to the list.
func (hl *HostList) Add(host string) error {
	if exists, _ := hl.search(host); exists {
		return fmt.Errorf("%w: %s", ErrExists, host)
	}

	hl.Hosts = append(hl.Hosts, host)

	return nil
}

// Remove deletes a host from the hosts list.
func (hl *HostList) Remove(host string) error {
	if exists, index := hl.search(host); exists {
		hl.Hosts = append(hl.Hosts[:index], hl.Hosts[index+1:]...)

		return nil
	}

	return fmt.Errorf("%w: %s", ErrNotExists, host)
}

// Load reads hosts from a hosts file.
func (hl *HostList) Load(hostsFile string) error {
	f, err := os.Open(hostsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts, scanner.Text())
	}

	return scanner.Err()
}

// Save saves the hosts to a hosts file.
func (hl *HostList) Save(hostsFile string) error {
	output := ""

	for _, h := range hl.Hosts {
		output += fmt.Sprintln(h)
	}

	return os.WriteFile(hostsFile, []byte(output), 0644)
}
