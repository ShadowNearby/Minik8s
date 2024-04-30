package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wxnacy/wgo/arrays"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/controller"
	"strings"
)

var describeCmd = &cobra.Command{
	Use:   "describe <resource> <name>/ describe <resource>s",
	Short: "Display one or many resources",
	Long:  "Display one or many resources",
	Args:  cobra.RangeArgs(1, 2),
	Run:   describeHandler,
}

func describeHandler(cmd *cobra.Command, args []string) {
	var kind string
	if len(args) == 2 {
		kind = strings.ToLower(args[0])
		name := strings.ToLower(args[1])
		/* validate if `kind` is in the resource list */
		if idx := arrays.ContainsString(core.ObjTypeAll, kind); idx != -1 {
			objType := core.ObjType(kind + "s")
			res := controller.GetObject(objType, "", name)
			log.Infoln(res)
		} else if len(args) == 1 {
			kind = strings.ToLower(args[0])
			kind = kind[0 : len(kind)-1]
			/* validate if `kind` is in the resource list */
			if idx := arrays.ContainsString(core.ObjTypeAll, kind); idx != -1 {
				objType := core.ObjType(kind)
				res := controller.GetObject(objType, "", "")
				log.Infoln(res)
			}
		} else {
			fmt.Printf("error: the server doesn't have a resource type \"%s\"", kind)
		}
	}
}
