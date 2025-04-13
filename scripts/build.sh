#!/bin/bash

# cleanup
# create temporary folder
mkdir ./build-tmp
# move storage data to temp folder
mv ./build/storage ./build-tmp
# remove build folder
rm -rf ./build
# create build folder
mkdir ./build
# move storage back
mv ./build-tmp/storage ./build
# remove temp dir
rm -rf ./build-tmp

# build server
sleep 0.2
./scripts/build_server.sh

# build plugins
sleep 0.2
./scripts/build_plugins.sh