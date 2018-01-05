package commands

import (
	"cnvy/apps"
	"github.com/spf13/cobra"
)

var (
	appsCmd = &cobra.Command{
		Use:   "apps",
		Short: "App related actions",
		Long:  ``,

		Run: appsRun,
	}
)

func appsRun(ccmd *cobra.Command, args []string) {
	apps.ListApps()
}
