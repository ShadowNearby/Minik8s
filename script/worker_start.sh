#!/bin/bash

source ./script/env.sh

nohup ./bin/kubelet > $LOGDIR/kubelet.log 2>&1 &

nohup ./bin/kubeproxy > $LOGDIR/kubeproxy.log 2>&1 &
