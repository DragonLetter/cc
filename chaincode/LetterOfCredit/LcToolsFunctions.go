package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"errors"
	"math"
	"fmt"
	"bytes"
)

const MIN = 0.000001

func getNextSequence(stub shim.ChaincodeStubInterface,formPrefix string) int {
	key:=formPrefix+"Sequence"
	lcSeqAsBytes, err := stub.GetState(key)
	if err != nil {
		shim.Error("Failed to get Sequence: " + err.Error())
	}
	seq, _ := strconv.Atoi(string(lcSeqAsBytes))
	stub.PutState(key,[]byte(strconv.Itoa(seq+1)))
	return seq
}

func decodeApplicationForm(jsonStr string) (ApplicationForm,error){
	var form ApplicationForm
	err:=json.Unmarshal([]byte (jsonStr),&form)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to ApplicationForm}"
		return ApplicationForm{},errors.New(jsonResp)
	}
	return form,nil
}
func decodeLetterOfCredit(jsonStr string) (LetterOfCredit, error)  {
	var letterOfCredit LetterOfCredit
	err:=json.Unmarshal([]byte (jsonStr),&letterOfCredit)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to LetterOfCredit}"
		return LetterOfCredit{},errors.New(jsonResp)
	}
	return letterOfCredit,nil
}
func decodeLCTransDeposit(jsonStr string) (LCTransDeposit, error)  {
	var deposit LCTransDeposit
	err:=json.Unmarshal([]byte (jsonStr),&deposit)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to LCTransDeposit}"
		return LCTransDeposit{},errors.New(jsonResp)
	}
	return deposit,nil
}
func decodeTransaction(jsonStr string) (Transaction, error)  {
	var transaction Transaction
	err:=json.Unmarshal([]byte (jsonStr),&transaction)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to Transaction}"
		return Transaction{},errors.New(jsonResp)
	}
	return transaction,nil
}
func decodeTransProgress(jsonStr string) (TransProgress, error)  {
	var transProgress TransProgress
	err:=json.Unmarshal([]byte (jsonStr),&transProgress)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to TransProgress}"
		return TransProgress{},errors.New(jsonResp)
	}
	return transProgress,nil
}
func decodeCorp(jsonStr string) (Corporation,error){
	var corp Corporation
	err:=json.Unmarshal([]byte (jsonStr),&corp)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to Corporation}"
		return Corporation{},errors.New(jsonResp)
	}
	return corp,nil
}
func decodeContract(jsonStr string) (Contract,error){
	var contract Contract
	err:=json.Unmarshal([]byte (jsonStr),&contract)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to Contract}"
		return Contract{},errors.New(jsonResp)
	}
	return contract,nil
}
func decodeBank(jsonStr string) (Bank,error){
	var bank Bank
	err:=json.Unmarshal([]byte (jsonStr),&bank)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to Bank}"
		return Bank{},errors.New(jsonResp)
	}
	return bank,nil
}
func decodeLCTransDocsReceive(jsonStr string) (LCTransDocsReceive,error){
	var lCTransDocsReceive LCTransDocsReceive
	err:=json.Unmarshal([]byte (jsonStr),&lCTransDocsReceive)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to LCTransDocsReceive}"
		return LCTransDocsReceive{},errors.New(jsonResp)
	}
	return lCTransDocsReceive,nil
}
func decodeBillOfLanding(jsonStr string) (BillOfLanding,error){
	var billOfLanding BillOfLanding
	err:=json.Unmarshal([]byte (jsonStr),&billOfLanding)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to BillOfLanding}"
		return BillOfLanding{},errors.New(jsonResp)
	}
	return billOfLanding,nil
}
func decodeGoods(jsonStr string) (GoodsInfo,error){
	var goods GoodsInfo
	err:=json.Unmarshal([]byte (jsonStr),&goods)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to Goods}"
		return GoodsInfo{},errors.New(jsonResp)
	}
	return goods,nil
}
func decodeLegalEntity(jsonStr string) (LegalEntity,error){
	var bank LegalEntity
	err:=json.Unmarshal([]byte (jsonStr),&bank)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to LegalEntity}"
		return LegalEntity{},errors.New(jsonResp)
	}
	return bank,nil
}
func decodeDocument(jsonStr string) (Document,error){
	var document Document
	err :=json.Unmarshal([]byte (jsonStr),&document)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to document}"
		return document,errors.New(jsonResp)
	}
	return document,nil
}
func decodeDocuments(jsonStr string) ([]Document,error){
	var documents []Document
	err :=json.Unmarshal([]byte (jsonStr),&documents)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to documents}"
		return documents,errors.New(jsonResp)
	}
	return documents,nil
}
func decodeLCLetter(jsonStr string) (LetterOfCredit,error){
	var lcLetter LetterOfCredit
	err :=json.Unmarshal([]byte (jsonStr),&lcLetter)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to lcLetter}"
		return lcLetter,errors.New(jsonResp)
	}
	return lcLetter,nil
}
func decodeAmendForm(jsonStr string) (AmendForm,error){
	var amendForm AmendForm
	err :=json.Unmarshal([]byte (jsonStr),&amendForm)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to AmendForm}"
		return amendForm,errors.New(jsonResp)
	}
	return amendForm,nil
}
func decodeBCSData(jsonStr string) (DataOfBCS,error){
	var data DataOfBCS
	err:=json.Unmarshal([]byte (jsonStr),&data)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + jsonStr+ "\" to DataOfBCS}"
		return DataOfBCS{},errors.New(jsonResp)
	}
	return data,nil
}

// MIN 为用户自定义的比较精度
func isEqual(f1, f2 float64) bool {
	return math.Dim(f1, f2) < MIN
}

//Query
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
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

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}
