#!/bin/bash

source ./script/env.sh

for program in "${PROGRAMS[@]}"; do
    getpid $program
    if [ "$pid" != "" ]; then
        echo "${program} is running (pid: ${pid})! kill it first"
        exit 1
    fi
done

cp $CONFDIR/nginx.conf.bk $CONFDIR/nginx.conf
cp $CONFDIR/prometheus/sd_node.json.bk $CONFDIR/prometheus/sd_node.json
cp $CONFDIR/prometheus/sd_pod.json.bk $CONFDIR/prometheus/sd_pod.json

nginx -c $CONFDIR/nginx.conf
nohup ./bin/apiserver > $LOGDIR/apiserver.log 2>&1 & 
sleep 3
nohup ./bin/kubelet > $LOGDIR/kubelet.log 2>&1 &
nohup ./bin/kubeproxy > $LOGDIR/kubeproxy.log 2>&1 &
nohup prometheus --config.file $CONFDIR/prometheus/prometheus.yml > $LOGDIR/prometheus.log 2>&1 & 
nohup ./bin/kmonitor > $LOGDIR/kmonitor.log 2>&1 &
