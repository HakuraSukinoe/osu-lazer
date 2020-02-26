package commands

import (
	"github.com/mattermost/viper"
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
	RootCmd.PersistentFlags().StringP("config", "c", "config.json", "Configuration file to use.")
	RootCmd.PersistentFlags().Bool("disableconfigwatch", false, "When set config.json will not be loaded from disk when the file is changed.")
	RootCmd.PersistentFlags().Bool("platform", false, "This flag signifies that the user tried to start the command from the platform binary, so we can log a mssage")
	RootCmd.PersistentFlags().MarkHidden("platform")

	viper.SetEnvPrefix("mm")
	viper.BindEnv("config")
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
}
