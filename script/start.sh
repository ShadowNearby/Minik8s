#!/bin/bash

source ./script/env.sh

nohup ./bin/apiserver > $LOGDIR/apiserver.log 2>&1 & 
sleep 3
nohup ./bin/kubelet > $LOGDIR/kubelet.log 2>&1 &
nohup prometheus --config.file ./config/prometheus/prometheus.yml > $LOGDIR/prometheus.log 2>&1 & 
nohup ./bin/monitor > $LOGDIR/monitor.log 2>&1 &
