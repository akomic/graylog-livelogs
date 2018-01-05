package commands

import (
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"cnvy/accounts"
	"cnvy/logs"
	"fmt"
	"os"
)

var (
	livelogsCmd = &cobra.Command{
		Use:   "livelogs",
		Short: "App Live Logs",
		Long:  ``,

		Run: livelogs,
	}
)

var (
	cluster string
)

func init() {
	livelogsCmd.Flags().StringVarP(&cluster, "cluster", "c", "", "Container Cluster Name")
	livelogsCmd.Flags().StringSliceP("filter", "f", nil, "Filter e.g. -f stack_name=idea1")
	livelogsCmd.Flags().BoolP("rawOutput", "r", false, "Dump complete messages as json")
	viper.BindPFlag("filter", livelogsCmd.Flags().Lookup("filter"))
	viper.BindPFlag("rawOutput", livelogsCmd.Flags().Lookup("rawOutput"))
}

func livelogs(ccmd *cobra.Command, args []string) {
	accounts.Acc.Verify()
	filters := viper.GetStringSlice("filter")

	if accounts.Acc.Livelogs_address == "" {
		fmt.Fprintln(os.Stderr, "Error: No sure where livelogs are.")
		os.Exit(2)
	}

	u := url.URL{Scheme: "ws", Host: accounts.Acc.Livelogs_address, Path: "/filter"}
	querySlice := []string{}

	if accounts.Acc.Livelogs_token != "" {
		querySlice = append(querySlice, "token="+accounts.Acc.Livelogs_token)
	}

	querySlice = append(querySlice, filters...)
	u.RawQuery = strings.Join(querySlice, "&")

	logs.Livelogs(u.String())
}
