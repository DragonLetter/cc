package main

// 引入必要的依赖包
import (
	"bytes"
	"fmt"
	"time"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/pem"
	"crypto/x509"
	"strings"
	"github.com/looplab/fsm"
)

// 声明一个结构体
type SimpleChaincode struct {
	FSM *fsm.FSM
}

// 主函数，需要调用shim.Start()方法
func main() {
	t := new(SimpleChaincode)
	t.FSM = InitFSM(LCStepText[LCStart])
	err := shim.Start(t)
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// 为SimpleChaincode结构体添加Init方法，该方法实现链码初始化或升级时的处理逻辑
// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	stub.PutState("LCSequence", []byte(strconv.Itoa(1)))
	return shim.Success(nil)
}

// 为SimpleChaincode结构体添加Invoke方法，该方法实现链码运行中被调用或查询时的处理逻辑
// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	//fmt.Println("invoke is running " + function)
	// Handle different functions
	if function == "issueLc" { //create a new Lc
		//return t.issueLc(stub, args)
	} else if function == "getLcByNo" { //read a Lc
		return t.getLcByNo(stub, args)
	} else if function == "getLCsByRangeNo" { //find Lcs for owner X using rich query
		return t.getLCsByRangeNo(stub, args)
	} else if function == "getLcByOwner" { //find Lcs for owner X using rich query
		return t.getLcListByOwner(stub, args)
	} else if function == "getLcListByIssuingBank" { //find Lcs for issuing bank X using rich query
		return t.getLcListByIssuingBank(stub, args)
	} else if function == "getLcListByAdvisingBank" { //find Lcs for issuing bank X using rich query
		return t.getLcListByAdvisingBank(stub, args)
	} else if function == "getLcListByApplicant" { //find Lcs for applicant X using rich query
		return t.getLcListByApplicant(stub, args)
	} else if function == "getLcListByBeneficiary" { //find Lcs for applicant X using rich query
		return t.getLcListByBeneficiary(stub, args)
	} else if function == "getLcListByBankId" { //find Lcs for bank X using rich query
		return t.getLcListByBankId(stub, args)
	} else if function == "getLcListByCorpId" { //find Lcs for corp X using rich query
		return t.getLcListByCorpId(stub, args)
	} else if function == "getLcList" { //find Lcs for current user using rich query
		return t.getLcList(stub)
	} else if function == "getLcListByApplicationFormNo" { //find Lcs for application form No using rich query
		return t.getLcListByApplicationFormNo(stub, args)
	} else if function == "getLcListByLetterOfCreditNo" { //find Lcs for letter of credit No using rich query
		return t.getLcListByLetterOfCreditNo(stub, args)
	} else if function == "getLcListByLetterOfCreditLCNo" { //find Lcs for letter of credit LCNo using rich query
		return t.getLcListByLetterOfCreditLCNo(stub, args)
	} else if function == "saveLCApplication" {
		return t.saveLCApplication(stub, args)
	} else if function == "submitLCApplication" {
		return t.submitLCApplication(stub, args)
	} else if function == "bankConfirmApplication" {
		return t.bankConfirmApplication(stub, args)
	} else if function == "deposit" {
		return t.deposit(stub, args)
	} else if function == "issueLetterOfCredit" {
		return t.issueLetterOfCredit(stub, args)
	} else if function == "advisingBankReceiveLCNotice" {
		return t.advisingBankReceiveLCNotice(stub, args)
	} else if function == "beneficiaryReceiveLCNotice" {
		return t.beneficiaryReceiveLCNotice(stub, args)
	} else if function == "handOverBills" {
		return t.handOverBills(stub, args)
	} else if function == "appliantCheckBills" {
        return t.appliantCheckBills(stub, args)
	} else if function == "lcAmendSubmit" {
		return t.lcAmendSubmit(stub, args)
	} else if function == "lcAmendConfirm" {
		return t.lcAmendConfirm(stub, args)
	} else if function == "reviewBills" {
		return t.reviewBills(stub, args)
	} else if function == "lcAcceptOrReject" {
		return t.lcAcceptOrReject(stub, args)
	} else if function == "retireShippingBills" {
		return t.retireShippingBills(stub, args)
	} else if function == "reviewRetireBills" {
		return t.reviewRetireBills(stub, args)
	} else if function == "lcClose" {
		return t.lcClose(stub, args)
	} else if function == "saveBCSInfo" {
		return t.saveBCSInfo(stub, args)
	} else if function == "getBCSList" {
		return t.getBCSList(stub, args)
	} else if function == "getBCSByBCSNo" {
		return t.getBCSByBCSNo(stub, args)
	} else if function == "getBCSsByBCID" {
		return t.getBCSsByBCID(stub, args)
	} else if function == "issueLetterOfAmend" {
		return t.issueLetterOfAmend(stub, args)
	} else if function == "advisingLetterOfAmend" {
		return t.advisingLetterOfAmend(stub, args)
	} else if function == "beneficiaryLetterOfAmend" {
		return t.beneficiaryLetterOfAmend(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

//获得当前序号，并且将序列+1
func getNumber(stub shim.ChaincodeStubInterface) string {
	return strconv.Itoa(getNextSequence(stub, "NO"))
}

//获得LC的内部编号，并且将序列+1,并作为transactionId
func getLcNumber(stub shim.ChaincodeStubInterface) string {
	return "LC" + time.Now().Format("20060102") + strconv.Itoa(getNextSequence(stub, "LC"))
}

/****
	query function
 */
//根据信用证号查询信用证信息
func (t *SimpleChaincode) getLcByNo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	lcNum := args[0]
	lcBytes, err := stub.GetState(lcNum)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + lcNum)
	}
	return shim.Success(lcBytes)
}

//查询某个区间范围内所有的信用证信息
func (t *SimpleChaincode) getLCsByRangeNo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	startKey := args[0]
	endKey := args[1]
	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getLCsByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//查询开证行相关的信用证
func (t *SimpleChaincode) getLcListByIssuingBank(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	entityNo := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"ApplicationForm.IssuingBank.No\":\"%s\"}}", entityNo)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//查询通知行相关的信用证
func (t *SimpleChaincode) getLcListByAdvisingBank(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	entityNo := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"ApplicationForm.AdvisingBank.No\":\"%s\"}}", entityNo)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//查询银行相关的信用证
func (t *SimpleChaincode) getLcListByBankId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	entityNo := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"$or\":[{\"ApplicationForm.IssuingBank.No\":\"%s\"},{\"ApplicationForm.AdvisingBank.No\":\"%s\"}]}}", entityNo,entityNo)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}
//查询信用证申请人相关的信用证
func (t *SimpleChaincode) getLcListByApplicant(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	entityNo := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"ApplicationForm.Applicant.No\":\"%s\"}}", entityNo)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//查询信用证受益人相关的信用证
func (t *SimpleChaincode) getLcListByBeneficiary(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	entityNo := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"ApplicationForm.Beneficiary.No\":\"%s\"}}", entityNo)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}
//查询企业相关的信用证
func (t *SimpleChaincode) getLcListByCorpId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	entityNo := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"$or\":[{\"ApplicationForm.Applicant.No\":\"%s\"},{\"ApplicationForm.Beneficiary.No\":\"%s\"}]}}", entityNo,entityNo)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}
//查询当前用户相关的信用证
func (t *SimpleChaincode) getLcList(stub shim.ChaincodeStubInterface) pb.Response {

	_, _, domain := identity(stub)
	queryString := fmt.Sprintf("{\"selector\":{\"$or\":[{\"ApplicationForm.Applicant.Domain\":\"%s\"},{\"ApplicationForm.IssuingBank.Domain\":\"%s\"},{\"ApplicationForm.AdvisingBank.Domain\":\"%s\"},{\"ApplicationForm.Beneficiary.Domain\":\"%s\"}]}}", domain, domain, domain, domain)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//根据信用证申请单编号查询信用证
func (t *SimpleChaincode) getLcListByApplicationFormNo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	entityNo := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"ApplicationForm.No\":\"%s\"}}", entityNo)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//根据信用证正本编号（内部编号）查询信用证
func (t *SimpleChaincode) getLcListByLetterOfCreditNo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	entityNo := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"LetterOfCredit.No\":\"%s\"}}", entityNo)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//根据信用证正本号（非内部编号）查询信用证
func (t *SimpleChaincode) getLcListByLetterOfCreditLCNo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	entityNo := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"LetterOfCredit.LCNo\":\"%s\"}}", entityNo)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// 查看某Entity下有哪些信用证正本
func (t *SimpleChaincode) getLcListByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	entityNo := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"Owner.No\":\"%s\"}}", entityNo)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

/****
	invoke function
 */
//初始化状态机
func InitFSM(initStatus string) *fsm.FSM {
	f := fsm.NewFSM(
		initStatus,
		fsm.Events{
			{Name: "applicantSaveLCApplication", Src: []string{LCStepText[LCStart]}, Dst: LCStepText[ApplicantSaveLCApplyFormStep]},

			{Name: "applicantSubmitLCApplication", Src: []string{LCStepText[ApplicantSaveLCApplyFormStep]}, Dst: LCStepText[BankConfirmApplyFormStep]},

			{Name: "issuingBankApproveApplicantSubmitLCApplication", Src: []string{LCStepText[BankConfirmApplyFormStep]}, Dst: LCStepText[ApplicantFillLCDraftStep]},
			{Name: "issuingBankRejectApplicantSubmitLCApplication", Src: []string{LCStepText[BankConfirmApplyFormStep]}, Dst: LCStepText[ApplicantSaveLCApplyFormStep]},

			{Name: "applicantSubmitLCDraft", Src: []string{LCStepText[ApplicantFillLCDraftStep]}, Dst: LCStepText[BankIssueLCStep]},

			{Name: "issuingBankApproveIssueLC", Src: []string{LCStepText[BankIssueLCStep]}, Dst: LCStepText[AdvisingBankReceiveLCNoticeStep]},
			{Name: "issuingBankRejectIssueLC", Src: []string{LCStepText[BankIssueLCStep]}, Dst: LCStepText[ApplicantFillLCDraftStep]},

			{Name: "advisingBankApproveLCNotice", Src: []string{LCStepText[AdvisingBankReceiveLCNoticeStep]}, Dst: LCStepText[BeneficiaryReceiveLCStep]},
			{Name: "advisingBankRejectLCNotice", Src: []string{LCStepText[AdvisingBankReceiveLCNoticeStep]}, Dst: LCStepText[BankIssueLCStep]},

			{Name: "beneficiaryApproveLC", Src: []string{LCStepText[BeneficiaryReceiveLCStep]}, Dst: LCStepText[ApplicantRetireBillsStep]},
			{Name: "beneficiaryRejectLC", Src: []string{LCStepText[BeneficiaryReceiveLCStep]}, Dst: LCStepText[ApplicantLCAmendStep]},

			{Name: "applicantSubmitLCAmend", Src: []string{LCStepText[ApplicantLCAmendStep]}, Dst: LCStepText[MultiPartyCountersignStep]},

			{Name: "MultiPartyCountersignApprove", Src: []string{LCStepText[MultiPartyCountersignStep]}, Dst: LCStepText[BeneficiaryReceiveLCStep]},
			{Name: "MultiPartyCountersignReject", Src: []string{LCStepText[MultiPartyCountersignStep]}, Dst: LCStepText[ApplicantLCAmendStep]},

			// {Name: "beneficiaryHandOverBills", Src: []string{LCStepText[BeneficiaryHandOverBillsStep]}, Dst: LCStepText[AdvisingBankReviewBillsStep]},

			// {Name: "advisingBankApproveBills", Src: []string{LCStepText[AdvisingBankReviewBillsStep]}, Dst: LCStepText[IssuingBankAcceptOrRejectStep]},
			// {Name: "advisingBankRejectBills", Src: []string{LCStepText[AdvisingBankReviewBillsStep]}, Dst: LCStepText[BeneficiaryHandOverBillsStep]},

			// {Name: "issuingBankAccept", Src: []string{LCStepText[IssuingBankAcceptOrRejectStep]}, Dst: LCStepText[ApplicantRetireBillsStep]},
			// {Name: "issuingBankReject", Src: []string{LCStepText[IssuingBankAcceptOrRejectStep]}, Dst: LCStepText[AdvisingBankReviewBillsStep]},

			{Name: "applicantRetireBills", Src: []string{LCStepText[ApplicantRetireBillsStep]}, Dst: LCStepText[IssuingBankReviewRetireBillsStep]},

			{Name: "issuingBankApproveRetireBills", Src: []string{LCStepText[IssuingBankReviewRetireBillsStep]}, Dst: LCStepText[IssuingBankCloseLCStep]},
			{Name: "issuingBankRejectRetireBills", Src: []string{LCStepText[IssuingBankReviewRetireBillsStep]}, Dst: LCStepText[ApplicantRetireBillsStep]},

			{Name: "issuingBankCloseLC", Src: []string{LCStepText[IssuingBankCloseLCStep]}, Dst: LCStepText[LCEnd]},

		},
		fsm.Callbacks{
			//"enter_state": func(e *fsm.Event) { lc.enterState(e) },
		},
	)
	return f
}

/**
	Role:申请人
	OP:保存信用证申请
	status:申请
	Description:发起申请,签名.核实信息、保存申请材料。
	Return：申请状态
 */
func (t *SimpleChaincode) saveLCApplication(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	identity(stub)
	lc := &LCLetter{}
	//查询链上是否已经存在相同数据，若存在则对由原有数据进行更新，若不存在，则新增数据
	//信用证内部流水码
	no := args[0]
	//申请单
	applicationForm, err := decodeApplicationForm(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	applicationForm.IsApproved = true
	//信用证信息根据申请单得到，此时信用证为草稿状态 1:草稿 2:正本
	letterOfCredit, err := decodeLetterOfCredit(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	letterOfCredit.Status = 1
	// ==== Check if letter of credit already exists ====
	lcAsBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error(err.Error())
	}
	//判断save操作是上传附件还是保存基本信息
	var description string
	if applicationForm.Contract.FileHash == ""{
		description = "保存信用证申请单及附件"
	} else{
		description = "保存信用证申请单"
	}

	if lcAsBytes == nil {
		//此时不存在该lc的信息
		//初始化值
		acceptAmount := 0.0
		nopayAmount := applicationForm.Amount - acceptAmount
		lcStatus := Apply
		amendTimes := 0
		aBTimes := 0
		isValid := false
		isClose := false
		isCancel := false
		ti := time.Now() // 获取当前时间
		acceptDate := time.Date(1000, 01, 01, 01, 00, 00, 00, ti.Location())
		countersign := map[string]bool{
		}
		var transProgressFlow []TransProgress

		var amendFormFlow []AmendForm //发起修改

		//保证金内容在填申请时为空
		lcTransDeposit, err := decodeLCTransDeposit("{\"DepositAmount\":\"0.0\", \"CommitAmount\":\"0.0\", \"DepositDoc\":{\"FileName\":\"\",\"FileUri\":\"\",\"FileHash\":\"\",\"FileSignature\":\"\",\"Uploader\":\"\"}}")
		if err != nil {
			return shim.Error(err.Error())
		}
		//交单内容在填申请时为空
		var lcTransDocsReceive []LCTransDocsReceive
//		lCTransDocsReceive, err:= decodeLCTransDocsReceive("{\"ReceivedAmount\":\"0.0\"}")
		//货运单内容在填申请时为空
		//billOfLanding, err := decodeBillOfLanding("{\"GoodsNo\":\"\",\"GoodsDesc\":\"\",\"LoadPortName\":\"\",\"TransPortName\":\"\",\"LatestShipDate\":\"\",\"PartialShipment\":\"false\",\"TrackingNo\":\"\",\"Carrier\":{\"No\":\"\",\"Name\":\"\",\"Domain\":\"\"},\"ShippingTime\":\"\",\"Owner\":{\"NO\":\"\",\"Name\":\"\",\"Domain\":\"\"}}")
		if err != nil {
			return shim.Error(err.Error())
		}
		owner, err := decodeLegalEntity("{\"NO\":\"\",\"Name\":\"\",\"Domain\":\"\"}")
		if err != nil {
			return shim.Error(err.Error())
		}

		//申请人填写申请表时，此时还没有信用证号，执行保存操作
		ApplicantPaidAmount := 0.0
		lc = &LCLetter{no, "", applicationForm, letterOfCredit, lcTransDocsReceive, lcTransDeposit, acceptAmount, nopayAmount, acceptDate, int64(amendTimes), int64(aBTimes), false, ApplicantPaidAmount, isValid, isClose, isCancel, lcStatus, countersign, []string{}, owner, transProgressFlow, "", amendFormFlow, 0}

	} else { // 链上已经存在该LC信息
		err = json.Unmarshal(lcAsBytes, &lc) //unmarshal it aka JSON.parse()
		if err != nil {
			return shim.Error(err.Error())
		}
		lc.ApplicationForm = applicationForm
		lc.LetterOfCredit = letterOfCredit
	}
	t.FSM.SetCurrent("LCStart")

	err = t.FSM.Event("applicantSaveLCApplication") //触发状态机的事件
	if err != nil {
		return shim.Error(err.Error())
	}
	lc.CurrentStep = t.FSM.Current()
	transProgress := TransProgress{applicationForm.Applicant.Name, applicationForm.Applicant.Domain, time.Now(), description, Approve, lc.CurrentStep}
	lc.TransProgressFlow = append(lc.TransProgressFlow, transProgress)
	//=== Marshal LC ===
	lcJSONasBytes, err := json.Marshal(lc)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === Save LC to state ===
	err = stub.PutState(no, lcJSONasBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s, error: %s", args[0], err.Error()))
	}
	return shim.Success([]byte(no))
}

/**
	Role：申请人
	OP:提交信用证申请
	Status：正本
	Description：申请人提交信用证申请
	Return:申请状态
 */
func (t *SimpleChaincode) submitLCApplication(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1. no")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.ApplicationForm.Applicant.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not lc apply corporation domain:" + lc.ApplicationForm.Applicant.Domain)
	}
	t.FSM.SetCurrent(lc.CurrentStep)

	transProgress := TransProgress{userName, domain, time.Now(), "申请人提交信用证申请", Approve, lc.CurrentStep}
	lc.TransProgressFlow = append(lc.TransProgressFlow, transProgress)
	err = t.FSM.Event("applicantSubmitLCApplication") //触发状态机的事件
	if err != nil {
		return shim.Error(err.Error())
	}
	lc.LcStatus = Apply
	lc.CurrentStep = t.FSM.Current()

	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role:开证行
	OP:确认开证申请 或者拒绝开证申请
	Status:草稿
	Description：银行确认开证申请，并核实信息、签名。根据申请书和贸易合同，生成信用证草稿，银行提供保证金信息
	Return：信用证草稿
 */
func (t *SimpleChaincode) bankConfirmApplication(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5. no lcNo depositAmount opinion approveOrReject")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.ApplicationForm.IssuingBank.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not issuing bank domain:" + lc.ApplicationForm.IssuingBank.Domain)
	}
	t.FSM.SetCurrent(lc.CurrentStep)
	lcNo := args[1]
	depositAmt, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return shim.Error("2nd argument must be a numeric string")
	}
	opinionString := args[3]
	choice, err := strconv.ParseBool(args[4])
	if err != nil {
		return shim.Error("4th arguments must be bool")
	}
	var status, operation int
	if choice {
		err = t.FSM.Event("issuingBankApproveApplicantSubmitLCApplication") //触发状态机的事件
		if err != nil {
			return shim.Error(err.Error())
		}
		operation = Approve
		status = Draft
		lc.ApplicationForm.IsApproved = true
		lc.LetterOfCredit.Status = 1
		lc.LetterOfCredit.LCNo = lcNo
		lc.LCNo = lcNo
		lc.LCTransDeposit.DepositAmount = depositAmt
	} else {
		err = t.FSM.Event("issuingBankRejectApplicantSubmitLCApplication") //触发状态机的事件
		if err != nil {
			return shim.Error(err.Error())
		}
		lc.LCNo = ""
		operation = Overrule
		status = Apply
		lc.ApplicationForm.IsApproved = false
	}
	transProgress := &TransProgress{userName, domain, time.Now(), opinionString, operation, lc.CurrentStep}

	lc.TransProgressFlow = append(lc.TransProgressFlow, *transProgress)
	lc.LcStatus = status
	//lc.LCNo = lcNo
	lc.CurrentStep = t.FSM.Current()
	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role：申请人
	OP：提交保证金缴纳证明
	Status：草稿
	Description：申请人线下缴纳保证金，并提交保证金证明。
	保证金证明支持电子渠道的汇款流水号或者纸质汇款回执单
	Return：
 */
func (t *SimpleChaincode) deposit(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. no LCTransDeposit")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.ApplicationForm.Applicant.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not apply corporation domain:" + lc.ApplicationForm.IssuingBank.Domain)
	}
	t.FSM.SetCurrent(lc.CurrentStep)
	//保证金信息
	lcTransDeposit, err := decodeLCTransDeposit(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	transProgress := &TransProgress{userName, domain, time.Now(), "申请人提交保证金单据", Approve, lc.CurrentStep}
	lc.TransProgressFlow = append(lc.TransProgressFlow, *transProgress)
	lc.LCTransDeposit.CommitAmount = lcTransDeposit.CommitAmount
	lc.LCTransDeposit.DepositDoc = lcTransDeposit.DepositDoc
	//触发状态机的事件,若果已完成缴纳，则进入下一个状态
	err = t.FSM.Event("applicantSubmitLCDraft")
	if err != nil {
		return shim.Error(err.Error())
	}
	lc.CurrentStep = t.FSM.Current()
	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role：开证行
	OP:信用证开立或驳回
	Status：正本
	Description：银行发起信用证开立，同时上传信用证正本附件，返回信用证正本，同时将信用证发给通知行
	Return:信用证正本
 */
func (t *SimpleChaincode) issueLetterOfCredit(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4. no opinion approveOrReject LCOriginalAttachment")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.ApplicationForm.IssuingBank.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not issuing bank domain:" + lc.ApplicationForm.IssuingBank.Domain)
	}
	t.FSM.SetCurrent(lc.CurrentStep)

	opinionString := args[1]
	choice, err := strconv.ParseBool(args[2])
	if err != nil {
		return shim.Error("3rd arguments must be bool")
	}
	var status, operation int
	if choice {
		err = t.FSM.Event("issuingBankApproveIssueLC") //触发状态机的事件
		if err != nil {
			return shim.Error(err.Error())
		}
		operation = Approve
		status = Original
		lCOriginalAttachment, err := decodeDocument(args[3])
		if err != nil {
			return shim.Error(err.Error())
		}
		lc.LetterOfCredit.LCOriginalAttachment = lCOriginalAttachment
		lc.LetterOfCredit.Status = 2
		lc.LetterOfCredit.IssuingDate = time.Now()

	} else {
		err = t.FSM.Event("issuingBankRejectIssueLC") //触发状态机的事件
		if err != nil {
			return shim.Error(err.Error())
		}
		operation = Overrule
		status = Draft
	}
	transProgress := &TransProgress{userName, domain, time.Now(), opinionString, operation, lc.CurrentStep}
	lc.TransProgressFlow = append(lc.TransProgressFlow, *transProgress)
	lc.LcStatus = status
	lc.CurrentStep = t.FSM.Current()
	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role：通知行
	OP:通知行同意通知或者拒绝通知
	Status：正本
	Description：通知行同意通知或者拒绝通知受益人
	Return:信用证正本
 */
func (t *SimpleChaincode) advisingBankReceiveLCNotice(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3. no opinion approveOrReject")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
		if !strings.EqualFold(lc.LetterOfCredit.AdvisingBank.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not advising bank domain:" + lc.LetterOfCredit.AdvisingBank.Domain)
	}
	t.FSM.SetCurrent(lc.CurrentStep)

	opinionString := args[1]
	choice, err := strconv.ParseBool(args[2])
	if err != nil {
		return shim.Error("3rd arguments must be bool")
	}
	var operation int
	if choice {
		err = t.FSM.Event("advisingBankApproveLCNotice") //触发状态机的事件
		if err != nil {
			return shim.Error(err.Error())
		}
		operation = Approve
	} else {
		err = t.FSM.Event("advisingBankRejectLCNotice") //触发状态机的事件
		if err != nil {
			return shim.Error(err.Error())
		}
		operation = Overrule
	}
	transProgress := &TransProgress{userName, domain, time.Now(), opinionString, operation, lc.CurrentStep}
	lc.TransProgressFlow = append(lc.TransProgressFlow, *transProgress)
	lc.CurrentStep = t.FSM.Current()
	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role：受益人
	OP:受益人同意同意或者拒绝LC内容
	Status：正本
	Description：受益人同意同意或者拒绝LC内容，若同意则流转至交单步骤，若拒绝，则退回申请人修改
	Return:信用证生效 信用证正本修改
 */
func (t *SimpleChaincode) beneficiaryReceiveLCNotice(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3. no opinion approveOrReject")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.LetterOfCredit.Beneficiary.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not lc beneficiary domain:" + lc.LetterOfCredit.Beneficiary.Domain)
	}
	t.FSM.SetCurrent(lc.CurrentStep)

	opinionString := args[1]
	choice, err := strconv.ParseBool(args[2])
	if err != nil {
		return shim.Error("2nd arguments must be bool")
	}
	var operation int
	if choice {
		err = t.FSM.Event("beneficiaryApproveLC") //触发状态机的事件
		if err != nil {
			return shim.Error(err.Error())
		}
		operation = Approve
		lc.Owner = lc.LetterOfCredit.Beneficiary.LegalEntity
		lc.LcStatus = Effective
	} else {
		err = t.FSM.Event("beneficiaryRejectLC") //触发状态机的事件
		if err != nil {
			return shim.Error(err.Error())
		}
		operation = Overrule
		lc.LcStatus = Original
	}
	transProgress := &TransProgress{userName, domain, time.Now(), opinionString, operation, lc.CurrentStep}
	lc.TransProgressFlow = append(lc.TransProgressFlow, *transProgress)
	lc.CurrentStep = t.FSM.Current()
	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role:受益人
	OP:受益人执行交单操作
	Status:生效
	Owner：开证行
	Description：受益人执行交单操作，收益人将信用证、提货单及其他相关单据交付开证行
 */
func (t *SimpleChaincode) handOverBills(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3. no, BillOfLanding array, BillOfDoc array")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.LetterOfCredit.Beneficiary.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not lc beneficiary domain:" + lc.LetterOfCredit.Beneficiary.Domain)
	}

	// t.FSM.SetCurrent(lc.CurrentStep)
	// err = t.FSM.Event("beneficiaryHandOverBills") //触发状态机的事件
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// transProgress := &TransProgress{userName, domain, time.Now(), "受益人执行交单", Approve, lc.CurrentStep}
	// lc.TransProgressFlow = append(lc.TransProgressFlow, *transProgress)

	// 交单信息
	// 获取已有的交单的数量
	var billLen int = len(lc.LCTransDocsReceive)
	lCTransDocsReceive, err := decodeLCTransDocsReceive(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	BillOfLadingDocs, err := decodeDocuments(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}
	lCTransDocsReceive.No = strconv.Itoa(billLen + 1)
	lCTransDocsReceive.ReceivedDate = time.Now()
	lCTransDocsReceive.BillOfLadingDocs = BillOfLadingDocs

	// 设置受益人交单状态变化，记录在交单子结构中
	lCTransDocsReceive.HandOverBillStep = HandOverBillStep[IssuingBankCheckBillStep]
    transProgress := &TransProgress{userName, domain, time.Now(), "受益人执行交单", Approve, HandOverBillStep[BeneficiaryHandOverBillsStep]}
    lCTransDocsReceive.TransProgressFlow = append(lCTransDocsReceive.TransProgressFlow, *transProgress)
	
	lc.LCTransDocsReceive = append(lc.LCTransDocsReceive, lCTransDocsReceive)

	lc.Owner = lc.LetterOfCredit.IssuingBank.LegalEntity
	// lc.LcStatus = HandOverBill
	// lc.CurrentStep = t.FSM.Current()

	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role:申请人
	OP:申请人执行审单操作
	Status:生效
	Owner：申请人
	Description：申请人执行审核受益人交单操作，根据受益人提交的单据决定是否同意开证行进行承兑
 */
func (t *SimpleChaincode) appliantCheckBills(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 3. no, BillOfLanding array, BillOfDoc array")
	}
	lcNo := args[0]
	lcBytes, err := stub.GetState(lcNo)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + lcNo)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.LetterOfCredit.Applicant.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not lc Applicant domain:" + lc.LetterOfCredit.Applicant.Domain)
	}

	// 申请人审单结果
	billNo := args[1]
	opinionString := args[2]
	choice, err := strconv.ParseBool(args[3])
	if err != nil {
		return shim.Error("2nd arguments must be bool")
	}
	var operation int
	var handleStep string
	if choice {
		operation = Approve
		handleStep = HandOverBillStep[IssuingBankAcceptanceStep]
		lc.Owner = lc.LetterOfCredit.IssuingBank.LegalEntity
	} else {
		operation = Overrule
		handleStep = HandOverBillStep[ApplicantRejectStep]
		lc.Owner = lc.LetterOfCredit.Beneficiary.LegalEntity
	}

	// 设置交单状态变化，记录在交单子结构中
	for i := 0; i < len(lc.LCTransDocsReceive); i++ {
		if lc.LCTransDocsReceive[i].No == billNo {
			lc.LCTransDocsReceive[i].HandOverBillStep = handleStep
			lc.LCTransDocsReceive[i].Discrepancy = opinionString
			transProgress := &TransProgress{userName, domain, time.Now(), opinionString, operation, HandOverBillStep[ApplicantAcceptOrRejectStep]}
			lc.LCTransDocsReceive[i].TransProgressFlow = append(lc.LCTransDocsReceive[i].TransProgressFlow, *transProgress)
            break
		}
	}

	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(lcNo, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role：申请人
	OP:信用证修改
	Status：
	Description：申请人发起信用证修改。须由开证行、通知行、受益人重新确认
	Return:
 */

 func (t *SimpleChaincode) lcAmendSubmit(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. no")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.ApplicationForm.Applicant.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not lc apply corporation domain:" + lc.ApplicationForm.Applicant.Domain)
	}
	// t.FSM.SetCurrent(lc.CurrentStep)

	amendFormData, err := decodeAmendForm(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	// expireDate,_:= time.Parse("2018-08-07T04:56:07.000+00:00", amendFormData.AmendExpiryDate)
	lc.AmendNum = lc.AmendNum + 1

	var amendFormTransFlow []AmendFormProgress

	amendFormTransProgress := AmendFormProgress{userName, domain, time.Now(), "申请人提交发起修改申请", Approve, AmendStepText[AmendApplicantSubmitStep]}
	amendForm := AmendForm{lc.AmendNum, amendFormData.AmendTimes, amendFormData.AmendedCurrency, amendFormData.AmendedAmt, amendFormData.AddedDays, amendFormData.AmendExpiryDate, amendFormData.TransPortName, amendFormData.AddedDepositAmt, AmendStepText[AmendIssuingBankAcceptStep], amendFormTransFlow, time.Now()}
	amendForm.AmendFormProgressFlow = append(amendForm.AmendFormProgressFlow, amendFormTransProgress)

	lc.AmendFormFlow = append(lc.AmendFormFlow, amendForm)

	// err = t.FSM.Event("applicantSubmitLCApplication") //触发状态机的事件
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// lc.LcStatus = Apply
	// lc.CurrentStep = t.FSM.Current()

	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role：开证行
	OP:发起修改同意或拒绝
	Status：发起修改
	Description：申请人提交发起修改，开证行操作（第一步）
	Return:
 */
 func (t *SimpleChaincode) issueLetterOfAmend(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4. no opinion approveOrReject issueLetterOfAmend")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.ApplicationForm.IssuingBank.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not issuing bank domain:" + lc.ApplicationForm.IssuingBank.Domain)
	}
	// t.FSM.SetCurrent(lc.

	amendNo := args[1]
	
	// queryString := fmt.Sprintf("{\"selector\":{\"$and\":[{\"no\":\"%s\"},{\"AmendForm.AmendNo\":\"%s\"}]}}", no, amendNo)
	// queryResults, err := getQueryResultForQueryString(stub, queryString)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	// amend := AmendForm{}
	// err = json.Unmarshal(lcBytes, &amend) //unmarshal it aka JSON.parse()
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	opinionString := args[2]
	choice, err := strconv.ParseBool(args[3])
	if err != nil {
		return shim.Error("3rd arguments must be bool")
	}	
	
    for i := 0; i < len(lc.AmendFormFlow); i++{
		if amendNo == strconv.Itoa(lc.AmendFormFlow[i].AmendNo){
			if choice {
				lc.AmendFormFlow[i].Status = AmendStepText[AmendAdvisingBankAcceptStep]	
				amendFormTransProgress := AmendFormProgress{userName, domain, time.Now(), opinionString, Approve, AmendStepText[AmendIssuingBankAcceptStep]}
				lc.AmendFormFlow[i].AmendFormProgressFlow = append(lc.AmendFormFlow[i].AmendFormProgressFlow, amendFormTransProgress)		
			} else {
				lc.AmendFormFlow[i].Status = AmendStepText[AmendEnd]
				amendFormTransProgress := AmendFormProgress{userName, domain, time.Now(), opinionString, Overrule, AmendStepText[AmendIssuingBankRejectStep]}
				lc.AmendFormFlow[i].AmendFormProgressFlow = append(lc.AmendFormFlow[i].AmendFormProgressFlow, amendFormTransProgress)
			}
			break;
		}
	}

	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role：通知行
	OP:发起修改同意或拒绝
	Status：发起修改
	Description：申请人提交发起修改，通知行操作（开证行同意后到通知行）
	Return:
 */
 func (t *SimpleChaincode) advisingLetterOfAmend(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4. no opinion approveOrReject advisingLetterOfAmend")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.LetterOfCredit.AdvisingBank.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not advising bank domain:" + lc.LetterOfCredit.AdvisingBank.Domain)
	}
	// t.FSM.SetCurrent(lc.

	amendNo := args[1]
	
	// queryString := fmt.Sprintf("{\"selector\":{\"$and\":[{\"no\":\"%s\"},{\"AmendForm.AmendNo\":\"%s\"}]}}", no, amendNo)
	// queryResults, err := getQueryResultForQueryString(stub, queryString)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	// amend := AmendForm{}
	// err = json.Unmarshal(lcBytes, &amend) //unmarshal it aka JSON.parse()
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	opinionString := args[2]
	choice, err := strconv.ParseBool(args[3])
	if err != nil {
		return shim.Error("3rd arguments must be bool")
	}	
	
    for i := 0; i < len(lc.AmendFormFlow); i++{
		if amendNo == strconv.Itoa(lc.AmendFormFlow[i].AmendNo){
			if choice {
				lc.AmendFormFlow[i].Status = AmendStepText[AmendBeneficiaryAcceptStep]	
				amendFormTransProgress := AmendFormProgress{userName, domain, time.Now(), opinionString, Approve, AmendStepText[AmendAdvisingBankAcceptStep]}
				lc.AmendFormFlow[i].AmendFormProgressFlow = append(lc.AmendFormFlow[i].AmendFormProgressFlow, amendFormTransProgress)		
			} else {
				lc.AmendFormFlow[i].Status = AmendStepText[AmendEnd]
				amendFormTransProgress := AmendFormProgress{userName, domain, time.Now(), opinionString, Approve, AmendStepText[AmendAdvisingBankRejectStep]}
				lc.AmendFormFlow[i].AmendFormProgressFlow = append(lc.AmendFormFlow[i].AmendFormProgressFlow, amendFormTransProgress)
			}
			break;
		}
	}

	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}


/**
	Role：受益人
	OP:发起修改同意或拒绝
	Status：发起修改
	Description：申请人提交发起修改，通知行操作（开证行同意后到通知行）
	Return:
 */
 func (t *SimpleChaincode) beneficiaryLetterOfAmend(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4. no opinion approveOrReject beneficiaryLetterOfAmend")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.LetterOfCredit.Beneficiary.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not lc beneficiary domain:" + lc.LetterOfCredit.Beneficiary.Domain)
	}
	// t.FSM.SetCurrent(lc.

	amendNo := args[1]
	
	// queryString := fmt.Sprintf("{\"selector\":{\"$and\":[{\"no\":\"%s\"},{\"AmendForm.AmendNo\":\"%s\"}]}}", no, amendNo)
	// queryResults, err := getQueryResultForQueryString(stub, queryString)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	// amend := AmendForm{}
	// err = json.Unmarshal(lcBytes, &amend) //unmarshal it aka JSON.parse()
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	opinionString := args[2]
	choice, err := strconv.ParseBool(args[3])
	if err != nil {
		return shim.Error("3rd arguments must be bool")
	}	
	
    for i := 0; i < len(lc.AmendFormFlow); i++{
		if amendNo == strconv.Itoa(lc.AmendFormFlow[i].AmendNo){
			if choice {
				lc.AmendFormFlow[i].Status = AmendStepText[AmendEnd]	
				amendFormTransProgress := AmendFormProgress{userName, domain, time.Now(), opinionString, Approve, AmendStepText[AmendBeneficiaryAcceptStep]}
				lc.AmendFormFlow[i].AmendFormProgressFlow = append(lc.AmendFormFlow[i].AmendFormProgressFlow, amendFormTransProgress)
				
				//更新正本信息
				// expireDate,_:= time.Parse("2018-08-07T04:56:07.000+00:00", lc.AmendFormFlow[i].AmendExpiryDate)				
				lc.AmendTimes = lc.AmendFormFlow[i].AmendTimes

				lc.LetterOfCredit.Currency = lc.AmendFormFlow[i].AmendedCurrency
				lc.ApplicationForm.Currency = lc.AmendFormFlow[i].AmendedCurrency

				lc.LetterOfCredit.Amount = lc.AmendFormFlow[i].AmendedAmt
				lc.ApplicationForm.Amount = lc.AmendFormFlow[i].AmendedAmt

				lc.LetterOfCredit.DocDelay = lc.LetterOfCredit.DocDelay + lc.AmendFormFlow[i].AddedDays
				lc.ApplicationForm.DocDelay = lc.ApplicationForm.DocDelay + lc.AmendFormFlow[i].AddedDays

				lc.LetterOfCredit.ExpiryDate = lc.AmendFormFlow[i].AmendExpiryDate
				lc.ApplicationForm.ExpiryDate = lc.AmendFormFlow[i].AmendExpiryDate

				lc.LetterOfCredit.GoodsInfo.ShippingPlace = lc.AmendFormFlow[i].TransPortName
				lc.ApplicationForm.GoodsInfo.ShippingPlace = lc.AmendFormFlow[i].TransPortName

			
				//   lc.LCTransDeposit.DepositAmount = lc.LCTransDeposit.DepositAmount + lc.AmendFormFlow[i].AddedDepositAmt		
				lc.LetterOfCredit.EnsureAmount = lc.LetterOfCredit.EnsureAmount + lc.AmendFormFlow[i].AddedDepositAmt
				lc.ApplicationForm.EnsureAmount = lc.ApplicationForm.EnsureAmount + lc.AmendFormFlow[i].AddedDepositAmt
			} else {
				lc.AmendFormFlow[i].Status = AmendStepText[AmendEnd]
				amendFormTransProgress := AmendFormProgress{userName, domain, time.Now(), opinionString, Approve, AmendStepText[AmendBeneficiaryRejectStep]}
				lc.AmendFormFlow[i].AmendFormProgressFlow = append(lc.AmendFormFlow[i].AmendFormProgressFlow, amendFormTransProgress)
			}
			break;
		}
	}

	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role：开证行、通知行
	OP:参与方签名
	Status：正本
	Description：由多方确认信用证正本，确认并加签,多方会签，若均签名同意，则发给受益人。这么做的原因是因为受益人在并行流程中选择同意后，并填写发货信息后，如果开证行或者通知行不同意，则流程仍然回到申请人修改，流程效率较为低下
	Return：信用证正本
 */
func (t *SimpleChaincode) lcAmendConfirm(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3. no opinion approveOrReject")
	}
	_, userName, domain := identity(stub)
	if strings.EqualFold("", userName) {
		return shim.Error("Error stub.GetCreator")
	}

	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	if !(strings.EqualFold(lc.LetterOfCredit.IssuingBank.Domain, domain) || strings.EqualFold(lc.LetterOfCredit.AdvisingBank.Domain, domain)) {
		return shim.Error("Current operator domain:" + domain + " is not issuing bank or advising bank:" + lc.LetterOfCredit.IssuingBank.Domain + "or" + lc.LetterOfCredit.AdvisingBank.Domain)
	}
	t.FSM.SetCurrent(lc.CurrentStep)

	//上一个状态必须是正本修改状态，才能执行会签操作
	// if !(lc.LcStatus == OriginalModify) {
	// 	return shim.Error(no + "LC's Last Status is not originalModify state, can not counter sign")
	// }
	opinionString := args[1]
	choice, err := strconv.ParseBool(args[2])
	if err != nil {
		return shim.Error("2nd arguments must be bool")
	}
	var operation int
	if choice {
		operation = Approve
		transProgress := &TransProgress{userName, domain, time.Now(), opinionString, operation, lc.CurrentStep}
		counterSign := lc.Countersign
		counterSign[userName] = choice
		lc.Countersign = counterSign
		lc.TransProgressFlow = append(lc.TransProgressFlow, *transProgress)

		//如果counterSign的数量等于2，说明之前的审批人也同意，此时开证行、通知行均同意，同时步骤到受益人收取信用证同时发货步骤
		if len(counterSign) == 2 {
			err = t.FSM.Event("MultiPartyCountersignApprove") //触发状态机的事件
			if err != nil {
				return shim.Error(err.Error())
			}
			// lc.LcStatus = Original
			lc.CurrentStep = t.FSM.Current()

			jsonB, _ := json.Marshal(lc)
			err = stub.PutState(no, jsonB) //rewrite the lc
			if err != nil {
				return shim.Error(err.Error())
			}
			return shim.Success(nil)
		}
		jsonB, _ := json.Marshal(lc)
		err = stub.PutState(no, jsonB) //rewrite the lc
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte(userName + "'s opinion is Yes, But CounterSign is not finished!Wait another approve"))

	} else {
		//如果开证行、通知行有一家不同意，回到上一步
		err = t.FSM.Event("MultiPartyCountersignReject") //触发状态机的事件
		if err != nil {
			return shim.Error(err.Error())
		}
		operation = Overrule
		transProgress := &TransProgress{userName, domain, time.Now(), opinionString, operation, lc.CurrentStep}
		lc.TransProgressFlow = append(lc.TransProgressFlow, *transProgress)
		// lc.LcStatus = OriginalModify
		lc.CurrentStep = t.FSM.Current()

		jsonB, _ := json.Marshal(lc)
		err = stub.PutState(no, jsonB) //rewrite the lc
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte(userName + "'s opinion is No, reject lc to amend step!"))

	}
}

/** TODO
	Role：
	OP:信用证撤销
	Status:
	Description：
	Return：
 */
func (t *SimpleChaincode) lcCancel(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	return shim.Success(nil)
}

/**
	Role:开证行
	OP:确认受益人提交的单
	Status:生效
	Description：开证行审核受益人提交的单据
 */
func (t *SimpleChaincode) reviewBills(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 3. no opinion approveOrReject")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.LetterOfCredit.IssuingBank.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not lc IssuingBank domain:" + lc.LetterOfCredit.IssuingBank.Domain)
	}
	billNo := args[1]
	opinionString := args[2]
	choice, err := strconv.ParseBool(args[3])
	if err != nil {
		return shim.Error("3nd arguments must be bool")
	}
	var operation int
	var handleStep string
	var curStep string
	if choice {
		if err != nil {
			return shim.Error(err.Error())
		}
		operation = Approve
		handleStep = HandOverBillStep[IssuingBankCheckBillStep]
		curStep = HandOverBillStep[ApplicantAcceptOrRejectStep]
		lc.Owner = lc.LetterOfCredit.IssuingBank.LegalEntity
	} else {
		if err != nil {
			return shim.Error(err.Error())
		}
		operation = Overrule
		handleStep = HandOverBillStep[IssuingBankRejectStep]
		lc.Owner = lc.LetterOfCredit.IssuingBank.LegalEntity
	}
	// 设置交单状态变化，记录在交单子结构中
	for i := 0; i < len(lc.LCTransDocsReceive); i++ {
		if lc.LCTransDocsReceive[i].No == billNo {
			lc.LCTransDocsReceive[i].Discrepancy = opinionString
        	lc.LCTransDocsReceive[i].HandOverBillStep = curStep
	        transProgress := &TransProgress{userName, domain, time.Now(), opinionString, operation, handleStep}
    	    lc.LCTransDocsReceive[i].TransProgressFlow = append(lc.LCTransDocsReceive[i].TransProgressFlow, *transProgress)
        	break
    	}
	}

	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role：开证行
	OP：开证行承兑
	Status：承兑/拒付
	Description：开证行承兑信用证，或者拒付信用证。拒付时需要提出不符点，然后由
	核实单证是否相符，如果相符银行承兑后，将提单所有权给开证申请人
	Return：承兑信息
 */
func (t *SimpleChaincode) lcAcceptOrReject(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5. no acceptAmount discrepancies opinion approveOrReject")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.LetterOfCredit.IssuingBank.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not issuing bank domain:" + lc.LetterOfCredit.IssuingBank.Domain)
	}
	billNo := args[1]
	acceptAmount, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return shim.Error("2st argument must be a numeric string")
	}
	opinionString := args[3]
	choice, err := strconv.ParseBool(args[4])
	if err != nil {
		return shim.Error("3nd arguments must be bool")
	}
	var operation int
	var handleStep string
	var curStep string
	if choice {
		lc.AcceptAmount = lc.AcceptAmount + acceptAmount
		lc.NotPayAmount = lc.LetterOfCredit.Amount - lc.AcceptAmount
		if err != nil {
			return shim.Error(err.Error())
		}
		operation = Approve
		handleStep = HandOverBillStep[IssuingBankAcceptanceStep]
		curStep = HandOverBillStep[HandoverBillSuccStep]
		lc.Owner = lc.LetterOfCredit.IssuingBank.LegalEntity
	} else {
		if err != nil {
			return shim.Error(err.Error())
		}
		operation = Overrule
		handleStep = HandOverBillStep[IssuingBankRejectStep]
		curStep = HandOverBillStep[IssuingBankRejectStep]
		lc.Owner = lc.LetterOfCredit.IssuingBank.LegalEntity
	}
	// 设置交单状态变化，记录在交单子结构中
	for i := 0; i < len(lc.LCTransDocsReceive); i++ {
		if lc.LCTransDocsReceive[i].No == billNo {
        	lc.LCTransDocsReceive[i].HandOverBillStep = curStep
	        transProgress := &TransProgress{userName, domain, time.Now(), opinionString, operation, handleStep}
			lc.LCTransDocsReceive[i].TransProgressFlow = append(lc.LCTransDocsReceive[i].TransProgressFlow, *transProgress)
        	break
    	}
	}
	
	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role：申请人
	OP：付款赎单
	Status：付款赎单
	Description：申请人进行付款赎单操作，赎回货运单。
 */
func (t *SimpleChaincode) retireShippingBills(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. no, payment")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.LetterOfCredit.Applicant.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not applicant domain:" + lc.LetterOfCredit.Applicant.Domain)
	}
	t.FSM.SetCurrent(lc.CurrentStep)

	if lc.LcStatus != Effective {
		return shim.Error("lc status must be Effective. LCNumber:" + no)
	}
	payment, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error("1st argument must be a numeric string")
	}

	//是否已支付为是，状态为付款赎单状态，货运单拥有者为信用证申请企业
	lc.IsApplicantPaid = true
	lc.ApplicantPaidAmount = payment
	lc.LcStatus = RetireBills
	transProgress := &TransProgress{userName, domain, time.Now(), "", Approve, lc.CurrentStep}
	lc.TransProgressFlow = append(lc.TransProgressFlow, *transProgress)
	err = t.FSM.Event("applicantRetireBills") //触发状态机的事件，付款赎单
	if err != nil {
		return shim.Error(err.Error())
	}
	lc.CurrentStep = t.FSM.Current()

	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role:开证行
	OP:开证行对申请人的付款申请进行审核受益人提交的单
	Status:付款赎单
	Description：开证行审核申请人的付款
 */
func (t *SimpleChaincode) reviewRetireBills(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3. no opinion approveOrReject")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.LetterOfCredit.IssuingBank.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not lc issuing bank domain:" + lc.LetterOfCredit.IssuingBank.Domain)
	}
	t.FSM.SetCurrent(lc.CurrentStep)

	opinionString := args[1]
	choice, err := strconv.ParseBool(args[2])
	if err != nil {
		return shim.Error("2nd arguments must be bool")
	}
	var operation int
	if choice {
		err = t.FSM.Event("issuingBankApproveRetireBills") //触发状态机的事件
		if err != nil {
			return shim.Error(err.Error())
		}
		operation = Approve
	} else {
		err = t.FSM.Event("issuingBankRejectRetireBills") //触发状态机的事件
		if err != nil {
			return shim.Error(err.Error())
		}
		operation = Overrule
	}
	transProgress := &TransProgress{userName, domain, time.Now(), opinionString, operation, lc.CurrentStep}
	lc.TransProgressFlow = append(lc.TransProgressFlow, *transProgress)
	lc.CurrentStep = t.FSM.Current()
	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
	Role:开证行
	OP：信用证闭卷
	Status:闭卷
	Description：自动审核信用证，提醒银行可以闭卷。由开证行发起信用证闭卷
	Return：
 */
func (t *SimpleChaincode) lcClose(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. no, comments")
	}
	no := args[0]
	lcBytes, err := stub.GetState(no)
	if err != nil {
		return shim.Error("query Letter of Credit fail. Number:" + no)
	}
	lc := LCLetter{}
	err = json.Unmarshal(lcBytes, &lc) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	_, userName, domain := identity(stub)
	if !strings.EqualFold(lc.LetterOfCredit.IssuingBank.Domain, domain) {
		return shim.Error("Current operator domain:" + domain + " is not issuing bank domain:" + lc.LetterOfCredit.IssuingBank.Domain)
	}
	t.FSM.SetCurrent(lc.CurrentStep)

	if lc.LcStatus != RetireBills {
		return shim.Error("lc status must be RetirementOfDocuments. LCNumber:" + no)
	}
	lc.LcStatus = Close
	lc.IsClose = true
	transProgress := &TransProgress{userName, domain, time.Now(), args[1], Approve, lc.CurrentStep}
	lc.TransProgressFlow = append(lc.TransProgressFlow, *transProgress)
	err = t.FSM.Event("issuingBankCloseLC") //触发状态机的事件，付款赎单
	if err != nil {
		return shim.Error(err.Error())
	}
	lc.CurrentStep = t.FSM.Current()

	jsonB, _ := json.Marshal(lc)
	err = stub.PutState(no, jsonB) //rewrite the lc
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)

}

//获取当前用户身份
func identity(stub shim.ChaincodeStubInterface) (pb.Response, string, string) {
	creatorByte, err := stub.GetCreator()
	if err != nil {
		return shim.Error("Error stub.GetCreator"), "", ""
	}
	fmt.Println(string(creatorByte))
	certStart := bytes.IndexAny(creatorByte, "-----") // Devin:I don't know why sometimes -----BEGIN is invalid, so I use -----
	if certStart == -1 {
		return shim.Error("No certificate found"), "", ""
	}
	certText := creatorByte[certStart:]
	fmt.Println("certStart:" + strconv.Itoa(certStart))
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return shim.Error("Could not decode the PEM structure"), "", ""
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return shim.Error("ParseCertificate failed"), "", ""
	}
	//fmt.Println(cert)
	//fmt.Println(cert.Subject.CommonName)
	//fmt.Println(cert.Subject)
	//domain :=strings.Join(cert.PermittedDNSDomains,"-")

	if len(cert.Subject.Organization) > 0 {
		return shim.Success(nil), cert.Subject.OrganizationalUnit[0], cert.Subject.Organization[0]
	}
	stringSlice := strings.Split(cert.Subject.CommonName, "@")
	return shim.Success(nil), stringSlice[0], stringSlice[1]
}

/**
	Role:管理员
	OP：添加银行、企业、签约信息
	Description：添加银行、企业、签约信息到链上
	Return：
 */
func (t *SimpleChaincode) saveBCSInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	identity(stub)
	no := args[0]
	bcs, err := decodeBCSData(args[1]);
	if err != nil {
		return shim.Error(err.Error())
	}
	//=== Marshal Data of BCS ===
	bcsJSONasBytes, err := json.Marshal(bcs)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === Save BCS to state ===
	err = stub.PutState(no, bcsJSONasBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s, error: %s", args[0], err.Error()))
	}
	return shim.Success([]byte(no))
}

/**
	Role:管理员
	OP：获取银行、企业、签约信息
	Description：通过类型从链上获取银行、企业、签约信息
	Return：
 */
func (t *SimpleChaincode) getBCSList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	dtype := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"%s\"}}", dtype)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

/**
	Role:管理员
	OP：获取银行、企业、签约信息
	Description：根据id从链上获取银行、企业、签约信息
	Return：
 */
func (t *SimpleChaincode) getBCSByBCSNo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	no := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"No\":\"%s\"}}", no)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

/**
	Role:管理员
	OP：获取银行、企业、签约信息
	Description：根据银行/企业id从链上获取银行、企业、签约信息
	Return：
 */
func (t *SimpleChaincode) getBCSsByBCID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	no := args[0]
	dtype := args[1]
	queryString := fmt.Sprintf("{\"selector\":{\"$or\":[{\"DataBank.No\":\"%s\"},{\"DataCorp.No\":\"%s\"}],\"and\":[{\"Type\":\"%s\"}]}}", no,no,dtype);
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}
