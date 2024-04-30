package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or many resources",
	Long:  "Display one or many resources",
	Run:   getHandler,
}

func getHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("resource name is required")
		cmd.Usage()
		return
	}
	kind := args[0]
	if len(args) == 1 {
		if kind == core.NodeKind {
			getNodes()
			return
		}
	}
}

func getNodes() {
	url := utils.ParseUrlMany(core.NodeKind, "nil")
	statusCode, bodyJson, err := utils.GetByTarget(url)
	if err != nil {
		log.Fatal(err)
	}
	if statusCode != 200 {
		log.Fatal("status code is :", statusCode)
	} else {
		log.Infoln(bodyJson)
	}

}
