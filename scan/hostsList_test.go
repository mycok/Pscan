package scan_test

import (
	"errors"
	"os"
	"testing"

	"github.com/mycok/Pscan/scan"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name        string
		host        string
		expected    int
		expectedErr error
	}{
		{
			name:        "AddNew",
			host:        "host2",
			expected:    2,
			expectedErr: nil,
		},
		{
			name:        "AddExisting",
			host:        "host1",
			expected:    1,
			expectedErr: scan.ErrExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostList{}

			// Initialize list
			if err := hl.Add("host1"); err != nil {
				t.Fatal(err)
			}

			err := hl.Add(tc.host)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("Expected error: %q, but got nil instead\n", tc.expectedErr)

					return
				}

				if !errors.Is(err, scan.ErrExists) {
					t.Errorf("Expected error: %q, but got %q instead\n", tc.expectedErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %q", err)

				return
			}

			if len(hl.Hosts) != tc.expected {
				t.Errorf("Expected hosts list to contain: %d hosts, but got %d hosts instead\n", tc.expected, len(hl.Hosts))
			}

			if hl.Hosts[1] != tc.host {
				t.Errorf(
					"Expected host name %s at index %d, but got %s at index %d instead \n",
					tc.host,
					1,
					hl.Hosts[1],
					1,
				)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	testCases := []struct {
		name        string
		host        string
		expected    int
		expectedErr error
	}{
		{
			name:        "RemoveExisting",
			host:        "host2",
			expected:    1,
			expectedErr: nil,
		},
		{
			name:        "RemoveNotFound",
			host:        "host3",
			expected:    1,
			expectedErr: scan.ErrNotExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostList{}

			// Initialize list
			for _, host := range []string{"host1", "host2"} {
				if err := hl.Add(host); err != nil {
					t.Fatal(err)
				}
			}

			err := hl.Remove(tc.host)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("Expected error: %q, but got nil instead\n", tc.expectedErr)

					return
				}

				if !errors.Is(err, scan.ErrNotExists) {
					t.Errorf("Expected error: %q, but got %q instead\n", tc.expectedErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %q", err)

				return
			}

			if len(hl.Hosts) != tc.expected {
				t.Errorf("Expected hosts list to contain: %d hosts, but got %d hosts instead\n", tc.expected, len(hl.Hosts))
			}

			if hl.Hosts[0] == tc.host {
				t.Errorf("Host name %s should not be in the hosts list", tc.host)
			}
		})
	}
}

func TestLoadNoFile(t *testing.T) {
	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}

	if err := os.Remove(tf.Name()); err != nil {
		t.Fatalf("Error removing temp file: %s", err)
	}

	hl := &scan.HostList{}

	if err := hl.Load(tf.Name()); err != nil {
		t.Errorf("Expected no error, but got %q instead\n", err)

		if !errors.Is(err, scan.ErrNotExists) {
			t.Errorf("Expected error: %q, but got %q instead\n", scan.ErrNotExists, err)
		}
	}
}
