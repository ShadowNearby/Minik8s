#!/bin/bash

source ./script/env.sh

nohup ./bin/kubelet > $LOGDIR/kubelet.log 2>&1 &
