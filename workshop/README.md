# Configure Hyperledger Fabric network and deploy chaincode
This page will try to guide you step by step through the process to create, install and deploy a Hyperledger Fabric network using Docker containers.
Also, include chaincode examples and will help you to install and instantiate those chaincodes in the network on both, go and nodejs, languages.

## Download platform specific binaries 
```
curl -sSL https://goo.gl/6wtTN5 | bash -s 1.1.0
```
When finish, navigate to the bin folder and verify the version of the binaries downloaded, it should be 1.1.0
```
$ cryptogen version
``` 
Its also better to add the binaries directory to the path in order to access them from any directory, one way to do so in Linux is copy the path to bin folder, open user profile file and append the directory to the PATH variable:
```
$ nano ~/.profile
```
Set the path variable in the file:
```
PATH="$PATH:/yourpathtodir/fabric-samples/bin"
```
You need to do logout and login again to update the entire system with this variable.

## Create crypto materials

At this point you should have downloaded this repo.
Navigate to **workshop** directory and run the following command assuming thath
```
$ cryptogen generate --config=./fabricnetwork/crypto-config.yaml
```
On success this will create a crypto-config directory with all the cryptographic material.

## Create channel artifacts
In this step we will create the required artifacts to configure the channel in the Fabric network and generate the "genesis" block.
To use the required tool it is necessary to define a variable to indicate where the configuration file is located, enter to **fabricnetwork** directory an run the following command
```
$ export FABRIC_CFG_PATH=$PWD
``` 
This will indicate that the file is in the current directory.
It will be better if you define a folder where you will put all the generated files, create a **channel-artifacts** directory:
```
$ mkdir channel-artifacts
```
Then we will generate the **genesis** block, this is the first block of our blockchain:
```
$ configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
```
Next, create the channel configuration transaction
```
$ configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID mychannel
```
Now, we need to define the anchor peers for Org1 and Org2:
```
$ configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP
$ configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID mychannel -asOrg Org2MSP
```
We are ready now to deploy our Fabric network.

## Configure docker-compose file
To development is more easy to deploy the network on docker containers, so this project contains a composer file to deploy those containers.
Also, in order to make more easy the configuration of this file, includes an script to set the directories required to deploy our network. This file is located on the same **workshop** directory and is named **preparedocker.sh**.
For this workshop, there is only 4 variables to defined:
 **CRYPTO_CONFIGROOT**: This variable needs to be set to the **crypto-config** directory, for this an the anothers variables is required to set the absolute path.
 **CHANN_ART_DIR**: Set this to the **channel-artifacts** directory.
 **ORDERER_CERTS**: This will be the absolute path to the orderer organization in crypto-config dir: (crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com)
 **CHAINCODE**: The path to the chaincode directory, in our case, this is located in workshop/chaincode directory.
Once this variables are set run the preparedocker.sh script (make sure that has execution permissions):
```
$ ./preparedocker.sh
```

## Start the network
Now we are ready to start the network, locate on **fabricnetwork** directory and run the next command:
```
$ docker-compose up
```
When finish you will have the docker containers running. On the first time, this may take a long time because it will download all the images required.

## Install the chaincode
To run the next commands it will be necessary to interact with one of the containers up and running, the CLI. Because it is initially configured for the Org1, we need to set some environment variables to change between one organization and another one.
To enter in the cli, run the next command:
```
$ docker exec -it cli bash
```
To make more easy and shorter the next command, we will set a new environment variable to set where the orderer certificate is located, consider that  this paths are related to the docker container file structure, not the local one.
```
$ TLSCADIR=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

## Configure the network
With the files generated previously we will create and configure the channel in the network.
### Create the channel
```
$ peer channel create -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/channel.tx --tls --cafile $TLSCADIR
```
This command will create a file. Now we will join the peer for the Organization 1 to the channel:
```
$ peer channel join -b mychannel.block
```
and update the anchor peer
```
peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/Org1MSPanchors.tx --tls --cafile $TLSCADIR
```
Now we need to change the Variables to interact with the Org2 peer.
```
$ CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
$ CORE_PEER_ADDRESS=peer0.org2.example.com:7051
$ CORE_PEER_LOCALMSPID="Org2MSP"
$ CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
```
Now we are ready to join the peer to the channel and update the anchor peer.
```
$ peer channel join -b mychannel.block
$ peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/Org2MSPanchors.tx --tls --cafile $TLSCADIR
```
With this, we are running our first blockchain network using Hyperledger Fabric.

## Install and instantiate the chaincode
The next steps are for install the chaincode on the network.
First we need to install the chaincode on the peer of the Org1. Set the variables in the CLI to interact with this peer:
 ```
$ CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
$ CORE_PEER_ADDRESS=peer0.org1.example.com:7051
$ CORE_PEER_LOCALMSPID="Org1MSP"
$ CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
```
### Install the chaincode
```
$ peer chaincode install -n licensecc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go
```

> **NodeJs** To install a nodejs chaincode we need to specify that on the command, also we need the absolute path to the directory
```
$ peer chaincode install -n licensecc -v 1.0 -l node -p /opt/gopath/src/github.com/hyperledger/fabric/examples/chaincode/javascript
```

### Instantiate the chaincode in the channel
```
$ peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile $TLSCADIR -C mychannel -n licensecc -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"
```
Now, the chaincode is ready to use in the network, but because it is only installed on the peer for the Organization 1, all the requests needs to be from this peer.

> **NodeJs** To instantiate a nodejs chaincode we need to specify that on the command, also we need the absolute path to the directory. If you previously deployed a go chaincode, it is probable that you face some troubles to deploy a nodejs chaincode. The best and easy solution is to delete all images and download again (you can do this running docker system prune -a, take in account that this command will delete ALL images)
```
$ peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile $TLSCADIR -C mychannel -n licensecc -l node -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"
```

## Interact with the chaincode 
The next step is interact with the chaincode. This will vary according your chaincode definition. This commands are only for the chaincode presented in this repository.

Set the initial data:
```
$ peer chaincode invoke -o orderer.example.com:7050  --tls --cafile $TLSCADIR  -C mychannel -n licensecc -c '{"Args":["create"]}'
```
Query data:
```
$ peer chaincode query -C mychannel -n licensecc -c '{"Args":["query","abc123"]}'
```
Change ownership
```
$ peer chaincode invoke -o orderer.example.com:7050  --tls --cafile $TLSCADIR  -C mychannel -n licensecc -c '{"Args":["transfer","abc123","moises"]}'
```

# Hyperledger Composer (TODO)
## INSTALL PREREQUISITES

* Grant execute permissions to file "installcomposer.sh"
```
$ chmod +rwx ./installcomposer.sh
```
* Run the script
```
$ HOME_PATH=/replace_your_home_path ./installcomposer.sh
```
$nbsp;
### Notes
* Do not run the script with root user or using sudo command, the script will prompt the root password when need it.
* The home path should be something like this "/home/nst"
* Once the script finish, logout and login again to refresh some path variables