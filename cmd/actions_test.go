package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/mycok/Pscan/scan"
)

func TestHostActions(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}

	testCases := []struct{
		name string
		args []string
		expected string
		initList bool
		actionFunc func(io.Writer, string, []string) error
	}{
		{
			name:     "AddAction",
			args:     hosts,
			expected: "Added host: host1\nAdded host: host2\nAdded host: host3\n",
			initList: false,
			actionFunc: addAction,
		},
		{
			name:     "ListAction",
			args:     hosts,
			expected: "host1\nhost2\nhost3\n",
			initList: true,
			actionFunc: listAction,
		},
		{
			name:     "DeleteAction",
			args:     []string{"host1", "host2"},
			expected: "Deleted host: host1\nDeleted host: host2\n",
			initList: true,
			actionFunc: deleteAction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer

			fileName := setup(t, tc.args, tc.initList)
			defer t.Cleanup(func() {
				os.Remove(fileName)
			})

			if err := tc.actionFunc(&out, fileName, tc.args); err != nil {
				t.Fatalf("Expected no error but got: %q instead\n", err)
			}

			if tc.expected != out.String() {
				t.Errorf("Expected output: %q, but got: %q instead", tc.expected, out.String())
			}
		})
	}
}

func setup(t *testing.T, hosts []string, initList bool) string {
	// create a temp file.
	f, err := os.CreateTemp("", "Pscan")
	if err != nil {
		t.Fatal(err)
	}

	f.Close()

	// Initialize the list if necessary.
	if initList {
		hl := &scan.HostList{}

		for _, h := range hosts {
			if err := hl.Add(h); err != nil {
				t.Fatal(err)
			}
		}

		if err := hl.Save(f.Name()); err != nil {
			t.Fatal(err)
		}
	}

	return f.Name()
}