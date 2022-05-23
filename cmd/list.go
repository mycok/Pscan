/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/mycok/Pscan/scan"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List all the hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}

		return listAction(os.Stdout, hostsFile, args)
	},
}

func listAction(w io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	if len(hl.Hosts) > 0 {
		for _, h := range hl.Hosts {
			if _, err := fmt.Fprintln(w, h); err != nil {
				return err
			}
		}
	} else {
		_, err := fmt.Fprintln(w, "No hosts found")

		return err
	}

	return nil
}

func init() {
	hostsCmd.AddCommand(listCmd)
}
