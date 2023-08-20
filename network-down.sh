#!/bin/bash

# Step 1: Delete files
rm -f ./fabcar.tar.gz ./token.tar.gz ./channel-artifacts/mychannel.block
rm -r -f ./artifacts/src/github.com/fabcar/go/vendor ./artifacts/src/github.com/token-chaincode/go/vendor

# Step 2: Move to directory
cd ./artifacts

# Step 3: Stop containers and wait
docker-compose down