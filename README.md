# Hyperledger Fabric Network

## Order for Running the Command Scripts

1. Run `create-artifacts.sh`
2. Run `docker-compose up -d`
3. Run `create-channel.sh`
4. Run `deploy-chaincode.sh`
5. Run `go mod tidy` on `/artifacts/src/github.com/fabcar/go`

---

## Guide to Creating a Hyperledger Fabric Network

### Configuration Files

- `/artifacts/channel/crypto-config.yaml`
  - Config file for creating identities for orderers and peers
  - Uses cryptogen tools
  - Defines multiple orderers and peer orgs
  - Template count: # of nodes
  - User count: 2 users (admin and two others)

### Shell Scripts

- `/artifacts/channel/create-artifacts.sh`
  - Generates crypto material based on `crypto-config.yaml` and `configtx.yaml`
  - Creates folders for orderer and peer orgs
  - Generates `genesis.block`, `mychannel.tx`, `Org1MSPanchors.tx`, `Org2MSPanchors.tx`
- `/artifacts/channel/configtx.yaml`
  - Creates channel artifacts, including `channel.tx` and `genesis.block`
  - Contains configurations for organizations, capabilities, application, orderer, channel, profile
- `/artifacts/docker-compose.yaml`
  - Defines Docker containers for the network components
  - Specifies environment variables and services for orderers, peer orgs, and databases

### Creating a Channel and Joining Peers

1. Ensure the network is running
2. Have the `channel-artifacts` folder
3. Run `/create-channel.sh`
   - Export required environment variables
   - Call necessary functions
   - Uses `core.yaml` from `artifacts/channel`
   - Executes `createChannel`, `joinChannel`, and `updateAnchorPeers`

---

## Dealing with Chaincode

1. Pre-setup (install dependencies)
2. Create a smart contract package
3. Install package on endorsing peers (`peer0` of `org1` and `org2`)
4. Query the installed chaincode on each peer
5. Approve the chaincode by each organization
6. Check commit readiness
7. Commit the chaincode to the peer
8. Query committed chaincode
9. Invoke initialization function
10. Invoke transactions
11. Query transactions
12. These steps are done in `/deployChaincode.sh`

---

## Creating a Network with Multiple Orderers

1. Modify `crypto-config.yaml`
   - Add more hostnames (`orderer2`, `orderer3`)
2. In `configtx.yaml`, add more orderers under `Orderer: EtcdRaft: Consenters`
3. Run `create-artifacts.sh` after changing config files
4. In `docker-compose.yaml`, define multiple orderers
   - Set `ORDERER_GENERAL_LISTEN_PORT` for each orderer
   - Map volumes for MSP and TLS certs
5. Run the usual commands (`docker-compose up`, `createChannel`, `deployChaincode`)

---

## Further Learning

- Starting from **Video 20**, explore Fabric SDK.

---

When the `/user` API is hit for user registration, the `./crypto/org1` and `fabric-client-kv-org1` folders are created. These folders contain public and private keys of the registered user.

