sshpass -p Goodluck7 ssh fabric@p1 "cd /home/fabric/go/src/github.com/hyperledger/fabric/examples/e2e_cli;./network_setup.sh down"
sshpass -p Goodluck7 ssh fabric@p2 "cd /home/fabric/go/src/github.com/hyperledger/fabric/examples/e2e_cli;./network_setup.sh down"
sshpass -p Goodluck7 ssh fabric@p3 "cd /home/fabric/go/src/github.com/hyperledger/fabric/examples/e2e_cli;./network_setup.sh down"
sshpass -p Goodluck7 ssh fabric@p4 "cd /home/fabric/go/src/github.com/hyperledger/fabric/examples/e2e_cli;./network_setup.sh down"
cd /home/fabric/go/src/github.com/hyperledger/fabric/examples/e2e_cli
./network_setup.sh down
