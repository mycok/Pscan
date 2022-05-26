/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mycok/Pscan/scan"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List all the hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")

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
