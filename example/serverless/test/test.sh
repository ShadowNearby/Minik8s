#!/bin/bash

for i in $(seq 1 20);
do
    curl -X POST http://192.168.1.12:8090/api/v1/functions/get-sum/trigger \
        -H "Content-Type: application/json" \
        -d '{"kind": "trigger", "name": "get-sum", "params": "{\"x\": 3, \"y\": 4}"}'
    # sleep 2
done
