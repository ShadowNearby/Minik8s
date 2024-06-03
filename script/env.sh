#!/bin/bash

ENV=$1
LOGDIR=$(pwd)/log
PIDDIR=/var/run/minik8s
BUILDDIR=$(pwd)/build
CONFDIR=$(pwd)/config

if [ "$ENV" = "CI" ]; then
  LOGDIR=/var/log/minik8s
fi

PROGRAMS=("apiserver" "kubelet" "kubeproxy" "prometheus" "kmonitor" "nginx")

pid=""
getpid () {
  pid=$(ps -ef | grep $1  | grep root  | grep -v grep | awk '{print $2}')
}
