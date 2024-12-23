#!/bin/bash

# cleanup
rm -rf ./build

# create dir for build
mkdir ./build

# build server
sleep 0.2
./scripts/build_server.sh

# build plugins
sleep 0.2
./scripts/build_plugins.sh