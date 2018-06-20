TIMEOUT=6000
CHANNEL_NAME=mychannel
ORDERER_DOMAIN=example.com
ORG1=org1.example.com
ORG2=org2.example.com
CRYPTO_CONFIGROOT=/home/moises/Documents/NST-Blockchain4Business/workshop/crypto-config
CHANN_ART_DIR=/home/moises/Documents/NST-Blockchain4Business/workshop/fabricnetwork/channel-artifacts
ORDERER_CERTS=/home/moises/Documents/NST-Blockchain4Business/workshop/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com
CHAINCODE=/home/moises/Documents/NST-Blockchain4Business/workshop/chaincode

sed -i -e 's%${CHANN_ART_DIR}%'"$CHANN_ART_DIR"'%g' docker-compose.yaml
sed -i -e 's%${ORDERER_CERTS}%'"$ORDERER_CERTS"'%g' docker-compose.yaml
sed -i -e 's%${CRYPTO_CONFIGROOT}%'"$CRYPTO_CONFIGROOT"'%g' docker-compose.yaml
sed -i -e 's%${ORG1}%'"$ORG1"'%g' docker-compose.yaml
sed -i -e 's%${ORG2}%'"$ORG2"'%g' docker-compose.yaml
sed -i -e 's%${CHANNEL_NAME}%'"$CHANNEL_NAME"'%g' docker-compose.yaml
sed -i -e 's%${TIMEOUT}%'"$TIMEOUT"'%g' docker-compose.yaml
sed -i -e 's%${CHAINCODE}%'"$CHAINCODE"'%g' docker-compose.yaml
