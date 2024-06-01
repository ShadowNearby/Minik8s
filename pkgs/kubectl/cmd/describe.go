package cmd

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	} else if len(args) == 1 {
		kind := strings.ToLower(args[0])
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
	haveNamespace, ok := core.ObjTypeNamespace[objType]
	if !ok {
		fmt.Printf("wrong type %s", objType)
	}
	var resp string
	if haveNamespace {
		resp = utils.GetObject(objType, namespace, name)
	} else {
		resp = utils.GetObjectWONamespace(objType, name)
	}
	fmt.Println(resp)
}
