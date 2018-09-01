# dragonpay
Dragon Ledger Letter of Credit ChainCode usage
## 一、chaincode安装及更新
### 1.1. 安装：
```bash
peer chaincode install -n mylc -v 1.3 -p github.com/hyperledger/fabric/examples/chaincode/go/LetterOfCredit
```
### 1.2. 初始化：
```bash
peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -v 1.0 -c '{"Args":[]}' -P "OR('Org1MSP.member','Org2MSP.member')"
```
### 1.3. 更新代码，说明需要执行安装操作，并将-v参数进行升级，然后执行如下操作
```bash
peer chaincode upgrade -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -v 1.2 -c '{"Args":[]}' -P "OR('Org1MSP.member','Org2MSP.member')"
```
## 二、信用证业务流转功能说明
### 2.1. 申请人保存信用证开证申请，填写基本信息，参数为3个，包括编号、信用证编号、申请单
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["saveLCApplication", "trans003", "{\"No\":\"AF001\", \"Applicant\":{\"No\":\"1\",\"Name\":\"Lixiaohu\",\"Domain\":\"org1.example.com\",\"Account\":\"6222000017854934\",\"DepositBank\":\"ICBC1\",\"Address\":\"1 Fuxingmen Nei Dajie, Beijing 100818, China BKCH CN BJ\",\"PostCode\":\"10010\",\"Telephone\":\"20010\",\"Telefax\":\"30010\"}, \"Beneficiary\":{\"No\":\"2\",\"Name\":\"zengyi\",\"Domain\":\"org1.example.com\",\"Account\":\"62220000178123\",\"DepositBank\":\"BCC\",\"Address\":\"8 Yabao Lu, Chaoyang District, Beijing 100020, China BKCH CN BJ 110\",\"PostCode\":\"10010\",\"Telephone\":\"20010\",\"Telefax\":\"30010\"},\"IssuingBank\":{\"No\":\"ICBC1\",\"Name\":\"ICBC Bank\",\"Domain\":\"org1.example.com\",\"Address\":\"Beijing University of Posts and Telecommunications No.10 Xitucheng Road, Haidian District, Beijing, China\",\"AccountNo\":\"6222000017854934\",\"AccountName\":\"LIXIAOHU\",\"Remark\":\"nothing\",\"PostCode\":\"10010\",\"Telephone\":\"20010\",\"Telefax\":\"30010\"}, \"AdvisingBank\":{\"No\":\"CMB123\",\"Name\":\"ICBC Bank\",\"Domain\":\"org1.example.com\",\"Address\":\"Beijing University of Posts and Telecommunications No.10 Xitucheng Road, Haidian District, Beijing, China\",\"AccountNo\":\"6222000017854934\",\"AccountName\":\"LIXIAOHU\",\"Remark\":\"nothing\",\"PostCode\":\"10010\",\"Telephone\":\"20010\",\"Telefax\":\"30010\"}, \"ExpiryDate\":\"2017-12-31T16:09:51.692226358+08:00\", \"ExpiryPlace\":\"Beijing\", \"IsAtSight\":\"false\", \"AfterSight\":\"90\", \"GoodsInfo\":{\"GoodsNo\":\"1234567\",\"AllowPartialShipment\":\"1\", \"AllowTransShipment\":\"0\", \"LatestShipmentDate\":\"2017-08-10T16:09:51.692226358+08:00\", \"ShippingWay\":\"airplane\", \"ShippingPlace\":\"shanghai\", \"ShippingDestination\":\"beijing\", \"TradeNature\":\"1\", \"GoodsDescription\":\"goodsDescription\"}, \"DocumentRequire\":\"1\",\"Currency\":\"CNY\", \"Amount\":\"123.12\", \"ApplyTime\":\"2017-11-09T16:09:51.692226358+08:00\", \"ChargeInIssueBank\":\"1\", \"ChargeOutIssueBank\":\"1\", \"DocDelay\":\"90\", \"OtherRequire\":\"none\", \"Contract\":{\"FileName\":\"lixiaohucontract\",\"FileUri\":\"c:qwerty\",\"FileHash\":\"56d34dad234bbdabcb3213\",\"FileSignature\":\"lixiaohussignature\",\"Uploader\":\"lihuichi\"}, \"Attachments\":[{\"FileName\":\"dragonLedgerDoc\",\"FileUri\":\"file_path\",\"FileHash\":\"cer324234fsfdgeergfgdfd\",\"FileSignature\":\"lixiaohuSign\",\"Uploader\":\"lixiaohu\"},{\"FileName\":\"dragonLedgerDoc1\",\"FileUri\":\"file_path1\",\"FileHash\":\"cer324234fsfdgeergfgdfd1\",\"FileSignature\":\"lixiaohuSign1\",\"Uploader\":\"chaincodelixiaohu1\"}]}"]}'
```
### 2.2. 申请人提交信用证开证申请，填写基本信息，参数为1个，分别是信用证交易号，执行完操作后，信用证状态为"申请中"
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["submitLCApplication", "trans002"]}'
```
### 2.3. 银行确认信用证申请，并返回一个信用证编号，参数为：5个，分别是信用证交易号、信用证编号、保证金金额、银行意见、是否同意（true false），执行完操作后，信用证状态为"草稿"
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["bankConfirmApplication", "trans002","lc20171108","200000.0","agree","true"]}'
```
### 2.4. 申请人提交保证金，同时上传相关单据,参数为：3个,分别是信用证交易号、保证金对象（已交金额、单据），执行完操作后，信用证状态为"草稿"
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["deposit", "trans002","{\"CommitAmount\":\"150000.0\",\"DepositDoc\":{\"FileName\":\"DepositDoc\",\"FileUri\":\"file_path\",\"FileHash\":\"cer324234fsfdgeergfgdfd\",\"FileSignature\":\"lixiaohuSign\",\"Uploader\":\"lixiaohu\"}}"]}'
```
### 2.5. 开证行开立信用证，上传信用证正本，参数为4个，信用证交易号，意见、是否同意（true false）、信用证文本，执行完操作后，信用证状态为"正本"
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["issueLetterOfCredit", "trans002","agree lc issue","true","{\"FileName\":\"DepositDoc\",\"FileUri\":\"file_path\",\"FileHash\":\"cer324234fsfdgeergfgdfd\",\"FileSignature\":\"lixiaohuSign\",\"Uploader\":\"lixiaohu\"}"]}'
```
### 2.6. 通知行同意开证行通知或者不同意，参数为3个，信用证交易号，意见、是否同意（true false）执行完操作后，信用证状态为"正本"
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["advisingBankReceiveLCNotice", "trans002","agree lc receive","true"]}'
```
### 2.7. 受益人同意同意或者拒绝LC内容，参数为3个，信用证交易号、意见、是否同意（true false）执行完操作后，信用证状态为"正本"
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["beneficiaryReceiveLCNotice", "trans002","agree lc receive","true"]}'
```
### 2.8. 申请人进行信用证修改操作，参数为2个，信用证交易号，修改单（修改次数、修改货币类型、修改金额、期限增减、有效日期、发货地修改、保证金增减），执行完修改操作后，信用证状态为"正本修改"
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["lcAmendSubmit", "trans002","{\"AmendTimes\":\"1\",\"AmendedCurrency\":\"USD\",\"AmendedAmt\":\"9876.54\",\"AddedDays\":\"-20\",\"AmendExpiryDate\":\"20170930\",\"TransPortName\":\"Xian\",\"AddedDepositAmt\":\"-100\"}"]}'
```
### 2.9. 开证行、通知行进行信用证修改后会签，参数为3个，信用证交易号，意见、是否同意（true false），信用证状态为"正本"
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["lcAmendConfirm", "0","buxing","false"]}'
```
### 2.10. 受益人执行交单操作，将相关单据交付给通知行。参数为3个，信用证交易号,货运单信息,提货单附件。执行完成后，信用证状态为"交单"
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["handOverBills", "trans002","{\"BillOfLandings\":[{\"BolNO\":\"20170192102\",\"GoodsNo\":\"20170192102\",\"GoodsDesc\":\"expensive goods\",\"ShippingTime\":\"20170910\"}]}","[{\"FileName\":\"DepositDoc\",\"FileUri\":\"file_path\",\"FileHash\":\"cer324234fsfdgeergfgdfd\",\"FileSignature\":\"lixiaohuSign\",\"Uploader\":\"lixiaohu\"},{\"FileName\":\"DepositDoc\",\"FileUri\":\"file_path\",\"FileHash\":\"cer324234fsfdgeergfgdfd\",\"FileSignature\":\"lixiaohuSign\",\"Uploader\":\"lixiaohu\"}]"]}'
```
### 2.11. 通知行进行受益人交单审核。参数为3个，信用证交易号、意见、是否同意（true false）。
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["reviewBills", "trans002","no problem","true"]}'
```
### 2.12. 开证行进行承兑或者拒付操作，若承兑金额为0，说明拒付。参数为5个，信用证交易号、承兑金额、不符点、意见、是否同意（true false）。执行完成后，信用证状态为"承兑"或者回退至交单状态，若拒付时不符点不允许为空，承兑金额为0
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["lcAcceptOrReject", "trans002","123.12","","no problem","true"]}'
```
### 2.12. 若承兑，申请人需要进行付款赎单操作。参数为2个，信用证交易号、付款金额。需付款金额与未付金额相同才能判断为执行完成,目前逻辑是申请人需一次性付款，不允许多次付款，执行完成后，信用证状态为"付款赎单"
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["retireShippingBills", "trans002","123.12"]}'
```
### 2.13. 开证行可以审核申请人付款。参数为3个，信用证交易号，意见、是否同意（true false）。执行完成后，信用证状态为"付款赎单"
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["reviewRetireBills", "trans002","no problem","true"]}'
```
### 2.14. 若已经付款承兑，开证行可以进行闭卷操作。参数为2个，信用证交易号、闭卷意见。执行完成后，信用证状态为"闭卷"
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C all -n mylc -c'{"Args":["lcClose", "trans002","close lc"]}'
```
## 三、查询说明
### 3.1. 查询我们已经建立好的信用证信息
```bash
peer chaincode query -C all -n mylc -c '{"Args":["getLcByNo","136"]}'
```
### 3.2. 查询我们已经建立好的信用证信息,范围查询
```bash
peer chaincode query -C all -n mylc -c '{"Args":["getLCsByRangeNo","trans002","trans003"]}'
peer chaincode query -C all -n bcs -c '{"Args":["getBCSsByBCID","B001","Sign"]}'
peer chaincode query -C all -n bcs -c '{"Args":["getBCSList","Bank"]}'
```
### 3.3. 查询某开证行下所有的信用证,范围查询
```bash
peer chaincode query -C all -n mylc -c '{"Args":["getLcListByIssuingBank","ICBC"]}'
```
### 3.4. 查询某通知行下所有的信用证,范围查询
```bash
peer chaincode query -C all -n mylc -c '{"Args":["getLcListByAdvisingBank","ICBC"]}'
```
### 3.5. 查询某银行行下所有的信用证,范围查询
```bash
peer chaincode query -C all -n mylc -c '{"Args":["getLcListByBankId","ICBC"]}'
```
### 3.6. 查询某企业提交的下所有的信用证,范围查询
```bash
peer chaincode query -C all -n mylc -c '{"Args":["getLcListByApplicant","1"]}'
```
### 3.7. 查询受益企业相关的下所有的信用证,范围查询
```bash
peer chaincode query -C all -n mylc -c '{"Args":["getLcListByBeneficiary","1"]}'
```
### 3.8. 查询企业相关的下所有的信用证,范围查询
```bash
peer chaincode query -C all -n mylc -c '{"Args":["getLcListByCorpId","1"]}'
```
### 3.9. 查询当前用户下所有的信用证,范围查询
```bash
peer chaincode query -C all -n mylc -c '{"Args":["getLcList"]}'
```
### 3.10. 根据信用证正本拥有者查询信用证信息
```bash
peer chaincode query -C all -n mylc -c '{"Args":["getLcByOwner","1"]}'
```
### 3.11. 根据信用证申请单编号查询信用证
```bash
peer chaincode query -C all -n mylc -c '{"Args":["getLcListByApplicationFormNo","123"]}'
```
### 3.12. 根据信用证正本编号（交易号）查询信用证
```bash
peer chaincode query -C all -n mylc -c '{"Args":["getLcListByLetterOfCreditNo","123"]}'
```
### 3.13. 根据信用证正本号（非交易号）查询信用证
```bash
peer chaincode query -C all -n mylc -c '{"Args":["getLcListByLetterOfCreditLCNo","lc20171108"]}'
```