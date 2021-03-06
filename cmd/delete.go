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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete <host1>...<hostn>",
	Aliases: []string{"d"},
	Short:   "Delete host(s) to the hosts list",
	Long: `Delete any number of hosts to the hosts list. Do this by providing
	a comma separated list of host names.`,
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")

		return deleteAction(os.Stdout, hostsFile, args)
	},
}

func deleteAction(w io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	for _, host := range args {
		if err := hl.Remove(host); err != nil {
			return err
		}

		if _, err := fmt.Fprintln(w, "Deleted host:", host); err != nil {
			return err
		}
	}

	return hl.Save(hostsFile)
}

func init() {
	hostsCmd.AddCommand(deleteCmd)
}
