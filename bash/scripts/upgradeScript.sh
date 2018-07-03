#!/bin/bash
# Copyright London Stock Exchange Group All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
echo
echo " ____    _____      _      ____    _____   "
echo "/ ___|  |_   _|    / \    |  _ \  |_   _|  "
echo "\___ \    | |     / _ \   | |_) |   | |    "
echo " ___) |   | |    / ___ \  |  _ <    | |    "
echo "|____/    |_|   /_/   \_\ |_| \_\   |_|    "
echo

CHANNEL_NAME="$2"
: ${CHANNEL_NAME:="all"}
: ${TIMEOUT:="60"}
COUNTER=1
MAX_RETRY=5
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

echo "Channel name : "$CHANNEL_NAME

verifyResult () {
	if [ $1 -ne 0 ] ; then
		echo "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
                echo "================== ERROR !!! FAILED to execute End-2-End Scenario =================="
		echo
   		exit 1
	fi
}

setGlobals () {

	if [ $1 -eq 0 ] ; then
		CORE_PEER_LOCALMSPID="B1MSP"
		CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/B1.example.com/peers/peer0.B1.example.com/tls/ca.crt
		CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/B1.example.com/users/Admin@B1.example.com/msp
	        CORE_PEER_ADDRESS=peer0.B1.example.com:7051
	fi

        if [ $1 -eq 1 ] ; then
                CORE_PEER_LOCALMSPID="B2MSP"
                CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/B2.example.com/peers/peer0.B2.example.com/tls/ca.crt
                CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/B2.example.com/users/Admin@B2.example.com/msp
                CORE_PEER_ADDRESS=peer0.B2.example.com:7051
        fi

        if [ $1 -eq 2 ] ; then
                CORE_PEER_LOCALMSPID="B3MSP"
                CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/B3.example.com/peers/peer0.B3.example.com/tls/ca.crt
                CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/B3.example.com/users/Admin@B3.example.com/msp
                CORE_PEER_ADDRESS=peer0.B3.example.com:7051
        fi

        if [ $1 -eq 3 ] ; then
                CORE_PEER_LOCALMSPID="C1MSP"
                CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/C1.example.com/peers/peer0.C1.example.com/tls/ca.crt
                CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/C1.example.com/users/Admin@C1.example.com/msp
                CORE_PEER_ADDRESS=peer0.C1.example.com:7051
        fi

	env |grep CORE
}

installChaincode () {
	PEER=$1
	setGlobals $PEER
	peer chaincode install -n mylc -v $2 -p github.com/hyperledger/fabric/examples/chaincode/go/LetterOfCredit >&log.txt
	#peer chaincode install -n mycc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02 >&log.txt
	res=$?
	cat log.txt
        verifyResult $res "Chaincode installation on remote peer PEER$PEER has Failed"
	echo "===================== Chaincode is installed on remote peer PEER$PEER ===================== "
	echo
}

upgradeChaincode () {
	PEER=$1
	setGlobals $PEER
	# while 'peer chaincode' command can get the orderer endpoint from the peer (if join was successful),
	# lets supply it directly as we know it using the "-o" option
	if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
		peer chaincode upgrade -o orderer.example.com:7050 -C $CHANNEL_NAME -n mylc -v $2 -c '{"Args":[]}' >&log.txt
	else
		peer chaincode upgrade -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n mylc -v $2 -c '{"Args":[]}'  >&log.txt
		#peer chaincode upgrade -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C $CHANNEL_NAME -n mycc -v $2 -c '{"Args":["init","a","100","b","200"]}' >&log.txt
	fi
	res=$?
	cat log.txt
	verifyResult $res "Chaincode instantiation on PEER$PEER on channel '$CHANNEL_NAME' failed"
	echo "===================== Chaincode upgraded on PEER$PEER on channel '$CHANNEL_NAME' is successful ===================== "
	echo
}

## Install chaincode on Peer0/Org1 and Peer2/Org2
echo "Installing chaincode on B1/B2/B3..."
installChaincode 0 $1
installChaincode 1 $1
installChaincode 2 $1
echo "Install chaincode on C1..."
installChaincode 3 $1

#upgrade chaincode on Peer2/Org2
echo "upgrading chaincode on B1/peer2..."
upgradeChaincode 0 $1

## Install chaincode on Peer3/Org2
#echo "Installing chaincode on B3/peer3..."
#installChaincode 3 $1

echo
echo "===================== All GOOD, End-2-End execution completed ===================== "
echo

echo
echo " _____   _   _   ____            _____   ____    _____ "
echo "| ____| | \ | | |  _ \          | ____| |___ \  | ____|"
echo "|  _|   |  \| | | | | |  _____  |  _|     __) | |  _|  "
echo "| |___  | |\  | | |_| | |_____| | |___   / __/  | |___ "
echo "|_____| |_| \_| |____/          |_____| |_____| |_____|"
echo

exit 0
