package apps

import (
	"cnvy/accounts"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/ryanuber/columnize"
	"github.com/spf13/viper"
	"os"
	"sort"
	"strings"
)

var (
	ecssvc *ecs.ECS
)

func listServices() []*string {
	serviceArns := []*string{}

	var nextToken *string
	nextToken = nil

	list := func() {
		params := &ecs.ListServicesInput{
			Cluster:   &accounts.Acc.Cluster_name,
			NextToken: nextToken,
		}
		resp, err := ecssvc.ListServices(params)
		if err != nil {
			// A service error occurred.
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(2)
		}
		serviceArns = append(serviceArns, resp.ServiceArns...)
		nextToken = resp.NextToken
	}

	for {
		list()

		if nextToken == nil {
			break
		}
	}

	return serviceArns
}

func describeServices(serviceArns []*string) []*ecs.Service {
	serviceDescs := []*ecs.Service{}

	process := func(s []*string) {
		params := &ecs.DescribeServicesInput{
			Cluster:  &accounts.Acc.Cluster_name,
			Services: s,
		}

		resp, err := ecssvc.DescribeServices(params)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(2)
		}
		serviceDescs = append(serviceDescs, resp.Services...)
	}

	for i := 0; i < len(serviceArns); i += 10 {
		end := i + 10
		if end > len(serviceArns) {
			end = len(serviceArns)
		}
		process(serviceArns[i:end])
	}
	return serviceDescs
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func sortServices(serviceDescs []*ecs.Service) ([]string, map[string][]*ecs.Service) {
	sServiceDescs := map[string][]*ecs.Service{}
	list := []string{}
	for _, s := range serviceDescs {
		nameParts := strings.Split(*s.ServiceName, "-")

		stackName := "Other"
		if len(nameParts) >= 2 {
			stackName = nameParts[1]
		}
		sServiceDescs[stackName] = append(sServiceDescs[stackName], s)
		if stringInSlice(stackName, list) == false {
			list = append(list, stackName)
		}
	}
	sort.Strings(list)
	return list, sServiceDescs
}

func sortServiceArns(serviceArns []*string) ([]string, map[string][]string) {
	sServiceArns := map[string][]string{}
	list := []string{}
	for _, s := range serviceArns {
		nameParts := strings.Split(*s, "/")

		serviceName := *s
		if len(nameParts) >= 2 {
			serviceName = nameParts[1]
		}

		stackNameParts := strings.Split(serviceName, "-")
		stackName := "Other"
		if len(stackNameParts) >= 2 {
			stackName = stackNameParts[1]
		}

		sServiceArns[stackName] = append(sServiceArns[stackName], serviceName)
		if stringInSlice(stackName, list) == false {
			list = append(list, stackName)
		}
	}
	sort.Strings(list)
	return list, sServiceArns
}

func ListApps() {
	verbose := viper.GetBool("verbose")
	fmt.Println("Apps on Cluster", accounts.Acc.Cluster_name)

	ecssvc = ecs.New(session.New())
	serviceArns := listServices()

	if verbose {
		serviceDescs := describeServices(serviceArns)

		sortList, sServiceDescs := sortServices(serviceDescs)

		for _, stackName := range sortList {
			fmt.Println()
			fmt.Println("STACKNAME:", stackName, "-", "Apps:", len(sServiceDescs[stackName]))

			output := []string{}
			if verbose {
				output = []string{"Name | desiredCount | runningCount | LB Count | status | Equilibrium | CreatedAt | TaskDefinition"}

				for _, s := range sServiceDescs[stackName] {
					output = append(output, fmt.Sprintf("%s | %d | %d | %d | %s | %t | %s | %s", *s.ServiceName, *s.DesiredCount, *s.RunningCount, len(s.LoadBalancers), *s.Status, (*s.DesiredCount == *s.RunningCount), *s.CreatedAt, *s.TaskDefinition))
				}
			} else {
				output = []string{"Name | desiredCount | runningCount | LB Count | status | Equilibrium"}

				for _, s := range sServiceDescs[stackName] {
					output = append(output, fmt.Sprintf("%s | %d | %d | %d | %s | %t", *s.ServiceName, *s.DesiredCount, *s.RunningCount, len(s.LoadBalancers), *s.Status, (*s.DesiredCount == *s.RunningCount)))
				}
			}
			result := columnize.SimpleFormat(output)
			fmt.Println(result)
		}
	} else {
		sortList, sServiceArns := sortServiceArns(serviceArns)

		for _, stackName := range sortList {
			fmt.Println(stackName)
			for _, s := range sServiceArns[stackName] {
				fmt.Println("\t", s)
			}
		}
	}
}
