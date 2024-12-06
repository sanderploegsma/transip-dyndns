package cmd

import "github.com/spf13/cobra"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Make sure all DNS records are up-to-date",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
