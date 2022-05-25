/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// hostsCmd represents the hosts command
var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Manage the hosts list",
	Long: `Manage the hosts list for Pscan

Add hosts with the [add] command
Delete hosts with the [delete] command
List hosts with the [list] command`,
}

func init() {
	rootCmd.AddCommand(hostsCmd)
}
