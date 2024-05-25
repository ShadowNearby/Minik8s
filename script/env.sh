#!/bin/bash

ENV=$1
LOGDIR=$(pwd)/log
PIDDIR=/var/run/minik8s
BUILDDIR=$(pwd)/build
CONFDIR=$(pwd)/config

if [ "$ENV" = "CI" ]; then
  LOGDIR=/var/log/minik8s
fi

PROGRAMS=("apiserver" "kubelet" "prometheus" "monitor" "nginx")

pid=""
getpid () {
  pid=$(ps -ef | grep $1 | grep -v grep | awk '{print $2}')
}
