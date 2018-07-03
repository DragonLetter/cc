echo "-----Upgrade ChainCode files to all peers-----"
cp /home/fabric/plume/chaincode/LetterOfCredit /home/fabric/go/src/github.com/hyperledger/fabric/examples/chaincode/go/ -R

echo "-----Start copy files to P1-----"
./expect_scp.sh p1 fabric Goodluck7 ./channel-artifacts /home/fabric/go/src/github.com/hyperledger/fabric/examples/e2e_cli/
./expect_scp.sh p1 fabric Goodluck7 ./crypto-config /home/fabric/go/src/github.com/hyperledger/fabric/examples/e2e_cli/
./expect_scp.sh p1 fabric Goodluck7 /home/fabric/plume/chaincode/LetterOfCredit /home/fabric/go/src/github.com/hyperledger/fabric/examples/chaincode/go/

echo "-----Start copy files to P2-----"
./expect_scp.sh p2 fabric Goodluck7 ./channel-artifacts /home/fabric/go/src/github.com/hyperledger/fabric/examples/e2e_cli/
./expect_scp.sh p2 fabric Goodluck7 ./crypto-config /home/fabric/go/src/github.com/hyperledger/fabric/examples/e2e_cli/
./expect_scp.sh p2 fabric Goodluck7 /home/fabric/plume/chaincode/LetterOfCredit /home/fabric/go/src/github.com/hyperledger/fabric/examples/chaincode/go/
echo "-----Start copy files to P3-----"
./expect_scp.sh p3 fabric Goodluck7 ./channel-artifacts /home/fabric/go/src/github.com/hyperledger/fabric/examples/e2e_cli/
./expect_scp.sh p3 fabric Goodluck7 ./crypto-config /home/fabric/go/src/github.com/hyperledger/fabric/examples/e2e_cli/
./expect_scp.sh p3 fabric Goodluck7 /home/fabric/plume/chaincode/LetterOfCredit /home/fabric/go/src/github.com/hyperledger/fabric/examples/chaincode/go/
echo "-----Start copy files to P4-----"
./expect_scp.sh p4 fabric Goodluck7 ./channel-artifacts /home/fabric/go/src/github.com/hyperledger/fabric/examples/e2e_cli/
./expect_scp.sh p4 fabric Goodluck7 ./crypto-config /home/fabric/go/src/github.com/hyperledger/fabric/examples/e2e_cli/
./expect_scp.sh p4 fabric Goodluck7 /home/fabric/plume/chaincode/LetterOfCredit /home/fabric/go/src/github.com/hyperledger/fabric/examples/chaincode/go/
echo "-----All Copy done!-----"