package cmd

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <resource> <name>/ get <resource>s",
	Short: "Display one or many resources",
	Long:  "Display one or many resources",
	Args:  cobra.RangeArgs(1, 2),
	Run:   getHandler,
}

func getHandler(cmd *cobra.Command, args []string) {
	var kind string
	var name string
	var objType core.ObjType
	logrus.Debugln(args)
	wrongType := true
	if len(args) == 2 {
		kind = strings.ToLower(args[0])
		name = strings.ToLower(args[1])
		for _, ty := range core.ObjTypeAll {
			if !strings.Contains(ty, kind) {
				continue
			}
			objType = core.ObjType(ty)
			wrongType = false
		}
	} else if len(args) == 1 {
		kind := strings.ToLower(args[0])
		for _, ty := range core.ObjTypeAll {
			if !strings.Contains(ty, kind) {
				continue
			}
			objType = core.ObjType(ty)
			wrongType = false
		}
	} else {
		fmt.Printf("error: the server doesn't have a resource type %s\n", kind)
		return
	}

	if wrongType {
		fmt.Printf("wrong type")
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

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	switch objType {
	case core.ObjPod:
		t.AppendHeader(table.Row{"NAME", "STATUS", "AGE", "POD-IP", "HOST"})
		pods := []core.Pod{}
		if resp != "" {
			if name == "" {
				err := utils.JsonUnMarshal(resp, &pods)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
			} else {
				pod := core.Pod{}
				err := utils.JsonUnMarshal(resp, &pod)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
				pods = append(pods, pod)
			}
			for _, pod := range pods {
				t.AppendRow(table.Row{pod.MetaData.Name, pod.Status.Condition, time.Since(pod.Status.StartTime).Round(time.Second), pod.Status.PodIP, pod.Status.HostIP})
			}
		}
	case core.ObjNode:
		t.AppendHeader(table.Row{"NAME", "IP", "STATUS"})
		nodes := []core.Node{}
		if resp != "" {
			if name == "" {
				err := utils.JsonUnMarshal(resp, &nodes)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
			} else {
				node := core.Node{}
				err := utils.JsonUnMarshal(resp, &node)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
				nodes = append(nodes, node)
			}
			for _, node := range nodes {
				t.AppendRow(table.Row{node.NodeMetaData.Name, node.Spec.NodeIP, node.Status.Phase})
			}
		}
	case core.ObjReplicaSet:
		t.AppendHeader(table.Row{"NAME", "DESIRED", "READY"})
		rslist := []core.ReplicaSet{}
		if resp != "" {
			if name == "" {
				err := utils.JsonUnMarshal(resp, &rslist)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
			} else {
				rs := core.ReplicaSet{}
				err := utils.JsonUnMarshal(resp, &rs)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
				rslist = append(rslist, rs)
			}
			for _, rs := range rslist {
				t.AppendRow(table.Row{rs.MetaData.Name, rs.Spec.Replicas, rs.Status.RealReplicas})
			}
		}
	case core.ObjService:
		t.AppendHeader(table.Row{"NAME", "TYPE", "SELECTOR", "CLUSTER-IP", "PORT(S)"})
		services := []core.Service{}
		if resp != "" {
			if name == "" {
				err := utils.JsonUnMarshal(resp, &services)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
			} else {
				service := core.Service{}
				err := utils.JsonUnMarshal(resp, &service)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
				services = append(services, service)
			}
			for _, service := range services {
				portStrs := []string{}
				if service.Spec.Type == core.ServiceTypeNodePort {
					for _, port := range service.Spec.Ports {
						portStrs = append(portStrs, fmt.Sprintf("%d/%s", port.NodePort, port.Protocol))
					}
				} else {
					for _, port := range service.Spec.Ports {
						portStrs = append(portStrs, fmt.Sprintf("%d/%s", port.Port, port.Protocol))
					}
				}
				selector := []string{}
				for k, v := range service.Spec.Selector.MatchLabels {
					selector = append(selector, fmt.Sprintf("%s=%s", k, v))
				}
				t.AppendRow(table.Row{service.MetaData.Name, service.Spec.Type, strings.Join(selector, ","), service.Spec.ClusterIP, strings.Join(portStrs, " ")})
			}
		}
	case core.ObjJob:
		t.AppendHeader(table.Row{"NAME", "POD", "STATUS"})

	case core.ObjHpa:
		t.AppendHeader(table.Row{"NAME", "REFERENCE", "TARGETS", "MINPODS", "MAXPODS"})
		hpalist := []core.HorizontalPodAutoscaler{}
		if resp != "" {
			if name == "" {
				err := utils.JsonUnMarshal(resp, &hpalist)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
			} else {
				hpa := core.HorizontalPodAutoscaler{}
				err := utils.JsonUnMarshal(resp, &hpa)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
				hpalist = append(hpalist, hpa)
			}
			for _, hpa := range hpalist {
				target := ""
				t.AppendRow(table.Row{hpa.MetaData.Name, fmt.Sprintf("%s/%s", hpa.Spec.ScaleTargetRef.Kind, hpa.Spec.ScaleTargetRef.Name), target, hpa.Spec.MinReplicas, hpa.Spec.MaxReplicas})
			}
		}
	case core.ObjFunction:
		t.AppendHeader(table.Row{"NAME", "PATH", "STATUS"})
		functions := []core.Function{}
		if resp != "" {
			if name == "" {
				err := utils.JsonUnMarshal(resp, &functions)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
			} else {
				function := core.Function{}
				err := utils.JsonUnMarshal(resp, &function)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
				functions = append(functions, function)
			}
			for _, function := range functions {
				t.AppendRow(table.Row{function.Name, function.Path, function.Status})
			}
		}
	case core.ObjWorkflow:
		t.AppendHeader(table.Row{"NAME", "STATUS"})
		workflows := []core.Workflow{}
		if resp != "" {
			if name == "" {
				err := utils.JsonUnMarshal(resp, &workflows)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
			} else {
				workflow := core.Workflow{}
				err := utils.JsonUnMarshal(resp, &workflow)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
				workflows = append(workflows, workflow)
			}
			for _, workflow := range workflows {
				t.AppendRow(table.Row{workflow.Name, workflow.Status})
			}
		}
	case core.ObjEndPoint:
		t.AppendHeader(table.Row{"NAME", "TARGETS"})
		endpoints := []core.Endpoint{}
		if resp != "" {
			if name == "" {
				err := utils.JsonUnMarshal(resp, &endpoints)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
			} else {
				endpoint := core.Endpoint{}
				err := utils.JsonUnMarshal(resp, &endpoint)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
				endpoints = append(endpoints, endpoint)
			}
			for _, endpoint := range endpoints {
				targets := ""
				for _, bind := range endpoint.Binds {
					dests := []string{}
					for _, dest := range bind.Destinations {
						dests = append(dests, fmt.Sprintf("%s:%d", dest.IP, dest.Port))
					}
					targets = fmt.Sprintf("%s:%d -> %s", endpoint.ServiceClusterIP, bind.ServicePort, strings.Join(dests, ","))
				}
				t.AppendRow(table.Row{endpoint.MetaData.Name, targets})
			}
		}
	case core.ObjDNS:
		t.AppendHeader(table.Row{"NAME", "HOST", "PATH(S)"})
		dnslist := []core.DNSRecord{}
		if resp != "" {
			if name == "" {
				err := utils.JsonUnMarshal(resp, &dnslist)
				if err != nil {
					fmt.Printf("error in unmarshal %s", err.Error())
				}
			} else {
				dns := core.DNSRecord{}
				err := utils.JsonUnMarshal(resp, &dns)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
				dnslist = append(dnslist, dns)
			}
			for _, dns := range dnslist {
				paths := []string{}
				for _, path := range dns.Paths {
					paths = append(paths, fmt.Sprintf("%s/%s:%d", path.Service, path.IP, path.Port))
				}
				t.AppendRow(table.Row{dns.MetaData.Name, dns.Host, strings.Join(paths, "\n")})
			}
		}
	case core.ObjVolume:
		t.AppendHeader(table.Row{"NAME", "CLASSNAME", "SERVER"})
		volumes := []core.PersistentVolume{}
		if resp != "" {
			if name == "" {
				err := utils.JsonUnMarshal(resp, &volumes)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
			} else {
				volume := core.PersistentVolume{}
				err := utils.JsonUnMarshal(resp, &volume)
				if err != nil {
					fmt.Printf("error in unmarshal %s\n", err.Error())
				}
				volumes = append(volumes, volume)
			}
			for _, volume := range volumes {
				t.AppendRow(table.Row{volume.MetaData.Name, volume.Spec.StorageClassName, volume.Spec.Nfs.Server})
			}
		}
	default:
	}
	t.Render()
}
