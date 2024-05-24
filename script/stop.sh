#!/bin/bash

source ./script/env.sh

for program in "${PROGRAMS[@]}"; do
    getpid $program
    if [ "$pid" != "" ]; then
        pkill $program
    fi
done
