package clusters

import (
	"fmt"
	"os"
	// "reflect"
	"strings"

	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/awserr"
	// "github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/ryanuber/columnize"
	"github.com/spf13/viper"
)

func Clusters() {
	verbose := viper.GetBool("verbose")
	// fmt.Println("Clusters")

	ecssvc := ecs.New(session.New())
	params := &ecs.ListClustersInput{}

	resp, err := ecssvc.ListClusters(params)

	if err != nil {
		// A service error occurred.
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(2)
	} else if err != nil {
		// A non-service error occurred.
		panic(err)
	}

	if verbose {
		params := &ecs.DescribeClustersInput{
			Clusters: resp.ClusterArns,
		}
		resp, err := ecssvc.DescribeClusters(params)

		if err != nil {
			// A service error occurred.
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(2)
		} else if err != nil {
			// A non-service error occurred.
			panic(err)
		}

		// Pretty-print the response data.
		// fmt.Println(awsutil.StringValue(resp))

		output := []string{"Name | Instances | Running Tasks | Pending Tasks | Status"}

		for _, cluster := range resp.Clusters {
			// fmt.Println(*cluster.ClusterName)
			output = append(output, fmt.Sprintf("%s | %d | %d | %d | %s", *cluster.ClusterName, *cluster.RegisteredContainerInstancesCount, *cluster.RunningTasksCount, *cluster.PendingTasksCount, *cluster.Status))
		}
		result := columnize.SimpleFormat(output)
		fmt.Println(result)
		// fmt.Println("Verbose? a?")
	} else {
		for _, clusterARN := range resp.ClusterArns {
			fmt.Println(strings.Split(*clusterARN, "/")[1])
		}
	}
}
