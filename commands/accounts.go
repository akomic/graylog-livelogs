package commands

import (
	"cnvy/accounts"
	"github.com/spf13/cobra"
)

var (
	accountsCmd = &cobra.Command{
		Use:   "accounts",
		Short: "Handle Accounts",
		Long:  ``,

		Run: accountsRun,
	}
)

var (
	accountName string
)

func init() {
	accountsCmd.Flags().StringVarP(&accountName, "su", "s", "", "Login to account")
}

func accountsRun(ccmd *cobra.Command, args []string) {
	if accountName == "" {
		accounts.Acc.ListAccounts()
	} else {
		accounts.Acc.SUAccount(accountName)
	}
}
