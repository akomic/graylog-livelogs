package commands

import (
	"cnvy/clusters"
	"github.com/spf13/cobra"
)

var (
	clustersCmd = &cobra.Command{
		Use:   "clusters",
		Short: "Cluster related actions",
		Long:  ``,

		Run: clustersRun,
	}
)

func clustersRun(ccmd *cobra.Command, args []string) {
	clusters.Clusters()
}
