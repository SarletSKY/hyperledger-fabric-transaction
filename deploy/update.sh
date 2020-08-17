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

echo "安装链码"
docker exec cli peer chaincode install -n food -v 1.3.0 -l golang -p github.com/chaincode/perishable-food
docker exec -e CORE_PEER_ADDRESS=word.node1.gdzc.com:7051 cli peer chaincode install -n food -v 1.3.0 -l golang -p github.com/chaincode/perishable-food

docker exec \
-e CORE_PEER_ADDRESS=zhao.node2.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org2MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG2_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER0_ORG2_CA} \
cli peer chaincode install -n food -v 1.3.0 -l golang -p github.com/chaincode/perishable-food

docker exec \
-e CORE_PEER_ADDRESS=weixiong.node2.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org2MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG2_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER1_ORG2_CA} \
cli peer chaincode install -n food -v 1.3.0 -l golang -p github.com/chaincode/perishable-food

docker exec \
-e CORE_PEER_ADDRESS=peer0.node3.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org3MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG3_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER0_ORG3_CA} \
cli peer chaincode install -n food -v 1.3.0 -l golang -p github.com/chaincode/perishable-food

docker exec \
-e CORE_PEER_ADDRESS=peer1.node3.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org3MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG3_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER1_ORG3_CA} \
cli peer chaincode install -n food -v 1.3.0 -l golang -p github.com/chaincode/perishable-food

echo "链码升级"
docker exec cli peer chaincode upgrade \
--tls --cafile ${ORDERER_CAFILE} -o center.gdzc.com:7050 -n food -v 1.3.0 -l golang -C transaction \
-c '{"Args":[""]}' -P 'AND("Org1MSP.member","Org2MSP.member","Org3MSP.member")'
docker exec -e CORE_PEER_ADDRESS=word.node1.gdzc.com:7051 \
cli peer chaincode upgrade \
--tls --cafile ${ORDERER_CAFILE} -o center.gdzc.com:7050 -n food -v 1.3.0 -l golang -C transaction \
-c '{"Args":[""]}' -P 'AND("Org1MSP.member","Org2MSP.member","Org3MSP.member")'

docker exec \
-e CORE_PEER_ADDRESS=zhao.node2.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org2MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG2_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER0_ORG2_CA} \
cli peer chaincode upgrade \
--tls --cafile ${ORDERER_CAFILE} -o center.gdzc.com:7050 -n food -v 1.3.0 -l golang -C transaction \
-c '{"Args":[""]}' -P 'AND("Org1MSP.member","Org2MSP.member","Org3MSP.member")'

docker exec \
-e CORE_PEER_ADDRESS=weixiong.node2.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org2MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG2_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER1_ORG2_CA} \
cli peer chaincode upgrade \
--tls --cafile ${ORDERER_CAFILE} -o center.gdzc.com:7050 -n food -v 1.3.0 -l golang -C transaction \
-c '{"Args":[""]}' -P 'AND("Org1MSP.member","Org2MSP.member","Org3MSP.member")'

docker exec \
-e CORE_PEER_ADDRESS=peer0.node3.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org3MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG3_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER0_ORG3_CA} \
cli peer chaincode upgrade \
--tls --cafile ${ORDERER_CAFILE} -o center.gdzc.com:7050 -n food -v 1.3.0 -l golang -C transaction \
-c '{"Args":[""]}' -P 'AND("Org1MSP.member","Org2MSP.member","Org3MSP.member")'

docker exec \
-e CORE_PEER_ADDRESS=peer1.node3.gdzc.com:7051 \
-e CORE_PEER_LOCALMSPID=Org3MSP \
-e CORE_PEER_MSPCONFIGPATH=${ORG3_MSPCONFIGPATH} \
-e CORE_PEER_TLS_ROOTCERT_FILE=${PEER1_ORG3_CA} \
cli peer chaincode upgrade \
--tls --cafile ${ORDERER_CAFILE} -o center.gdzc.com:7050 -n food -v 1.3.0 -l golang -C transaction \
-c '{"Args":[""]}' -P 'AND("Org1MSP.member","Org2MSP.member","Org3MSP.member")'