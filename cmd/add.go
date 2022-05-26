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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add <host1>....<hostn>",
	Aliases: []string{"a"},
	Short:   "Add new host(s) to the hosts list",
	Long: `Add any number of hosts to the hosts list. Do this by providing
a space / comma separated list of host names.`,
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")

		return addAction(os.Stdout, hostsFile, args)
	},
}

func addAction(w io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	for _, host := range args {
		if err := hl.Add(host); err != nil {
			return err
		}

		if _, err := fmt.Fprintln(w, "Added host:", host); err != nil {
			return err
		}
	}

	return hl.Save(hostsFile)
}

func init() {
	hostsCmd.AddCommand(addCmd)
}
