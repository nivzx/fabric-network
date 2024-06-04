#!/bin/bash

# Step 1: Delete files
rm -f ./level.tar.gz ./token.tar.gz ./channel-artifacts/mychannel.block
rm -r -f ./artifacts/src/github.com/fabcar/go/vendor ./artifacts/src/github.com/token-chaincode/go/vendor

# Step 2: Move to directory
cd ./artifacts

# Step 3: Stop containers and wait
docker-compose down
sleep 15

# Step 4: Start containers and wait
docker-compose up -d
sleep 15

# Step 5: Move back
cd ..

# Step 6: Run createChannel.sh and wait
./createChannel.sh
sleep 10

# Step 7: Vendoring Go dependencies for fabcar
echo "Vendoring Go dependencies for level ..."
pushd ./artifacts/src/github.com/fabcar/go
GO111MODULE=on go mod vendor
popd
echo "Finished vendoring Go dependencies for level"

# Step 8: Run deployChaincode.sh and wait
./deployChaincode.sh
sleep 25

# Step 9: Vendoring Go dependencies for token-chaincode
echo "Vendoring Go dependencies for token-chaincode ..."
pushd ./artifacts/src/github.com/token-chaincode/go
GO111MODULE=on go mod vendor
popd
echo "Finished vendoring Go dependencies for token-chaincode"

# Step 10: Run deployTokenChaincode.sh and wait
./deployTokenChaincode.sh
sleep 25

echo "Script execution completed."

# Step 11: Make API call
curl -X POST http://localhost:4000/users \
-H "Content-Type: application/json" \
-d '{
  "username": "boss",
  "orgName": "Org1"
}'
