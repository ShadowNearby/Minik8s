package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

/*
kubectl [command] [TYPE] [NAME] [flags]
*/
var commands = &cobra.Command{
	Use:   "kubectl",
	Short: "Kubectl is a tool for controlling minik8s cluster.",
	Long:  `Kubectl is a tool for controlling minik8s cluster. To see the help of a specific command, use: kubectl [command] --help`,
	Run:   runRoot,
}

func runRoot(cmd *cobra.Command, args []string) {
	fmt.Printf("execute %s args:%v \n", cmd.Name(), args)
	fmt.Println("kubectl is for better control of minik8s")
	fmt.Println(cmd.UsageString())

}
func Execute() {
	if err := commands.Execute(); err != nil {
		fmt.Println(err)
	}
}
func init() {
	commands.AddCommand(applyCmd)
}
