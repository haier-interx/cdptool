package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func main() {

	var configFile string

	var cmdCheck = &cobra.Command{
		Use:   "check [check configure file is valid]",
		Short: "Check configure file syntax",
		Long:  "check the configure file is valid",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", strings.Join(args, " "))
		},
	}

	cmdCheck.Flags().StringVarP(&configFile, "conf", "c", "config.yml", "times to echo the input")
}
