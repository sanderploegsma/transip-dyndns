package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Periodically update DNS records",
	Run: func(cmd *cobra.Command, args []string) {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		slog.Info("Server started")

		c := cron.New()
		schedule := viper.GetString("server.schedule")
		slog.Info("Scheduling periodic update", "schedule", schedule)
		if _, err := c.AddFunc(schedule, func() { fmt.Println("Hello cron") }); err != nil {
			slog.Error("Error while scheduling periodic update", "error", err)
			os.Exit(1)
		}

		slog.Info("Starting scheduler")
		c.Start()

		<-sigs
		slog.Info("Stopping scheduler")
		c.Stop()
		slog.Info("Server stopped")
	},
}

func init() {
	serverCmd.PersistentFlags().String("schedule", "@hourly", "Cron schedule")
	viper.BindPFlag("server.schedule", serverCmd.Flags().Lookup("schedule"))
}
