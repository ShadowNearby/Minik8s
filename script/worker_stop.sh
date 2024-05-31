#!/bin/bash

source ./script/env.sh

getpid kubelet
if [ "$pid" != "" ]; then
    pkill kubelet
fi

getpid kubeproxy
if [ "$pid" != "" ]; then
    pkill kubeproxy
fi