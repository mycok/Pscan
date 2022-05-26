/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate bash completion for your command",
	Long: `To load your completions run source<(Pscan completion)

To load completions automatically on login, add this line to your
.bashrc file:
$ ~/.zshrc
source <(Pscan completion)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return completionAction(os.Stdout)
	},
}

func completionAction(w io.Writer) error {
	return rootCmd.GenZshCompletion(w)
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
