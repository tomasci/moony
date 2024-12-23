#!/bin/bash

# cleanup
rm -rf ./build

# create dir for build
sleep 0.2
mkdir ./build

# build server
sleep 0.2
./scripts/build_server.sh

# build plugins
sleep 0.2
./scripts/build_plugins.sh