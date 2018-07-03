

package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
//信用证信息
type LCLetter struct {
	//L/C NO 信用证号(20字段) 
	LCNo string
	//发报行
	sendBank Bank
	//收报行
	recvBank Bank
	//开证申请人
	applyCorp Corporation
	//受益人
	benefCorp Corporation
	//开证币种
	currency string
	//开证金额
	amount float32
	//承兑金额
	acceptAmount float32
	//未付金额
	notPayAmount float32
	//开证日期
	issuseDate Date
	//到期日期
	expiryDate Date
	//货物描述
	goods GoodsInfo
	//单据要求
	documents [] LCDocument
	//交单期限，（天数）
	presentPeriod int
	//保证金
	depositAmt float32
	//改证次数
	amendTimes int
	//费用承担方
	chargeTaker Corporation
	//到单次数 
	ABTimes int
	//是否有效
	isValid bool
	//是否闭卷
	isClose bool
	//是否撤销
	isCancel bool

}	
//信用证申请信息
type LCApply struct{
	//申请编号
	applyNo string
	//信用证编号
	LCNO string 
	//合同
	contractNo Contract
	//申请人
	applyCorp Corporation
	//申请金额
	applyAmount float32
	//保证金比例
	depositPct float32
	//是否缴纳保证金
	isDeposited bool
	//申请日期
	applyDate Date
	
}
//银行
type Bank struct{
	//银行编号
	bankNo string
	//swift code
	bankSwiftCode string
	//名称
	bankName string
	//地址
	bankAddress string
	//平台签约编号
	chainNo string
}
//企业
type Corporation struct{
	//企业编号	
	corpNo string
	//名称
	corpName string
	//地址
	corpAddress string
	//swift code
	corpSwiftCode string
	//组织机构代码(9位)
	orgNo
	//统一社会信用代码(18位)
	uniNo

}
//货物
type GoodsInfo struct{
	//货物编号
	goodsNo string
	//货物描述
	goodsDesc string
	//装船发运地
	loadPortName string
	//目的地
	transPortName string
	//最迟装船日 
	latestShipDate string
	//是否分批装运
	partialShipment bool
}
//单据
type LCDocument struct{
	//单据编号
	docNo string
	//单据名称
	docName string
	//单据类型
	docType string
	//是否必须
	isRequired bool
	//影像Id
	fileHash string

}
//合同
type Contract struct{
	//继承自LCDocument
	LCDocument
	//合同编号
	contractNo string
	//买方
	purchaser Corporation
	//卖方
	vendor Corporation
	//合同金额
	amount float32
	//货物名称
	commodity  string
	//模板
	contractTemplateNo string
	//合同实体
	contractInstanceNo string

}
//企业账户信息
type Account struct{
	accountNo string 
	bank Bank
	corp Corporation
}
//信用证修改
type LCAmend struct{

}
//信用证撤销
type LCCancel struct{

}
//信用证闭卷
type LCClose struct{

}

//信用证不符点
type Discrepancy struct{
	discrepancyNo string 
	LC LCLetter
	//关联单据
	relatedDoc LCDocument
	//不符点
	discrepancies string 
	//是否接受
	isAccept bool

}

//信用证到单
type LCAB struct{
	ABNo string
	LC 	LCLetter
	//到单次数
	ABTimes int16
	//到单币种
	ABCurrency string
	//到单金额
	ABAmount float32
	//到单单据
	Document LCDocument
	//到单日期
	ABDate Dat
}

//信用证承兑
type LCAccept struct{

}
//信用证付款
type LCPay struct{

}
//信用证偿付
type LC
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
