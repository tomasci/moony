#!/bin/bash

# Source the environment variables from the .env file
if [ -f .env ]; then
    set -a
    source .env
    set +a
fi