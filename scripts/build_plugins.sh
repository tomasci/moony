#!/bin/bash

echo "Building plugins"

PLUGIN_DIR="./plugins"
BUILD_DIR="./build/plugins"

mkdir -p $BUILD_DIR

for plugin_path in $PLUGIN_DIR/*; do
    if [ -d "$plugin_path" ]; then
        plugin_name=$(basename "$plugin_path")
        output_dir="$BUILD_DIR/$plugin_name"

        mkdir -p $output_dir

        go build -buildmode=plugin -o $output_dir/${plugin_name}.so $plugin_path/${plugin_name}.go

        cp $plugin_path/plugin.json $output_dir

        echo "  \xE2\x9C\x94 $plugin_name"
    fi
done