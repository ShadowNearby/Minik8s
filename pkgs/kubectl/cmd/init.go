package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

/*
kubectl [command] [TYPE] [NAME] [flags]
*/
var RootCommand = &cobra.Command{
	Use:   "kubectl",
	Short: "Kubectl is a tool for controlling minik8s cluster.",
	Long:  `Kubectl is a tool for controlling minik8s cluster. To see the help of a specific command, use: kubectl [command] --help`,
	Run:   runRoot,
}

func Execute() {
	if err := RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func runRoot(cmd *cobra.Command, args []string) {
	fmt.Printf("execute %s args:%v \n", cmd.Name(), args)
	fmt.Println("kubectl is for better control of minik8s")
	fmt.Println(cmd.UsageString())
}

var NameSpace string
var filePath string

func init() {
	RootCommand.PersistentFlags().StringVarP(&NameSpace, "namespace", "n", "default", "kubectl (-n NAMESPACE)")
	applyCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "kubectl apply -f <FILENAME>")
	applyCmd.MarkFlagRequired("filePath")
	RootCommand.AddCommand(applyCmd)
	RootCommand.AddCommand(deleteCmd)
	RootCommand.AddCommand(getCmd)
	RootCommand.AddCommand(describeCmd)
}
