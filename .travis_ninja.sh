#!/bin/sh
# Installation of ninja

initdir=$(pwd)

mkdir ../../ninja-build
cd ../../ninja-build || echo 'Cannot open dir ../../ninja-build'

git clone https://github.com/ninja-build/ninja
cd ninja || echo 'Cannot cd ninja'

./configure.py --bootstrap

mkdir ./bin && mv ./ninja ./bin

export PATH=$(pwd)/bin:$PATH
cd "$initdir" || echo 'Cannot open dir ' "$initdir"

unset initdir
