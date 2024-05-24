#!/bin/bash

source ./script/env.sh

if [ ! -d $BUILDDIR ]; then
  mkdir ./build
fi

cd build
cmake ..
make -j`nproc`
cd ..