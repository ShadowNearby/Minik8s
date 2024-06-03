nerdctl stop `nerdctl ps -aq` && nerdctl rm `nerdctl ps -aq`

./bin/kubectl get node

./bin/kubectl apply -f ./example/pod/pod_multi_containers.yaml
./bin/kubectl get pod
./bin/kubectl delete pod pod-net1

./bin/kubectl apply -f ./example/volume/volume.yaml
./bin/kubectl get volume
./bin/kubectl delete volume test-pv
./bin/kubectl apply -f ./example/pod/pod_with_volume.yaml
./bin/kubectl get pod
nerdctl exec -it pod-net1-test-container /bin/bash
./bin/kubectl delete pod pod-volume

./bin/kubectl apply -f ./example/service/service_pod.yaml
./bin/kubectl apply -f ./example/pod/worker_pod.yaml
./bin/kubectl get pod 
nerdctl exec -it pod-on-worker-test-container /bin/bash
./bin/kubectl apply -f ./example/service/service_clusterip.yaml
./bin/kubectl apply -f ./example/service/service_nodeport.yaml
./bin/kubectl get service 
./bin/kubectl get endpoint
./bin/kubectl delete service service-clusterip
./bin/kubectl delete service service-nodeport
./bin/kubectl delete pod service-pod
./bin/kubectl delete pod pod-on-worker

./bin/kubectl apply -f ./example/replicaset/replicaset.yaml
./bin/kubectl get pod
./bin/kubectl get replicas
./bin/kubectl apply -f ./example/replicaset/rs_service.yaml
./bin/kubectl apply -f ./example/replicaset/replicaset.yaml -u ip-return-deployment
./bin/kubectl get service
./bin/kubectl delete replicas ip-return-deployment
./bin/kubectl delete service rs-service