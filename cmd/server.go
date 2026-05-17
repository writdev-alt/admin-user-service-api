package cmd

import (
	"github.com/writdev-alt/admin-user-service/internal/api"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the web API server",
	Run: func(cmd *cobra.Command, args []string) {
		api.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
