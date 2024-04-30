package cmd

import "github.com/spf13/cobra"

var describeCmd = &cobra.Command{
	Use:   "describe <resource> <name>/ describe <resource>s",
	Short: "Display one or many resources",
	Long:  "Display one or many resources",
	Args:  cobra.RangeArgs(1, 2),
	Run:   describeHandler,
}

func describeHandler(cmd *cobra.Command, args []string) {
	if len(args) == 1 {

	}
}
