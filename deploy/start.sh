#!/bin/bash

COMPOSE_FILE=docker-compose-cli.yaml
ORDERER_CAFILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/gdzc.com/orderers/center.gdzc.com/msp/tlscacerts/tlsca.gdzc.com-cert.pem

ORG2_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/node2.gdzc.com/users/Admin@node2.gdzc.com/msp
ORG3_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/node3.gdzc.com/users/Admin@node3.gdzc.com/msp

PEER0_ORG1_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/node1.gdzc.com/peers/hello.node1.gdzc.com/tls/ca.crt
PEER1_ORG1_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/node1.gdzc.com/peers/word.node1.gdzc.com/tls/ca.crt
PEER0_ORG2_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/node2.gdzc.com/peers/zhao.node2.gdzc.com/tls/ca.crt
PEER1_ORG2_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/node2.gdzc.com/peers/weixiong.node2.gdzc.com/tls/ca.crt
PEER0_ORG3_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/node3.gdzc.com/peers/peer0.node3.gdzc.com/tls/ca.crt
PEER1_ORG3_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/node3.gdzc.com/peers/peer1.node3.gdzc.com/tls/ca.crt

echo "环境清理"
mkdir -p crypto-config
mkdir -p channel-artifacts
rm -fr crypto-config/*
rm -fr channel-artifacts/×

echo "生成初始区块与证书"
cryptogen generate --config=./crypto-config.yaml
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block

echo "生成tx文件与更新锚节点的文件"
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/transaction.tx -channelID transaction

configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID transaction -asOrg Org1MSP
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID transaction -asOrg Org2MSP
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org3MSPanchors.tx -channelID transaction -asOrg Org3MSP


docker-compose -f $COMPOSE_FILE up -d

echo "创建通道"
docker exec cli peer channel create -o center.gdzc.com:7050 -c transaction --tls --cafile $ORDERER_CAFILE -f ./channel-artifacts/transaction.tx

echo "加入通道"
docker exec cli peer channel join -b transaction.block
docker exec -e CORE_PEER_ADDRESS=word.node1.gdzc.com:7051 cli peer channel join -b transaction.block

docker exec \
-e CORE_PEER_ADDRESS=zhao.node2.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org2MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG2_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER0_ORG2_CA} \
cli peer channel join -b transaction.block

docker exec \
-e CORE_PEER_ADDRESS=weixiong.node2.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org2MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG2_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER1_ORG2_CA} \
cli peer channel join -b transaction.block

docker exec \
-e CORE_PEER_ADDRESS=peer0.node3.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org3MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG3_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER0_ORG3_CA} \
cli peer channel join -b transaction.block

docker exec \
-e CORE_PEER_ADDRESS=peer1.node3.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org3MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG3_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER1_ORG3_CA} \
cli peer channel join -b transaction.block

echo "安装链码"
docker exec cli peer chaincode install -n food -v 1.0.0 -l golang -p github.com/chaincode/perishable-food
docker exec -e CORE_PEER_ADDRESS=word.node1.gdzc.com:7051 cli peer chaincode install -n food -v 1.0.0 -l golang -p github.com/chaincode/perishable-food

docker exec \
-e CORE_PEER_ADDRESS=zhao.node2.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org2MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG2_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER0_ORG2_CA} \
cli peer chaincode install -n food -v 1.0.0 -l golang -p github.com/chaincode/perishable-food

docker exec \
-e CORE_PEER_ADDRESS=weixiong.node2.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org2MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG2_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER1_ORG2_CA} \
cli peer chaincode install -n food -v 1.0.0 -l golang -p github.com/chaincode/perishable-food

docker exec \
-e CORE_PEER_ADDRESS=peer0.node3.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org3MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG3_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER0_ORG3_CA} \
cli peer chaincode install -n food -v 1.0.0 -l golang -p github.com/chaincode/perishable-food

docker exec \
-e CORE_PEER_ADDRESS=peer1.node3.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org3MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG3_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER1_ORG3_CA} \
cli peer chaincode install -n food -v 1.0.0 -l golang -p github.com/chaincode/perishable-food

echo "实例化链码"
docker exec cli peer chaincode instantiate --tls --cafile ${ORDERER_CAFILE} -o center.gdzc.com:7050 -n food -v 1.0.0 -l golang -C transaction \
-c '{"Args":["init","a","100","b","200"]}' -P 'AND("Org1MSP.member","Org2MSP.member","Org3MSP.member")'

#sleep 3
#
#echo "测试查询"
#docker exec cli peer chaincode query -n food -C transaction -c '{"Args":["queryCommodityList"]}'
#
#echo "转账"
#docker exec cli peer chaincode invoke -o center.gdzc.com:7050  -n food --tls --cafile ${ORDERER_CAFILE} -C transaction \
#--peerAddresses hello.node1.gdzc.com:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} \
#--peerAddresses zhao.node2.gdzc.com:7051 --tlsRootCertFiles ${PEER0_ORG2_CA} \
#--peerAddresses peer0.node3.gdzc.com:7051 --tlsRootCertFiles ${PEER0_ORG3_CA} \
#-c '{"Args":["createCommodity","2","2746eef4-7f44-4b65-a221-ca661fc0f1a2","芭蕉","中国","4","商家创建商品～芭蕉"]}'
#
#sleep 3
#echo "再次查询"
#docker exec cli peer chaincode query -n food -C transaction -c '{"Args":["queryCommodityList"]}'

