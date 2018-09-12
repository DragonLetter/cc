package main

import (
	"time"
)

//角色
const(
	Applicant = iota //0
	IssuingBank //1
	AdvisingBank //2
	Beneficiary //3
)
//信用证状态
const (
	Apply = iota
	Draft
	Original
	OriginalModify
	Effective
	HandOverBill
	Accept
	RetireBills
	Reject
	Close
)
//步骤
const (
	LCStart = iota
	ApplicantSaveLCApplyFormStep
	//ApplicantSubmitLCApplyFormStep
	BankConfirmApplyFormStep
	ApplicantFillLCDraftStep
	BankIssueLCStep
	AdvisingBankReceiveLCNoticeStep
	BeneficiaryReceiveLCStep
	ApplicantLCAmendStep
	MultiPartyCountersignStep
	// BeneficiaryHandOverBillsStep
	// AdvisingBankReviewBillsStep
	// IssuingBankAcceptOrRejectStep
	ApplicantRetireBillsStep
	IssuingBankReviewRetireBillsStep
	IssuingBankCloseLCStep
	LCEnd
)

var LCStepText  = map[int] string{
	LCStart : "LCStart",
	ApplicantSaveLCApplyFormStep : "ApplicantSaveLCApplyFormStep",
	//ApplicantSubmitLCApplyFormStep : "ApplicantSubmitLCApplyFormStep",
	BankConfirmApplyFormStep : "BankConfirmApplyFormStep",
	ApplicantFillLCDraftStep : "ApplicantFillLCDraftStep",
	BankIssueLCStep : "BankIssueLCStep",
	AdvisingBankReceiveLCNoticeStep : "AdvisingBankReceiveLCNoticeStep",
	BeneficiaryReceiveLCStep : "BeneficiaryReceiveLCStep",
	ApplicantLCAmendStep : "ApplicantLCAmendStep",
	MultiPartyCountersignStep : "MultiPartyCountersignStep",
	// BeneficiaryHandOverBillsStep : "BeneficiaryHandOverBillsStep",
	// AdvisingBankReviewBillsStep : "AdvisingBankReviewBillsStep",
	// IssuingBankAcceptOrRejectStep : "IssuingBankAcceptOrRejectStep",
	ApplicantRetireBillsStep : "ApplicantRetireBillsStep",
	IssuingBankReviewRetireBillsStep : "IssuingBankReviewRetireBillsStep",
	IssuingBankCloseLCStep : "IssuingBankCloseLCStep",
	LCEnd : "LCEnd",
}

var HandOverBillStep = map[int] string{
	BeneficiaryHandOverBillsStep : "BeneficiaryHandOverBillsStep",    // 受益人交单，初始状态
	IssuingBankCheckBillStep : "IssuingBankCheckBillStep",    // 开证行审单
	ApplicantAcceptOrRejectStep : "ApplicantAcceptOrRejectStep",    // 申请人接受或拒绝审单结果
	IssuingBankAcceptanceStep : "IssuingBankAcceptanceStep",    // 开证行承兑
	ApplicantRejectStep : "ApplicantRejectStep",    // 申请人拒付，结束状态
	IssuingBankRejectStep : "IssuingBankRejectStep",    // 开证行拒付，结束状态
	HandoverBillSuccStep : "HandoverBillSuccStep",    // 交单成功，结束状态
}

//操作状态
const (
	Approve = iota //同意
	Overrule  //驳回
	Processing //处理中
)

//链上存储的信用证信息
type LCLetter struct {
	No string `json:"no"`
	//L/C NO 信用证号(20字段)
	LCNo string `json:"lcNo"`
	//信用证申请单
	ApplicationForm ApplicationForm
	//信用证正本
	LetterOfCredit LetterOfCredit
	//交单信息
	LCTransDocsReceive []LCTransDocsReceive
	//保证金
	LCTransDeposit LCTransDeposit
	//承兑金额
	AcceptAmount float64
	//未付金额
	NotPayAmount float64
	//承兑日期
	AcceptDate time.Time `json:"acceptDate,string,omitempty"`
	//改证次数
	AmendTimes int64
	//到单次数
	ABTimes int64
	//申请人是否已经付款
	IsApplicantPaid bool `json:"isApplyCorpPaid,string,omitempty"`
	//是否有效
	IsValid bool
	//是否闭卷
	IsClose bool
	//是否撤销
	IsCancel bool
	//信用证状态
	LcStatus int `json:"lcStatus"`
	//会签单
	Countersign map[string]bool
	//不符点
	Discrepancy []string
	//当前谁拥有这个信用证
	Owner LegalEntity
	//审批记录
	TransProgressFlow []TransProgress
	//信用证当前步骤
	CurrentStep string
}

//信用证申请表
type ApplicationForm struct {
	No string
	//开证企业
	Applicant Corporation
	//受益企业
	Beneficiary Corporation
	//开证行
	IssuingBank Bank
	//通知行
	AdvisingBank Bank
	//到期日
	ExpiryDate time.Time `json:"expiryDate,string,omitempty"`
	//到期地点
	ExpiryPlace string
	//是否即期
	IsAtSight bool `json:"isAtSight,string,omitempty"`
	//远期付款期限
	AfterSight int `json:"afterSight,string,omitempty"`
	//货物信息
	GoodsInfo GoodsInfo
	//单据需求
	DocumentRequire int64 `json:"documentRequire,string,omitempty"`
	//结算货币
	Currency string
	//金额
	Amount float64 `json:"amount,string,omitempty"`
	//保证金金额
	EnsureAmount float64 `json:"EnsureAmount,string,omitempty"`
	//是否可议付
	Negotiate int `json:"Negotiate,string,omitempty"`
	//是否可转让
	Transfer int `json:"Transfer,string,omitempty"`
	//是否可保兑
	Confirmed int `json:"Confirmed,string,omitempty"`
	//短装
	Lowfill float64 `json:"Lowfill,string,omitempty"`
	//溢装
	Overfill float64 `json:"Overfill,string,omitempty"`
	//申请时间
	ApplyTime time.Time `json:"applyTime,string,omitempty"`
	//在开证行产生的费用由谁承担：1买方、2卖方
	ChargeInIssueBank int64 `json:"chargeInIssueBank,string,omitempty"`
	//在开证行外产生的费用由谁承担：1买方、2卖方
	ChargeOutIssueBank int64 `json:"chargeOutIssueBank,string,omitempty"`
	//单据最晚在签发多少日后提交
	DocDelay int64 `json:"docDelay,string,omitempty"`
	//其他需求
	OtherRequire string
	//贸易合同
	Contract Contract
	//附件
	Attachments []Document
	//是否通过
	IsApproved bool `json:"isApproved,string,omitempty"`
	//IssueBank Bank
}

//信用证文本
type LetterOfCredit struct {
	//编号
	No string
	//信用证编号
	LCNo string
	//开证企业
	Applicant Corporation
	//受益企业
	Beneficiary Corporation
	//开证行
	IssuingBank Bank
	//通知行
	AdvisingBank Bank
	//开立日期
	IssuingDate time.Time `json:"issuingDate,string,omitempty"`
	//到期日
	ExpiryDate time.Time `json:"expiryDate,string,omitempty"`
	//到期地点
	ExpiryPlace string
	//是否即期
	IsAtSight bool `json:"isAtSight,string,omitempty"`
	//远期付款期限
	AfterSight int `json:"afterSight,string,omitempty"`
	//货物信息
	GoodsInfo GoodsInfo
	//单据需求
	DocumentRequire string
	//结算货币
	Currency string
	//金额
	Amount float64 `json:"amount,string,omitempty"`
	//保证金金额
	EnsureAmount float64 `json:"EnsureAmount,string,omitempty"`
	//是否可议付
	Negotiate int `json:"Negotiate,string,omitempty"`
	//是否可转让
	Transfer int `json:"Transfer,string,omitempty"`
	//是否可保兑
	Confirmed int `json:"Confirmed,string,omitempty"`
	//短装
	Lowfill float64 `json:"Lowfill,string,omitempty"`
	//溢装
	Overfill float64 `json:"Overfill,string,omitempty"`
	//申请时间
	ApplyTime time.Time `json:"applyTime,string,omitempty"`
	//在开证行产生的费用由谁承担：1买方、2卖方
	ChargeInIssueBank int64 `json:"chargeInIssueBank,string,omitempty"`
	//在开证行外产生的费用由谁承担：1买方、2卖方
	ChargeOutIssueBank int64 `json:"chargeOutIssueBank,string,omitempty"`
	//单据最晚在签发多少日后提交
	DocDelay int64 `json:"docDelay,string,omitempty"`
	//其他需求
	OtherRequire string
	//贸易合同
	Contract Contract
	//附件
	Attachments []Document
	//状态：1.草稿 2.正本
	Status int64 `json:"status,string,omitempty"`
	//正本附件
	LCOriginalAttachment Document
}

//保证金要求
type LCTransDeposit struct {
	//保证金总额
	DepositAmount float64 `json:"depositAmount,string,omitempty"`
	//已交金额
	CommitAmount float64 `json:"commitAmount,string,omitempty"`
	//保证金单据
	DepositDoc Document
}

//交易
type Transaction struct {
	//交易号
	TransactionId string
	//申请企业
	Applicant string
	//受益企业
	Beneficiary string
	//开证行
	IssuingBank string
	//通知行
	AdvisingBank string
	//金额
	Amount string
	//币种
	Currency string
	//状态
	Status string
	//操作日期
	IssueDate string
}

//交易进度
type TransProgress struct {
	Name string
	Role string
	Time	time.Time `json:"time,string,omitempty"`
	Description	string
	Operation int `json:"operation,string,omitempty"`
	Status	string
}

//交单
type LCTransDocsReceive struct{
	No int `json:"No,string,omitempty"`
	ReceivedAmount float64 `json:"ReceivedAmount,string,omitempty"`
	ReceivedDate time.Time `json:"ReceivedDate,string,omitempty"`
	BillOfLandings []BillOfLanding
	//提货单
	BillOfLadingDocs []Document
	HandOverBillStep string
	//描述
	Discrepancy string
	TransProgressFlow []TransProgress
}

//货运单
type BillOfLanding struct{
	//货运单编号
	BolNO string
	//货物编号
	GoodsNo string
	//货物描述
	GoodsDesc string
	//装船发运地
	//LoadPortName string
	//目的地
	//TransPortName string
	//最迟装船日
	//LatestShipDate string
	//是否分批装运
	//PartialShipment bool `json:"partialShipment,string"`
	//物流号
	//TrackingNo string
	//物流公司
	//Carrier Carrier
	//实际发货时间
	ShippingTime string
	//类别
	//Kind string
}

//货物
type GoodsInfo struct{
	//货物编号
	GoodsNo string
	//允许部分装运
	AllowPartialShipment int64 `json:"allowPartialShipment,string,omitempty"`
	//允许转船装运
	AllowTransShipment int64 `json:"allowTransShipment,string,omitempty"`
	//最迟装运日期
	LatestShipmentDate time.Time `json:"latestShipmentDate,string,omitempty"`
	//装运方式
	ShippingWay string
	//装运地点
	ShippingPlace string
	//目的地
	ShippingDestination string
	//贸易性质
	TradeNature int64 `json:"tradeNature,string,omitempty"`
	//货物描述
	GoodsDescription string
}

//合同
type Contract struct{
	Document
}

//信用证修改单
type AmendForm struct{
	//所属信用证
	//LC LCLetter
	//改证次数
	AmendTimes int64 `json:"amendTimes,string"`
	//修改后币种
	AmendedCurrency string
	//修改后金额
	AmendedAmt float64 `json:"amendedAmt,string"`
	//下浮允许度
	//LCAmtTolerDown int
	//上浮允许度
	//LCAmtTolerUp int
	//期限增减
	AddedDays int64 `json:"addedDays,string"`
	//改证后有效日期
	AmendExpiryDate string
	//改证日期
	//LCAmendDate time.Time
	//货物发送最终目的地
	TransPortName string
	//保证金增减金额
	AddedDepositAmt	float64 `json:"addedDepositAmt,string"`
	//详细描述
	//Details	string
}

//文件
type Document struct {
	//名称
	FileName string
	//文件的路径
	FileUri string
	//文件的Hash
	FileHash string
	//文件的签名
	FileSignature string
	//上传人
	Uploader string
}

//银行
type Bank struct{
	LegalEntity
	Address string
	AccountNo string
	AccountName string
	PostCode string
	Telephone string
	Telefax string
	Remark string
}
//企业
type Corporation struct {
	LegalEntity
	Account string
	DepositBank string
	Address string
	Nation string
	Contact string
	Email string
	PostCode string
	Telephone string
	Telefax string
	CreateTime string
}
type LegalEntity struct{
	No string
	Name string
	Domain string
}
type Carrier struct{
	LegalEntity
}

//银行、企业、银行企业签约的数据
type DataOfBCS struct{
	//类型+编号  类型 B=>银行 C=>企业 S=>签约
	No string
	//类型  类型说明 Bank=>银行 Corp=>企业 Sign=>签约
	Type string
	//银行信息
	DataBank Bank
	//企业信息
	DataCorp Corporation
	//签约状态 0 申请 1.通过 -1.拒绝
	StateSign int
	//签约时间
	SignDate string
}
