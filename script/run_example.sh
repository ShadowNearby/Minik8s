./bin/kubectl apply -f ./example/volume/volume.yaml

./bin/kubectl get volume

./bin/kubectl delete volume test-pv

./bin/kubectl apply -f ./example/pod/pod_multi_containers.yaml

./bin/kubectl get pod

./bin/kubectl delete pod pod-net1