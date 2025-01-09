#!/bin/bash

echo "Running Redis on Docker"

docker run -d --name redis-stack -p 6379:6379 -p 8001:8001 -v $(pwd)/database/redis/:/data redis/redis-stack:latest
