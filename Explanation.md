
# Building a Hyperledger Fabric Network

Here, every command run to bring up a simple Hyperledger Fabric network is shown.
You must have following Fabric platform-specific binaries installed

 - configtxgen
 - configtxlator
 - cryptogen
 - discover
 - idemixgen
 - orderer
 - peer
 - fabric-ca-client
 - fabric-ca-server
**********************************************************************************
## Generating Crypto Materials & Channel Configuration

### Generate Crypto Artifacts for Organizations:

`cryptogen generate --config=./crypto-config.yaml --output=./crypto-config/`

 - This command receives configuration details from the crypto-config.yaml file and generates the crypto material for organizations at `./crypto-config/` directory
- More info on cryptogen command and sample crypto-config.yaml file can be seen [here](https://hyperledger-fabric.readthedocs.io/en/release-1.1/commands/cryptogen-commands.html) 


### Generate System Genesis Block
`configtxgen -profile OrdererGenesis -configPath path/to/configtx.yaml channelID "sys-channel" -outputBlock ./genesis.block`
- This configtxgen command receives path to the configtx.yaml file as -configPath parameter
- Also it receives the profile from configtx.yaml file to use for generation

### Generate the Channel Configuration Block
`configtxgen -profile BasicChannel -configPath path/to/configtx.yaml -outputCreateChannelTx ./mychannel.tx -channelID mychannel`
- This command creates also receives configtx.yaml configuration file and creates channel named "mychannel" and mychannel.tx file

### Generate Anchor Peer Updates for Peer Organizations
`configtxgen -profile BasicChannel -configPath path/to/configtx.yaml -outputAnchorPeersUpdate ./Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP`
- Similarly, this command generates anchor peer tx file for Org1.
- Here the `-asOrg` parameter specifies the organization ID. It performs the config generation as this specified organization, only including values in the write set that org has privilege to set
- Repeat this command for all the organizations in your specific network


`configtx.yaml` file is used to build the initial channel configuration that is stored in the genesis block. It includes following sections
- Organizations
- Capabilities
- Application
- Orderer
- Channel
- Profiles

More info about configtxgen command : [Hyperlegder docs on configtxgen](https://hyperledger-fabric.readthedocs.io/en/latest/commands/configtxgen.html)
More info about configtx.yaml : [Hyperledger docs on configtx.yaml](https://hyperledger-fabric.readthedocs.io/en/latest/create_channel/create_channel_config.html)

By now, all the crypto material and channel artifacts required to bring up the network are generated successfuly. You can provide these generated files in the docker-compose as env variable to bring up all required containers

Next is to use the docker-compose file to download all the containters from docker hub.

## Docker Compose File

`networks`:
- defines the network used to link between containers (ex: test)

`services`:
- `ca-org1`, `ca-org2`, etc. - `image: hyperledger/fabric-ca`
- `orderer1.example.com`, 2, 3 - `image: hyperledger/fabric-orderer`
- `couchdb0`, 1, 2, (for each peer) - `image: hyperledger/fabric-couchdb`
- `peer0.org1.example.com`, ... - `image: hyperledger/fabric-peer`

In addition, all these services need a bunch of environment variables, port definition, volume definition, dns_search (orderer)

Sample docker-compose file can bee seen the the git repository.

Now the network is up and running with all the docker containers. What's left is to create and join a channel and to deploy chaincode on it.

## Channels

### Channel Creation
We use `peer` binaries for channel creation.

`peer channel create -o localhost:7050 -c mychannel --ordererTLSHostNameOverride order.example.com -f path/to/mychannel.tx --outputBlock location/to/save --tls true --cafile path/to/tls`

Okay a lot is going on in the above command.
- `-o` flag specifies the address of the ordering service to which the channel creation request will be sent
- `-c` channel name
- `-ordererTLSHostNameOverride` flag is used to specify the hostname to be used when communicating with the orderer over TLS
- `-f` location of channel configuration .tx file. We created this in previous step
- `-outputBlock` location where the channel creation block .block file will be saved
- `--tls` indicates whether TLS should be enabled for communication
- `--cafile` CA certificate file that should be used to verify the TLS connection to the orderer. We generated this in the first step. (`crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem`)


### Joining the Channel

`peer channel join -b ./path/to/mychannel.block`
- Here the `-b` flag receives the path to the block file. We created this in the previous step.
- The block file contains information about the channel, such as its name, configuration, and initial state.
- Peer uses this information to configure itself to communicate and interact with the specific channel.
- It establishes connections with other peers already in the channel, as well as with the ordering service responsible for validating and ordering transactions on the channel.

*You should run this command for all the peers who needs to join that channel, by changing the global variables for each at a time*


### Updating the anchor peers
`peer channel update` command allows an organization to change channel configuration. This is used to add or update anchor peers.

`peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c mychannel -f ./path/to/anchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA`

- In this command, most of them is same as the channel creation. Only new thing is `-f` flag.
- `-f` flag points to the anchor peer update transaction .tx file

*This command also should be run for anchor peers in every organization. (By changing globals)*

Now, the channel is created, peers have joined and anchor peers have been configured!
All thats left is to deploy the chaincode on the channel.


## Chaincode

To install the dependencies, at the location where chaincode.go and go.mod files are,
- `go mod tidy`
- `go mod vendor`

Then we need to package, install, commit, ... chaincode to peers. This is described in the Fabric chaincode lifecycle.
We use peer lifecycle chaincode command for this. This command has subcommands and flags we will be using here.

More info about peer lifecycle chaincode : [Fabric docs on peer lifecyle](https://hyperledger-fabric.readthedocs.io/en/release-2.5/commands/peerlifecycle.html)

### 1. Package

`setGlobalsForPeer0Org1`

    peer lifecycle chaincode package fabcar.tar.gz \
            --path ${CC_SRC_PATH} --lang golang \
            --label ${CC_NAME}_${VERSION}

First, the global variables are set to the required peer.

- `fabcar.tar.gz` file is the resulting tarball file.
- `--path` flac specifies the path to fabcar.go file
- `--lang` programming language used for the chaincode (golang)
- `--label` flac sets a label for the package

### 2. Install 

    setGlobalsForPeer0Org1
    peer lifecycle chaincode install ${CC_NAME}.tar.gz
    
    setGlobalsForPeer0Org2
    // do the same

This installs the previously package chaincode on peer0Org1. You should install this chaincode for each and every endorsing peers defined. You do this by running peer lifecycle chaincode install command by changing the globals

### 3. Query Installed

`setGlobalsForPeer0Org1`

    peer lifecycle chaincode queryinstalled >&log.txt
    cat log.txt
    PACKAGE_ID=$(sed -n "/${CC_NAME}_${VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)

- `peer lifecycle chaincode queryinstalled >&log.txt`: This command queries the installed chaincode packages on the targeted peer and redirects the output to a file named log.txt.
- `cat` command is used to display the content of the log.txt on the terminal
- Final command extracts the package ID of a specific chaincode version from the content of the log.txt file.


### 4. Approve for Org1, Org2, ...

`setGlobalsForPeer0Org1`

    peer lifecycle chaincode approveformyorg -o localhost:7050 \
            --ordererTLSHostnameOverride orderer.example.com --tls \
            --collections-config $PRIVATE_DATA_CONFIG \
            --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${VERSION} \
            --init-required --package-id ${PACKAGE_ID} \
            --sequence ${VERSION}

- `peer lifecycle chaincode approveformyorg` : This command is used to approve a chaincode definition for a specific organization
- `--tls`: This flag indicates that TLS (Transport Layer Security) should be used for secure communication.
- `--package-id ${PACKAGE_ID}`: This flag specifies the package ID of the chaincode that was previously packaged and installed.
- `--sequence ${VERSION}`: This flag sets the sequence number for the chaincode.

All the other flags are same as the ones we discussed during channel creation.

*We should approve this for all orgs by calling this command after setting globals for endorsing peer of each organization*


### 5. Check commit readyness

    peer lifecycle chaincode checkcommitreadiness \
            --collections-config $PRIVATE_DATA_CONFIG \
            --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${VERSION} \
            --sequence ${VERSION} --output json --init-required

- This command is used to verify if a chaincode definition is ready to be committed and instantiated on a channel
- `--collections-config $PRIVATE_DATA_CONFIG`: This flag specifies the path to a file that contains the private data collection configuration for the chaincode. Private data collections allow selective sharing of data among specific participants
- `--init-required`: This flag indicates that the chaincode requires an initialization step during instantiation.

Other flags are either self explanatory or ones we discussed before


### 6. Commit chaincode definitions

`setGlobalsForPeer0Org1`

    peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
            --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
            --channelID $CHANNEL_NAME --name ${CC_NAME} \
            --collections-config $PRIVATE_DATA_CONFIG \
            --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_ORG1_CA \
            --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_ORG2_CA \
            --version ${VERSION} --sequence ${VERSION} --init-required

- `peer lifecycle chaincode commit`: This command is used to commit a chaincode definition for deployment on the channel
- `peerAddresses localhost:7051` --tlsRootCertFiles $PEER0_ORG1_CA: These flags specify the peer address and its corresponding TLS root certificate for the peer belonging to "org1".

### 7. Invoking the chaincode

#### 7.1 Invoke init

    peer chaincode invoke -o localhost:7050 \
            --ordererTLSHostnameOverride orderer.example.com \
            --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
            -C $CHANNEL_NAME -n ${CC_NAME} \
            --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_ORG1_CA \
            --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_ORG2_CA \
            --isInit -c '{"Args":[]}'

- `peer chaincode invoke`: This command is used to invoke a chaincode on the blockchain network.
- `-o localhost:7050`: This flag specifies the address of the orderer service to which the invocation request will be sent. In this case, it's the orderer running on the local machine at port 7050.
- `-c '{"Args":[]}'`: This flag specifies the payload for the chaincode invocation. In this case, it's an empty array of arguments (Args). This payload is sent to the chaincode's Invoke function for processing.

We alread know the other flags.
One other thing is, we provide `peerAddresses` of endorsing peers of each organization. These are the peers we installed the chaincode on.

#### 7.2 Query data

`peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"Args":["getAllLevels"]}'`

By now, this must be self explanatory to you

#### 7.3 Invoke chaincode

We can use `-c` flag of the above command to invoke different chaincodes defined. 

- `-c '{"function": "initLedger","Args":[]}'` : this invokes initLedger function in the chaincode with no arguments passed.
- `-c '{"function": "writeLevel","Args":["80.9999_6.1111","-10.99"]}'` : thin invokes writeLevel function in the chaincode with above arguments

## Outro
For both Channel and Chaincode sections, it relies on some additional functions to set global variables and change the variables according to organizations. Below are some example codes for that for you to get an idea about that

    export  CORE_PEER_TLS_ENABLED=true
    export  ORDERER_CA=${PWD}/artifacts/channel/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
    export  PEER0_ORG1_CA=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt    
    export  PEER0_ORG2_CA=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
    
    export  FABRIC_CFG_PATH=${PWD}/artifacts/channel/config/      
    
    export  PRIVATE_DATA_CONFIG=${PWD}/artifacts/private-data/collections_config.json  
      
    export  CHANNEL_NAME=mychannel
    
      
    
    setGlobalsForOrderer() {
	    export  CORE_PEER_LOCALMSPID="OrdererMSP"
	    export  CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/artifacts/channel/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
	    export  CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/msp
    }
    
      
    
    setGlobalsForPeer0Org1() {
	    export  CORE_PEER_LOCALMSPID="Org1MSP"
	    export  CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
	    export  CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
	    export  CORE_PEER_ADDRESS=localhost:7051
    }
    
    setGlobalsForPeer1Org1() {
	    export  CORE_PEER_LOCALMSPID="Org1MSP"
	    export  CORE_PEER_TLS_ROOTCERT_FILE=$PEER1_ORG1_CA
	    export  CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
	    export  CORE_PEER_ADDRESS=localhost:8051
    }
      
    setGlobalsForPeer0Org2() {
	    export  CORE_PEER_LOCALMSPID="Org2MSP"
	    export  CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
	    export  CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
	    export  CORE_PEER_ADDRESS=localhost:9051
    }
    
    setGlobalsForPeer1Org2() {
	    export  CORE_PEER_LOCALMSPID="Org2MSP"
	    export  CORE_PEER_TLS_ROOTCERT_FILE=$PEER1_ORG2_CA
	    export  CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
	    export  CORE_PEER_ADDRESS=localhost:10051
    }
