package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "transip-dyndns",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		for key, value := range viper.GetViper().AllSettings() {
			fmt.Printf("%s=%v\n", key, value)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(updateCmd)

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose logging")
	viper.BindPFlag("verbose", rootCmd.Flags().Lookup("verbose"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
