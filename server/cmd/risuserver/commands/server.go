package commands

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:          "server",
	Short:        "Run the Mattermost server",
	SilenceUsage: true,
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
