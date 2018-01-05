package main

import (
	"fmt"
	"os"

	"cnvy/commands"
)

//
func main() {
	if err := commands.CnvyCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
