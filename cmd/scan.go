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

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the hosts",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")

		ports, err := cmd.Flags().GetIntSlice("ports")
		if err != nil {
			return err
		}

		return scanAction(os.Stdout, hostsFile, ports)
	},
}

func scanAction(w io.Writer, hostsFile string, ports []int) error {
	hl := &scan.HostList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	results := scan.Run(hl, ports)

	return printResults(w, results)
}

func printResults(w io.Writer, results []scan.Result) error {
	message := ""

	for _, r := range results {
		message += fmt.Sprintf("%s:", r.Host)

		if !r.Found {
			message += " Host not found\n"

			continue
		}

		message += fmt.Sprintln()

		for _, ps := range r.PortStates {
			message += fmt.Sprintf("\t%d: %s\n", ps.Port, ps.Open)
		}

		message += fmt.Sprintln()
	}

	_, err := fmt.Fprintln(w, message)

	return err
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().IntSliceP("ports", "p", []int{22, 80, 443}, "ports to scan")
}
