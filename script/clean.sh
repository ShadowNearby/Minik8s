#!/bin/bash

ps=$(nerdctl ps -aq)
if [ ! -z "$ps" ]; then
  nerdctl stop $ps
  nerdctl rm $ps
fi

redis-cli FLUSHALL

ETCDOBJS=("nodes" "pods" "services" "hpa" "replicas" "jobs" "functions" "workflows" "deployment" "endpoints" "dns" "volumes" "csivolumes")

for obj in "${ETCDOBJS[@]}"; do
  count=$(etcdctl del --prefix /$obj)
  echo "etcd: delete $obj count: $count"
done
