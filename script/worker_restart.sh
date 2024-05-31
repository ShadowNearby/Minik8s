#!/bin/bash

source ./script/env.sh

./script/worker_stop.sh

./script/build.sh

./script/worker_start.sh