package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cnvy/accounts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	// "reflect"
)

var (
	config  string //
	daemon  bool   //
	version bool   //
	Acc     accounts.Account

	// CnvyCmd ...
	CnvyCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			// get the filepath
			abs, err := filepath.Abs(filepath.Join(os.Getenv("HOME"), ".cnvy/config.yml"))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading filepath: ", err.Error())
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
				fmt.Fprintln(os.Stderr, "Failed to read config file: ", err.Error())
				os.Exit(1)
			}

			// fmt.Println(accounts.Acc.Name)
			// Acc = accounts.Login()
			if accounts.Acc.Name == "" {
				accounts.Acc.SUAccount("tf-idea-ecs")
			}
		},

		Run: func(ccmd *cobra.Command, args []string) {
			ccmd.HelpFunc()(ccmd, args)
		},
	}
)

func init() {
	CnvyCmd.PersistentFlags().StringP("livelogs_url", "", "", "LiveLogs server url")
	CnvyCmd.PersistentFlags().MarkHidden("livelogs_url")
	CnvyCmd.PersistentFlags().StringP("livelogs_token", "", "", "LiveLogs token")
	CnvyCmd.PersistentFlags().MarkHidden("livelogs_token")

	CnvyCmd.PersistentFlags().StringP("consul_address", "", "", "Consul ")
	CnvyCmd.PersistentFlags().MarkHidden("consul_address")
	CnvyCmd.PersistentFlags().StringP("consul_scheme", "", "", "Consul ")
	CnvyCmd.PersistentFlags().MarkHidden("consul_scheme")
	CnvyCmd.PersistentFlags().StringP("consul_user", "", "", "Consul ")
	CnvyCmd.PersistentFlags().MarkHidden("consul_user")
	CnvyCmd.PersistentFlags().StringP("consul_pass", "", "", "Consul ")
	CnvyCmd.PersistentFlags().MarkHidden("consul_pass")
	CnvyCmd.PersistentFlags().StringP("consul_token", "", "", "Consul ")
	CnvyCmd.PersistentFlags().MarkHidden("consul_token")

	CnvyCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbosity")
	viper.BindPFlag("livelogs_url", CnvyCmd.PersistentFlags().Lookup("livelogs_url"))
	viper.BindPFlag("livelogs_token", CnvyCmd.PersistentFlags().Lookup("livelogs_token"))

	viper.BindPFlag("consul_address", CnvyCmd.PersistentFlags().Lookup("consul_address"))
	viper.BindPFlag("consul_scheme", CnvyCmd.PersistentFlags().Lookup("consul_scheme"))
	viper.BindPFlag("consul_user", CnvyCmd.PersistentFlags().Lookup("consul_user"))
	viper.BindPFlag("consul_pass", CnvyCmd.PersistentFlags().Lookup("consul_pass"))
	viper.BindPFlag("consul_token", CnvyCmd.PersistentFlags().Lookup("consul_token"))

	viper.BindPFlag("verbose", CnvyCmd.PersistentFlags().Lookup("verbose"))

	CnvyCmd.AddCommand(accountsCmd)
	CnvyCmd.AddCommand(livelogsCmd)
	CnvyCmd.AddCommand(clustersCmd)
	CnvyCmd.AddCommand(appsCmd)
}
