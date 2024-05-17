pkill kubelet
pkill apiserver
pkill nginx
cd build && make -j && cd ..
./bin/apiserver > ./log/apiserver.log 2>&1 &
sleep 3
./bin/kubelet > ./log/kubelet.log 2>&1 &
