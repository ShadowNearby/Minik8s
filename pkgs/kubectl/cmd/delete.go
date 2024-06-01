package cmd

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <resource> <name>",
	Short: "Kubectl can delete resources and names",
	Long:  "Kubectl can delete resources and names",
	Run:   deleteHandler,
}

func deleteHandler(cmd *cobra.Command, args []string) {
	var kind string
	var name string
	var objType core.ObjType
	logrus.Debugln(args)
	if len(args) == 2 {
		kind = strings.ToLower(args[0])
		name = strings.ToLower(args[1])
		for _, ty := range core.ObjTypeAll {
			if !strings.Contains(ty, kind) {
				continue
			}
			objType = core.ObjType(ty)
		}
	} else {
		fmt.Printf("error: the server doesn't have a resource type %s\n", kind)
		return
	}
	err := utils.DeleteObject(objType, namespace, name)
	if err != nil {
		log.Error(err)
	}
}
