package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var RouteTable = [...]Route{
	{Path: "/api/v1/namespaces/:namespace/pods", Method: "POST", Handler: CreatePodHandler},                   // POST, create a pod
	{Path: "/api/v1/namespaces/:namespace/pods/:name", Method: "GET", Handler: GetPodHandler},                 // GET, get a pod
	{Path: "/api/v1/namespaces/:namespace/pods", Method: "GET", Handler: GetPodListHandler},                   // GET, list all pods
	{Path: "/api/v1/namespaces/:namespace/pods/:name", Method: "DELETE", Handler: DeletePodHandler},           // DELETE, delete a pod
	{Path: "/api/v1/namespaces/:namespace/pods/:name", Method: "PUT", Handler: UpdatePodHandler},              // PUT, update a pod
	{Path: "/api/v1/pods", Method: "GET", Handler: GetAllPodsHandler},                                         // GET, get alwatchPodl pods
	{Path: "/api/v1/namespaces/:namespace/pods/:name", Method: "POST", Handler: ReplacePodHandler},            // POST, replace a pod
	{Path: "/api/v1/namespaces/:namespace/pods/:name/status", Method: "PUT", Handler: UpdatePodStatusHandler}, // PUT, update pod status

	{Path: "/api/v1/nodes", Method: "POST", Handler: CreateNodeHandler},              // POST, create a node
	{Path: "/api/v1/nodes", Method: "GET", Handler: GetAllNodesHandler},              // GET, list all nodes
	{Path: "/api/v1/nodes/:name", Method: "GET", Handler: GetNodeHandler},            // GET, get a node
	{Path: "/api/v1/nodes/:name/pods", Method: "GET", Handler: GetPodsInNodeHandler}, // GET, get pods of a node
	{Path: "/api/v1/nodes/:name", Method: "DELETE", Handler: DeleteNodeHandler},      // DELETE, delete a node
	{Path: "/api/v1/nodes/:name", Method: "PUT", Handler: UpdateNodeHandler},         // PUT, update a node

	{Path: "/api/v1/namespaces/:namespace/services", Method: "POST", Handler: CreateServiceHandler},                            // POST, create a service
	{Path: "/api/v1/namespaces/:namespace/services/:name", Method: "GET", Handler: GetServiceHandler},                          // GET, get a service
	{Path: "/api/v1/namespaces/:namespace/services", Method: "GET", Handler: GetServiceListHandler},                            // GET, list all services
	{Path: "/api/v1/namespaces/:namespace/services/:name", Method: "DELETE", Handler: DeleteServiceHandler},                    // DELETE, delete a service
	{Path: "/api/v1/namespaces/:namespace/services/:name", Method: "PUT", Handler: UpdateServiceHandler},                       // PUT, update a service
	{Path: "/api/v1/namespaces/:namespace/services/:name/clusterip", Method: "GET", Handler: GetServiceClusterIPHandler},       // GET, get a service clusterIP
	{Path: "/api/v1/namespaces/:namespace/services/:name/clusterip", Method: "DELETE", Handler: DeleteServiceClusterIPHandler}, // DELETE, free a service clusterIP

	{Path: "/api/v1/namespaces/:namespace/endpoints", Method: "POST", Handler: CreateEndpointHandler},         // POST, create an endpoint
	{Path: "/api/v1/namespaces/:namespace/endpoints/:name", Method: "GET", Handler: GetEndpointHandler},       // GET, get an endpoint
	{Path: "/api/v1/namespaces/:namespace/endpoints", Method: "GET", Handler: GetEndpointListHandler},         // GET, list all endpoints in this namespace
	{Path: "/api/v1/namespaces/:namespace/endpoints/:name", Method: "DELETE", Handler: DeleteEndpointHandler}, // DELETE, delete a endpoint
	{Path: "/api/v1/namespaces/:namespace/endpoints/:name", Method: "PUT", Handler: UpdateEndpointHandler},    // PUT, update a endpoint

	{Path: "/api/v1/namespaces/:namespace/dns", Method: "POST", Handler: CreateDNSHandler},         // POST, create a dns
	{Path: "/api/v1/namespaces/:namespace/dns/:name", Method: "GET", Handler: GetDNSHandler},       // GET, get a dns
	{Path: "/api/v1/namespaces/:namespace/dns", Method: "GET", Handler: GetDNSListHandler},         // GET, list all dns
	{Path: "/api/v1/dns", Method: "GET", Handler: GetAllDNSHandler},                                //Get, get all dns
	{Path: "/api/v1/namespaces/:namespace/dns/:name", Method: "DELETE", Handler: DeleteDNSHandler}, // DELETE, delete a dns
	{Path: "/api/v1/namespaces/:namespace/dns/:name", Method: "PUT", Handler: UpdateDNSHandler},    // PUT, update a dns

	{Path: "/api/v1/namespaces/:namespace/replicas", Method: "POST", Handler: CreateReplicaHandler},                   // POST, create a replica
	{Path: "/api/v1/namespaces/:namespace/replicas/:name", Method: "GET", Handler: GetReplicaHandler},                 // GET, get a replica
	{Path: "/api/v1/namespaces/:namespace/replicas", Method: "GET", Handler: GetReplicaListHandler},                   // GET, list all replicas
	{Path: "/api/v1/namespaces/:namespace/replicas/:name", Method: "DELETE", Handler: DeleteReplicaHandler},           // DELETE, delete a replica
	{Path: "/api/v1/namespaces/:namespace/replicas/:name", Method: "PUT", Handler: UpdateReplicaHandler},              // PUT, update a replica
	{Path: "/api/v1/namespaces/:namespace/replicas/:name/status", Method: "PUT", Handler: UpdateReplicaStatusHandler}, // PUT, update replica status | owner-reference

	{Path: "/api/v1/namespaces/:namespace/jobs", Method: "POST", Handler: CreateJobHandler},         // POST, create a Job
	{Path: "/api/v1/namespaces/:namespace/jobs/:name", Method: "GET", Handler: GetJobHandler},       // GET, get a Job
	{Path: "/api/v1/namespaces/:namespace/jobs", Method: "GET", Handler: GetJobListHandler},         // GET, list all Jobs
	{Path: "/api/v1/namespaces/:namespace/jobs/:name", Method: "DELETE", Handler: DeleteJobHandler}, // DELETE, delete a Job
	{Path: "/api/v1/namespaces/:namespace/jobs/:name", Method: "PUT", Handler: UpdateJobHandler},    // PUT, update a Job

	{Path: "/api/v1/functions", Method: "POST", Handler: CreateFunctionHandler},                // POST, create a function
	{Path: "/api/v1/functions/:name", Method: "GET", Handler: GetFunctionHandler},              // GET, get a function
	{Path: "/api/v1/functions/:name", Method: "DELETE", Handler: DeleteFunctionHandler},        // DELETE, delete a function
	{Path: "/api/v1/functions/:name", Method: "PUT", Handler: UpdateFunctionHandler},           // PUT, update a function
	{Path: "/api/v1/functions/:name/trigger", Method: "POST", Handler: TriggerFunctionHandler}, // POST, trigger a function
	{Path: "/api/v1/functions", Method: "GET", Handler: GetAllFunctionsHandler},                // GET, list all functions
	{Path: "/api/v1/functions/tasks", Method: "POST", Handler: CreateTaskHandler},
	{Path: "/api/v1/functions/tasks/:name", Method: "GET", Handler: GetTaskHandler},
	{Path: "/api/v1/functions/tasks", Method: "GET", Handler: GetAllTaskHandler},
	{Path: "/api/v1/functions/tasks/:name", Method: "POST", Handler: UpdateTaskHandler},
	{Path: "/api/v1/functions/tasks/:name", Method: "DELETE", Handler: DeleteTaskHandler},

	{Path: "/api/v1/namespaces/:namespace/hpa", Method: "POST", Handler: CreateHpaHandler},                   // POST, create a hpa
	{Path: "/api/v1/namespaces/:namespace/hpa/:name", Method: "GET", Handler: GetHpaHandler},                 // GET, get a hpa
	{Path: "/api/v1/namespaces/:namespace/hpa", Method: "GET", Handler: GetHpaListHandler},                   // GET, list all hpa
	{Path: "/api/v1/namespaces/:namespace/hpa/:name", Method: "DELETE", Handler: DeleteHpaHandler},           // DELETE, delete a hpa
	{Path: "/api/v1/namespaces/:namespace/hpa/:name", Method: "PUT", Handler: UpdateHpaHandler},              // PUT, update a hpa
	{Path: "/api/v1/hpa", Method: "GET", Handler: GetAllHpaHandler},                                          // GET, get all hpa
	{Path: "/api/v1/namespaces/:namespace/hpa/:name/status", Method: "PUT", Handler: UpdateHpaStatusHandler}, // PUT, update hpa status

	{Path: "/api/v1/workflows", Method: "POST", Handler: CreateWorkflowHandler},                // POST, create a workflow
	{Path: "/api/v1/workflows/:name", Method: "GET", Handler: GetWorkflowHandler},              // GET, get a workflow
	{Path: "/api/v1/workflows", Method: "GET", Handler: GetWorkflowListHandler},                // GET, list all workflows
	{Path: "/api/v1/workflows/:name", Method: "DELETE", Handler: DeleteWorkflowHandler},        // DELETE, delete a workflow
	{Path: "/api/v1/workflows/:name", Method: "PUT", Handler: UpdateWorkflowHandler},           // PUT, update a workflow
	{Path: "/api/v1/workflows/:name/trigger", Method: "POST", Handler: TriggerWorkflowHandler}, // POST, trigger a workflow

	{Path: "/api/v1/volumes", Method: "POST", Handler: CreateVolumeHandler},         // POST, create a volume
	{Path: "/api/v1/volumes/:name", Method: "GET", Handler: GetVolumeHandler},       // GET, get a volume
	{Path: "/api/v1/volumes", Method: "GET", Handler: GetAllVolumeHandler},          // GET, get all volume
	{Path: "/api/v1/volumes/:name", Method: "DELETE", Handler: DeleteVolumeHandler}, // DELETE, delete a volume
	{Path: "/api/v1/csivolumes/:name", Method: "GET", Handler: GetCsiVolumeHandler}, // GET, get a csivolume
	{Path: "/api/v1/watch/:name", Method: "GET", Handler: GetDNSHandler},            // get a DNS watcj
	{Path: "/metrics", Method: "GET", Handler: gin.WrapH(promhttp.Handler())},       // GET, get metrics
}
