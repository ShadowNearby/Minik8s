nerdctl stop `nerdctl ps -aq` && nerdctl rm `nerdctl ps -aq`
service systemd-resolved stop && service coredns start

# node test
./bin/kubectl get node

# pod test
./bin/kubectl apply -f ./example/pod/pod_multi_containers.yaml
./bin/kubectl get pod
./bin/kubectl delete pod pod-net1

# pod and pv test
./bin/kubectl apply -f ./example/volume/volume.yaml
./bin/kubectl get volume
./bin/kubectl delete volume test-pv
./bin/kubectl apply -f ./example/pod/pod_with_volume.yaml
./bin/kubectl get pod
nerdctl exec -it pod-net1-test-container /bin/bash
./bin/kubectl delete pod pod-volume

# service test
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

# replicaset test
./bin/kubectl apply -f ./example/replicaset/replicaset.yaml
./bin/kubectl get pod
./bin/kubectl get replicas
./bin/kubectl apply -f ./example/replicaset/rs_service.yaml
./bin/kubectl apply -f ./example/replicaset/replicaset.yaml -u ip-return-deployment
./bin/kubectl get service
./bin/kubectl delete replicas ip-return-deployment
./bin/kubectl delete service rs-service

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
./bin/kubectl apply -f ./example/hpa/replicaset.yaml
./bin/kubectl apply -f ./example/hpa/hpa.yaml
./bin/kubectl apply -f ./example/hpa/service.yaml
./bin/kubectl get pod
./bin/kubectl get hpa
./bin/kubectl get replicas
./bin/kubectl get service
./bin/kubectl get endpoint
nerdctl exec -it --- /bin/bash
stress --vm 2 --vm-bytes 500M --vm-keep
./bin/kubectl delete hpa test-hpa
./bin/kubectl delete service hpa-service