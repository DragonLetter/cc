package main


import (
	"bytes"
	"fmt"
	"time"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type LCFunctionsInterface interface {
	GetState(key string) ([]byte, error)

	//Role:申请人
	//OP:发起信用证申请
	//status:申请
	//Description：发起申请，签名。核实信息、保存申请材料。
	//Return：申请状态
	SubmitApplication(form ApplyForm,ContractId string) pb.Response

	//Role:开证行
	//OP:确认开证申请
	//Status:草稿
	//Description：银行确认开证申请，并核实信息、签名。根据申请书和贸易合同，生成信用证草稿
	//Return：信用证草稿，信用证编号
	BankConfirmApplication (form ApplyForm,ContractId string) (LCDraft,LCNumber,error)
	
	//Role：申请人
	//OP：提交保证金缴纳证明
	//Status：草稿
	//Description：申请人线下缴纳保证金，并提交保证金证明。
	//保证金证明支持电子渠道的汇款流水号或者纸质汇款回执单
	//Return：草稿
	Deposit(deposit DepositProof) （LCDraft,error)

	//Role：开证行
	//OP:信用证开立
	//Status：正本
	//Description：银行发起信用证开立，并加签，返回信用证正本
	//Return:信用证正本
	GetLCIssuce(draft LCDraft) (LCIssue,error)

	//Role：开证行、开证申请人、通知行、受益人
	//OP:参与方签名
	//Status：正本
	//Description：由多方确认信用证正本，确认并加签
	//Return：信用证正本
	LCConfirm(issue LCIssue,isLCEdit bool,isApproved bool,commit string) (LCIssue,error)

	//TODO
	//信用证撤销
	LCCancel(issue LCIssue) pb.Response

	//Role：申请人
	//OP:信用证修改
	//Status:正本修改
	//Description：申请人发起信用证修改。须由开证行、通知行、收益人重新确认
	//Return：
	LCAmend(issue LCIssue,form AmendForm)（LCIssue,error)

	//Role:物流方
	//OP:生成数字提单
	//Status:生效
	//Owner：受益人
	//Description：物流方在受托运后，向收益人开出数字提单
	GetLadingBill(issue LCIssue) (LadingBill,error)

	//Role：受益人
	//OP：交单
	//Status：交单
	//Description：收益人将信用证、提货单及其他相关单据交付开证行
	//Return：交单信息
	LCSumbitInvoce(ab LCAB) (LCAB,error)

	//Role：开证行、申请人
	//OP：开证行承兑
	//Status：承兑/拒付
	//Description：开证行承兑信用证，或者拒付信用证。
	//核实单证是否相符，如果相符银行承兑后，将提单所有权给开证申请人
	//承兑需要申请人确认
	//Return：承兑信息
	LCAcceptOrReject(ab LCAB) (LCAccept,error)

	//Role:开证行
	//OP：信用证闭卷
	//Status:闭卷
	//Description：自动审核信用证，提醒银行可以闭卷。由开证行发起信用证闭卷
	//Return：
	LCClose(issue LCIssue) pb.Response
}

