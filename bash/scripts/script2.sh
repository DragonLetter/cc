#!/bin/bash
# Copyright London Stock Exchange Group All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
echo
echo " ____    _____      _      ____    _____           _____   ____    _____ "
echo "/ ___|  |_   _|    / \    |  _ \  |_   _|         | ____| |___ \  | ____|"
echo "\___ \    | |     / _ \   | |_) |   | |    _____  |  _|     __) | |  _|  "
echo " ___) |   | |    / ___ \  |  _ <    | |   |_____| | |___   / __/  | |___ "
echo "|____/    |_|   /_/   \_\ |_| \_\   |_|           |_____| |_____| |_____|"
echo

CHANNEL_NAME="$1"
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

createChannel() {
	setGlobals 0

        if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
		peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx >&log.txt
	else
		peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
	fi
	res=$?
	cat log.txt
	verifyResult $res "Channel creation failed"
	echo "===================== Channel \"$CHANNEL_NAME\" is created successfully ===================== "
	echo
}

createB1Channel() {
        setGlobals 0
	peer channel create -o orderer.example.com:7050 -c b1b2c1 -f ./channel-artifacts/b1b2c1.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
	peer channel join -b b1b2c1.block  >&log.txt
        peer channel create -o orderer.example.com:7050 -c b1b3c1 -f ./channel-artifacts/b1b3c1.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
        peer channel join -b b1b3c1.block  >&log.txt
        peer channel create -o orderer.example.com:7050 -c b1c1 -f ./channel-artifacts/b1c1.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
        peer channel join -b b1c1.block  >&log.txt
	cat log.txt
}
createB2Channel() {
        setGlobals 1
	peer channel join -b b1b2c1.block  >&log.txt
        peer channel create -o orderer.example.com:7050 -c b2b3c1 -f ./channel-artifacts/b2b3c1.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
        peer channel join -b b2b3c1.block  >&log.txt
        peer channel create -o orderer.example.com:7050 -c b2c1 -f ./channel-artifacts/b2c1.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
        peer channel join -b b2c1.block  >&log.txt
        cat log.txt
}
createB3Channel() {
        setGlobals 2
        peer channel join -b b1b3c1.block  >&log.txt
        peer channel join -b b2b3c1.block  >&log.txt
        peer channel create -o orderer.example.com:7050 -c b3c1 -f ./channel-artifacts/b3c1.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
        peer channel join -b b3c1.block  >&log.txt
        cat log.txt
}
createB4Channel() {
        setGlobals 3
        peer channel join -b b1b2c1.block  >&log.txt
        peer channel join -b b1b3c1.block  >&log.txt
        peer channel join -b b2b3c1.block  >&log.txt
        peer channel join -b b1c1.block  >&log.txt
        peer channel join -b b2c1.block  >&log.txt
        peer channel join -b b3c1.block  >&log.txt
        cat log.txt
}

updateAnchorPeers() {
        PEER=$1
        setGlobals $PEER

        if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
		peer channel update -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx >&log.txt
	else
		peer channel update -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
	fi
	res=$?
	cat log.txt
	verifyResult $res "Anchor peer update failed"
	echo "===================== Anchor peers for org \"$CORE_PEER_LOCALMSPID\" on \"$CHANNEL_NAME\" is updated successfully ===================== "
	sleep 5
	echo
}

## Sometimes Join takes time hence RETRY atleast for 5 times
joinWithRetry () {
	peer channel join -b $CHANNEL_NAME.block  >&log.txt
	res=$?
	cat log.txt
	if [ $res -ne 0 -a $COUNTER -lt $MAX_RETRY ]; then
		COUNTER=` expr $COUNTER + 1`
		echo "PEER$1 failed to join the channel, Retry after 2 seconds"
		sleep 2
		joinWithRetry $1
	else
		COUNTER=1
	fi
        verifyResult $res "After $MAX_RETRY attempts, PEER$ch has failed to Join the Channel"
}

joinChannel () {
	for ch in 0 1 2 3; do
		setGlobals $ch
		joinWithRetry $ch
		echo "===================== PEER$ch joined on the channel \"$CHANNEL_NAME\" ===================== "
		sleep 2
		echo
	done
}

installChaincode () {
	PEER=$1
	setGlobals $PEER
	peer chaincode install -n mylc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/LetterOfCredit >&log.txt
	peer chaincode install -n bcs -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/LetterOfCredit >&log.txt
	peer chaincode install -n mycc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02 >&log.txt
	res=$?
	cat log.txt
        verifyResult $res "Chaincode installation on remote peer PEER$PEER has Failed"
	echo "===================== Chaincode is installed on remote peer PEER$PEER ===================== "
	echo
}

instantiateChaincode () {
	PEER=$1
	setGlobals $PEER
	# while 'peer chaincode' command can get the orderer endpoint from the peer (if join was successful),
	# lets supply it directly as we know it using the "-o" option
	if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
		peer chaincode instantiate -o orderer.example.com:7050 -C $CHANNEL_NAME -n mylc -v 1.0 -c '{"Args":[]}' >&log.txt
		peer chaincode instantiate -o orderer.example.com:7050 -C $CHANNEL_NAME -n bcs -v 1.0 -c '{"Args":[]}' >&log.txt
	else
		peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n mylc -v 1.0 -c '{"Args":[]}'  >&log.txt
		peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n bcs -v 1.0 -c '{"Args":[]}'  >&log.txt
		peer chaincode instantiate -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C $CHANNEL_NAME -n mycc -v 1.0 -c '{"Args":["init","a","100","b","200"]}' >&log.txt
	fi
	res=$?
	cat log.txt
	verifyResult $res "Chaincode instantiation on PEER$PEER on channel '$CHANNEL_NAME' failed"
	echo "===================== Chaincode Instantiation on PEER$PEER on channel '$CHANNEL_NAME' is successful ===================== "
	echo
}

chaincodeQuery () {
  PEER=$1
  echo "===================== Querying on PEER$PEER on channel '$CHANNEL_NAME'... ===================== "
  setGlobals $PEER
  local rc=1
  local starttime=$(date +%s)

  # continue to poll
  # we either get a successful response, or reach TIMEOUT
  while test "$(($(date +%s)-starttime))" -lt "$TIMEOUT" -a $rc -ne 0
  do
     sleep 3
     echo "Attempting to Query PEER$PEER ...$(($(date +%s)-starttime)) secs"
     peer chaincode query -C $CHANNEL_NAME -n mylc -c '{"Args":["getLcByOwner","BOC"]}' >&log.txt
     peer chaincode query -C $CHANNEL_NAME -n bcs -c '{"Args":["getBCSList","Sign"]}' >&log.txt
     test $? -eq 0 && VALUE=$(cat log.txt | awk '/Query Result/ {print $NF}')
     test "$VALUE" = "$2" && let rc=0
  done
  echo
  cat log.txt
  if test $rc -eq 0 ; then
	echo "===================== Query on PEER$PEER on channel '$CHANNEL_NAME' is successful ===================== "
  else
	echo $VALUE
  fi
}

chaincodeInvoke () {
	PEER=$1
	setGlobals $PEER
	# while 'peer chaincode' command can get the orderer endpoint from the peer (if join was successful),
	# lets supply it directly as we know it using the "-o" option
	if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
		peer chaincode invoke -o orderer.example.com:7050 -C $CHANNEL_NAME -n mycc -c '{"Args":["invoke","a","b","10"]}' >&log.txt
	else
		peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["saveLCApplication", "trans003", "{\"No\":\"AF001\", \"Applicant\":{\"No\":\"1\",\"Name\":\"Lixiaohu\",\"Domain\":\"org1.example.com\",\"Account\":\"6222000017854934\",\"DepositBank\":\"ICBC1\",\"Address\":\"1 Fuxingmen Nei Dajie, Beijing 100818, China BKCH CN BJ\",\"PostCode\":\"10010\",\"Telephone\":\"20010\",\"Telefax\":\"30010\"}, \"Beneficiary\":{\"No\":\"2\",\"Name\":\"zengyi\",\"Domain\":\"org1.example.com\",\"Account\":\"62220000178123\",\"DepositBank\":\"BCC\",\"Address\":\"8 Yabao Lu, Chaoyang District, Beijing 100020, China BKCH CN BJ 110\",\"PostCode\":\"10010\",\"Telephone\":\"20010\",\"Telefax\":\"30010\"},\"IssuingBank\":{\"No\":\"ICBC1\",\"Name\":\"ICBC Bank\",\"Domain\":\"org1.example.com\",\"Address\":\"Beijing University of Posts and Telecommunications No.10 Xitucheng Road, Haidian District, Beijing, China\",\"AccountNo\":\"6222000017854934\",\"AccountName\":\"LIXIAOHU\",\"Remark\":\"nothing\",\"PostCode\":\"10010\",\"Telephone\":\"20010\",\"Telefax\":\"30010\"}, \"AdvisingBank\":{\"No\":\"CMB123\",\"Name\":\"ICBC Bank\",\"Domain\":\"org1.example.com\",\"Address\":\"Beijing University of Posts and Telecommunications No.10 Xitucheng Road, Haidian District, Beijing, China\",\"AccountNo\":\"6222000017854934\",\"AccountName\":\"LIXIAOHU\",\"Remark\":\"nothing\",\"PostCode\":\"10010\",\"Telephone\":\"20010\",\"Telefax\":\"30010\"}, \"ExpiryDate\":\"2017-12-31T16:09:51.692226358+08:00\", \"ExpiryPlace\":\"Beijing\", \"IsAtSight\":\"false\", \"AfterSight\":\"90\", \"GoodsInfo\":{\"GoodsNo\":\"1234567\",\"AllowPartialShipment\":\"1\", \"AllowTransShipment\":\"0\", \"LatestShipmentDate\":\"2017-08-10T16:09:51.692226358+08:00\", \"ShippingWay\":\"airplane\", \"ShippingPlace\":\"shanghai\", \"ShippingDestination\":\"beijing\", \"TradeNature\":\"1\", \"GoodsDescription\":\"goodsDescription\"}, \"DocumentRequire\":\"1\",\"Currency\":\"CNY\", \"Amount\":\"123.12\", \"ApplyTime\":\"2017-11-09T16:09:51.692226358+08:00\", \"ChargeInIssueBank\":\"1\", \"ChargeOutIssueBank\":\"1\", \"DocDelay\":\"90\", \"OtherRequire\":\"none\", \"Contract\":{\"FileName\":\"lixiaohucontract\",\"FileUri\":\"c:qwerty\",\"FileHash\":\"56d34dad234bbdabcb3213\",\"FileSignature\":\"lixiaohussignature\",\"Uploader\":\"lihuichi\"}, \"Attachments\":[{\"FileName\":\"dragonLedgerDoc\",\"FileUri\":\"file_path\",\"FileHash\":\"cer324234fsfdgeergfgdfd\",\"FileSignature\":\"lixiaohuSign\",\"Uploader\":\"lixiaohu\"},{\"FileName\":\"dragonLedgerDoc1\",\"FileUri\":\"file_path1\",\"FileHash\":\"cer324234fsfdgeergfgdfd1\",\"FileSignature\":\"lixiaohuSign1\",\"Uploader\":\"chaincodelixiaohu1\"}]}"]}' >&log.txt
		peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n bcs -c'{"Args":["saveBCSInfo", "bcs001", "{\"No\":\"B001\",\"Type\":\"Bank\", \"DataBank\":{\"No\":\"ICBC1\",\"Name\":\"ICBC Bank\",\"Domain\":\"org1.example.com\",\"Address\":\"Beijing University of Posts and Telecommunications No.10 Xitucheng Road, Haidian District, Beijing, China\",\"AccountNo\":\"6222000017854934\",\"AccountName\":\"LIXIAOHU\",\"Remark\":\"nothing\",\"PostCode\":\"10010\",\"Telephone\":\"20010\",\"Telefax\":\"30010\"},\"DataCorp\":{},\"StateSign\":0}"]}' >&log.txt
	fi
	res=$?
	cat log.txt
	verifyResult $res "Invoke execution on PEER$PEER failed "
	echo "===================== Invoke transaction on PEER$PEER on channel '$CHANNEL_NAME' is successful ===================== "
	echo
}

## Create channel
echo "Creating channel..."
createChannel

## Join all the peers to the channel
echo "Having all peers join the channel..."
joinChannel

createB1Channel
createB2Channel
createB3Channel
createB4Channel
## Set the anchor peers for each org in the channel
#echo "Updating anchor peers for B1..."
#updateAnchorPeers 0
#echo "Updating anchor peers for B3..."
#updateAnchorPeers 2

## Install chaincode on Peer0/Org1 and Peer2/Org2
echo "Installing chaincode on B1/B2/B3..."
installChaincode 0
installChaincode 1
installChaincode 2
echo "Install chaincode on C1..."
installChaincode 3

#Instantiate chaincode on Peer2/Org2
echo "Instantiating chaincode on B1/peer2..."
instantiateChaincode 0

#Query on chaincode on Peer0/Org1
echo "Querying chaincode on B1/peer0..."
chaincodeQuery 0 100

#Invoke on chaincode on Peer0/Org1
echo "Sending invoke transaction on B1/peer0..."
chaincodeInvoke 0

## Install chaincode on Peer3/Org2
#echo "Installing chaincode on B3/peer3..."
#installChaincode 3

#Query on chaincode on Peer3/Org2, check if the result is 90
echo "Querying chaincode on B2..."
chaincodeQuery 1 90
echo "Querying chaincode on B3..."
chaincodeQuery 2 90
echo "Querying chaincode on C1..."
chaincodeQuery 3 90
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
