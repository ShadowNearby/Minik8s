package cmd

import (
	"fmt"
	"minik8s/config"
	"os"

	"github.com/sirupsen/logrus"
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
var cfgFile string
var update string

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	RootCommand.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "kubectl (-n NAMESPACE)")
	RootCommand.PersistentFlags().StringVar(&cfgFile, "config", "./config/config.json", "config file (default is ./config/config.json)")
	applyCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "kubectl apply -f <FILENAME>")
	applyCmd.Flags().StringVarP(&update, "update", "u", "", "update (default is false)")
	applyCmd.MarkFlagRequired("filePath")
	if err := config.InitConfig(cfgFile); err != nil {
		logrus.Fatalf("Error initializing config: %s", err.Error())
	}
	RootCommand.AddCommand(applyCmd)
	RootCommand.AddCommand(deleteCmd)
	RootCommand.AddCommand(getCmd)
	RootCommand.AddCommand(describeCmd)
}
