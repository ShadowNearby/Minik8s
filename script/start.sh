./bin/apiserver > ./log/apiserver.log 2>&1 &
sleep 3
./bin/kubelet > ./log/kubelet.log 2>&1 &