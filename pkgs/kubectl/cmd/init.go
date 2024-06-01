package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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

var namespace string
var filePath string

func init() {
	log.Info("Program started")
	log.SetLevel(log.InfoLevel)            // 设置日志级别
	log.SetFormatter(&log.TextFormatter{}) // 设置为文本格式
	RootCommand.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "kubectl (-n NAMESPACE)")
	applyCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "kubectl apply -f <FILENAME>")
	triggerCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "kubectl trigger <resource> <name> -f <FILENAME>")
	applyCmd.MarkFlagRequired("filePath")
	triggerCmd.MarkFlagRequired("filePath")
	RootCommand.AddCommand(applyCmd)
	RootCommand.AddCommand(deleteCmd)
	RootCommand.AddCommand(getCmd)
	RootCommand.AddCommand(describeCmd)
	RootCommand.AddCommand(triggerCmd)
}
