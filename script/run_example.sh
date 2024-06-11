nerdctl stop `nerdctl ps -aq` && nerdctl rm `nerdctl ps -aq`
service systemd-resolved stop && service coredns start

# node test
./bin/kubectl get node

# pod test
./bin/kubectl apply -f ./example/pod/pod_multi_containers.yaml
nerdctl exec -it pod-net1-test-container curl localhost:8080
nerdctl exec -it pod-net1-test-container curl localhost:8090
./bin/kubectl get pod
./bin/kubectl delete pod pod-net1
./bin/kubectl apply -f ./example/pod/worker_pod.yaml
./bin/kubectl apply -f ./example/pod/master_pod.yaml
./bin/kubectl delete pod pod-on-worker
./bin/kubectl delete pod pod-on-master

# pod kill test
./bin/kubectl apply -f ./example/pod/system_load_pod.yaml
./bin/kubectl get pod
curl :7070/memory
./bin/kubectl delete pod pod-load

# pod and pv test
./bin/kubectl apply -f ./example/volume/volume.yaml
./bin/kubectl get volume
./bin/kubectl delete volume test-pv
./bin/kubectl apply -f ./example/pod/pod_with_volume.yaml
./bin/kubectl apply -f ./example/pod/pod_with_volume_worker.yaml
./bin/kubectl get pod
nerdctl exec -it pod-volume-test1 touch /test_mount/pod-pv/hello.txt
nerdctl exec -it pod-volume-test2 ls /test_mount/pod-pv/hello.txt
nerdctl exec -it pod-volume-worker-test1 ls /test_mount/pod-pv/hello.txt
nerdctl exec -it pod-volume-test2 rm /test_mount/pod-pv/hello.txt
./bin/kubectl delete pod pod-volume
./bin/kubectl delete pod pod-volume-worker

# service test
./bin/kubectl apply -f ./example/service/service_pod.yaml
./bin/kubectl apply -f ./example/pod/worker_pod.yaml
./bin/kubectl get pod 
./bin/kubectl apply -f ./example/service/service_clusterip.yaml
./bin/kubectl apply -f ./example/service/service_nodeport.yaml
./bin/kubectl get service
nerdctl exec -it pod-on-worker-test-container /bin/bash
./bin/kubectl get endpoint
./bin/kubectl delete service service-clusterip
./bin/kubectl delete service service-nodeport
./bin/kubectl delete pod service-pod
./bin/kubectl delete pod pod-on-worker

# replicaset test
./bin/kubectl apply -f ./example/replicaset/replicaset.yaml
./bin/kubectl get pod
./bin/kubectl get replicas
./bin/kubectl apply -f ./example/replicaset/rs_service.yaml
./bin/kubectl apply -f ./example/replicaset/replicaset.yaml -u ip-return-deployment
./bin/kubectl get service
./bin/kubectl get endpoint
./bin/kubectl apply -f ./example/pod/master_pod.yaml
nerdctl exec -it pod-on-master-test-container curl :20922
./bin/kubectl delete pod pod-on-master
./bin/kubectl delete service rs-service
./bin/kubectl delete replicas ip-return-deployment

# replicaset fail test
./bin/kubectl apply -f ./example/replicaset/rs_fail.yaml
./bin/kubectl get pod
./bin/kubectl apply -f ./example/replicaset/rs_fail_service.yaml
./bin/kubectl get service
./bin/kubectl get endpoint
curl :7070/memory
./bin/kubectl delete service rs-fail-service
./bin/kubectl delete replicas fail-deployment


# dns test
./bin/kubectl apply -f ./example/dns/pod1.yaml
./bin/kubectl apply -f ./example/dns/pod2.yaml
./bin/kubectl apply -f ./example/dns/service1.yaml
./bin/kubectl apply -f ./example/dns/service2.yaml
./bin/kubectl apply -f ./example/dns/dns.yaml
./bin/kubectl get pod
./bin/kubectl get service
./bin/kubectl get endpoint
./bin/kubectl get dns
curl dnstest.com/service1
curl dnstest.com/service2
nerdctl exec -it service-pod1-test-container /bin/bash
./bin/kubectl delete pod service-pod1
./bin/kubectl delete pod service-pod2
./bin/kubectl delete service service1
./bin/kubectl delete service service2
./bin/kubectl delete dns dns-test

# hpa test
cat ./config/config.json | grep HPA
./bin/kubectl apply -f ./example/hpa/replicaset.yaml
./bin/kubectl apply -f ./example/hpa/hpa.yaml
./bin/kubectl apply -f ./example/hpa/service.yaml
./bin/kubectl get pod
./bin/kubectl get hpa
./bin/kubectl get replicas
./bin/kubectl get service
./bin/kubectl get endpoint
./bin/kubectl apply -f ./example/pod/master_pod.yaml
nerdctl exec -it --- /bin/bash
stress --vm 2 --vm-bytes 1600M --vm-keep
./bin/kubectl delete hpa test-hpa
./bin/kubectl delete service hpa-service

# test registry
./bin/kubectl apply -f ./example/pod/registry_pod.yaml
./bin/kubectl get pod
./bin/kubectl delete pod pod-on-registry

# test monitor
./bin/kubectl apply -f ./example/monitor/monitor_pod.yaml
./bin/kubectl get pod
./bin/kubectl delete pod monitor-pod

# test serverless basic
./bin/kubectl apply -f ./example/volume/volume.yaml
./bin/kubectl apply -f ./example/serverless/common/getsum.yaml
./bin/kubectl trigger functions -f ./example/serverless/common/trigger_functions.yaml
./bin/kubectl result functions 5fdf
./bin/kubectl apply -f ./example/serverless/common/getsum.yaml -u getsum

# test workflow
./bin/kubectl apply -f ./example/serverless/application/infer.yaml
./bin/kubectl apply -f ./example/serverless/application/preprocess.yaml
./bin/kubectl apply -f ./example/serverless/application/application_workflow.yaml
./bin/kubectl apply -f ./example/volume/volume.yaml
RESULT_ID = ./bin/kubectl trigger workflows -f ./example/serverless/applications/trigger_workflow.yaml
echo $(RESULT_ID)
./bin/kubectl result workflows $(RESULT_ID)

# test function
./bin/kubectl apply -f ./example/serverless/application/preprocess.yaml
RESULT_ID = ./bin/kubectl trigger functions -f ./example/serverless/applications/trigger_workflow.yaml
echo $(RESULT_ID)
./bin/kubectl result functions $(RESULT_ID)

# test fault torrence
./bin/kubectl apply -f ./example/service/service_pod.yaml
./bin/kubectl get pod 
./bin/kubectl apply -f ./example/service/service_clusterip.yaml
./bin/kubectl get service
./bin/kubectl get endpoint
./script/stop.sh
./script/start.sh
./bin/kubectl delete service service-clusterip
./bin/kubectl delete pod service-pod
