package commands

import (
	"github.com/spf13/cobra"
)

type Command = cobra.Command

func Run(args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.Execute()
}

var RootCmd = &cobra.Command{
	Use:   "risuserver",
	Short: "Open source, self-hosted Slack-alternative",
	Long: `
          _.====.._                                 ___                _ 
       ,:.         ~-_                             / _ \   ___  _  _  | |
          \'\        ~-_                          | (_) | (_-< | || | |_|
            |          \'.                         \___/  /__/  \_,_| (_)
          ,/             ~-_                                
-..__..-''                 ~~--..__...----...--.....---.....--....---...

      Rhythm is just a click away. Open source game server for osu!lazer

                                     GitHub: github.com/deissh/osu-lazer                                                    
                                                       2019-2020, deissh`,
}

func init() {
	RootCmd.PersistentFlags().StringP("config", "c", "config.yaml", "Configuration file to use.")
}
