package cmd

import "github.com/spf13/cobra"

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Kubectl apply can get apiObject in a declarative way",
	Long:  "Kubectl apply can get apiObject in a declarative way, usage kubectl apply -f [file]",
	Run:   applyHandler,
}
