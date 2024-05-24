#!/bin/bash

source ./script/env.sh

./script/stop.sh

./script/build.sh

./script/start.sh
