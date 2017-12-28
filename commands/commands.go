package commands

import (
	"github.com/jcelliott/lumber"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config  string //
	daemon  bool   //
	version bool   //

	// CnvyCmd ...
	CnvyCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			// get the filepath
			abs, err := filepath.Abs("config.yml")
			if err != nil {
				lumber.Error("Error reading filepath: ", err.Error())
			}

			// get the config name
			base := filepath.Base(abs)

			// get the path
			path := filepath.Dir(abs)

			//
			viper.SetConfigName(strings.Split(base, ".")[0])
			viper.AddConfigPath(path)

			// Find and read the config file; Handle errors reading the config file
			if err := viper.ReadInConfig(); err != nil {
				lumber.Fatal("Failed to read config file: ", err.Error())
				os.Exit(1)
			}
		},

		Run: func(ccmd *cobra.Command, args []string) {
			ccmd.HelpFunc()(ccmd, args)
		},
	}
)

func init() {
	CnvyCmd.PersistentFlags().StringP("livelogs_url", "", "", "LiveLogs server url")
	CnvyCmd.PersistentFlags().StringP("livelogs_token", "", "", "LiveLogs token")
	viper.BindPFlag("livelogs_url", CnvyCmd.PersistentFlags().Lookup("livelogs_url"))
	viper.BindPFlag("livelogs_token", CnvyCmd.PersistentFlags().Lookup("livelogs_token"))

	CnvyCmd.AddCommand(livelogsCmd)
}
