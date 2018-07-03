package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"encoding/pem"
	"crypto/x509"
	"bytes"
	"strconv"
	"encoding/json"
	"strings"
	//"github.com/golang/protobuf/proto"
	//"golang.org/x/tools/go/gcimporter15/testdata"
)

type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}


func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)
	if function == "test1" {//自定义函数名称
		return t.testCertificate(stub, args)//定义调用的函数
	}else if function=="test21"{
		return t.testStateInsert(stub,args)
	}else if function=="test22"{
		return t.testStateQuery(stub,args)
	}else if function=="test23"{
		return t.testStateDelete(stub,args)
	}	else if function=="test3"{
	return t.testRangeQuery(stub,args)
	}	else if function=="test4"{
	return t.testRichQuery(stub,args)
	}else if function=="test5"{
		return t.testCompositeKey(stub,args)
	}else if function=="test6"{
		return t.testHistoryQuery(stub,args)
	}else if function=="test7"{
		return t.testInvokeChainCode(stub,args)
	}else if function=="test8"{
		return t.testGetProposal(stub,args)
	}else if function=="test9"{
		return t.testEvent(stub,args)
	}
	return shim.Error("Received unknown function invocation")
}
func (t *SimpleChaincode) testCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	creatorByte,_:= stub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----")// Devin:I don't know why sometimes -----BEGIN is invalid, so I use -----
	if certStart == -1 {
		fmt.Errorf("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		fmt.Errorf("Could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		fmt.Errorf("ParseCertificate failed")
	}
	uname:=cert.Subject.CommonName
	fmt.Println("Name:"+uname)
	return shim.Success([]byte("Called testCertificate "+uname))
}
type Student struct {
	Id int
	Name string
}
func (t *SimpleChaincode) testStateInsert(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	student1:=Student{1,"Devin Zeng"}
	key:="Student:"+strconv.Itoa(student1.Id)
	studentJsonBytes, err := json.Marshal(student1)
	if err != nil {
		return shim.Error(err.Error())
	}
	err= stub.PutState(key,studentJsonBytes)
	if(err!=nil){
		return shim.Error(err.Error())
	}
	student2:=Student{2,"Edward"}
	studentJsonBytes, err = json.Marshal(student2)
	stub.PutState("Student:2",studentJsonBytes)
	return shim.Success([]byte("Saved Student!"))
}
func (t *SimpleChaincode) testStateQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	key:="Student:1"
	dbStudentBytes,err:= stub.GetState(key)
	var dbStudent Student;
	err=json.Unmarshal(dbStudentBytes,&dbStudent)
	if err != nil {
		return shim.Error("{\"Error\":\"Failed to decode JSON of: " + string(dbStudentBytes)+ "\" to Student}")
	}
	fmt.Println("Read Student from DB, name:"+dbStudent.Name)
	return shim.Success(dbStudentBytes)
}
func (t *SimpleChaincode) testStateDelete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	key := "Student:1"
	err:= stub.DelState(key)
	if err != nil {
		return shim.Error("1 Failed to delete Student from DB, key is: "+key)
	}else{
		fmt.Println("删除Student成功")
	}
	err= stub.DelState(key)
	if err != nil {
		return shim.Error("2 Failed to delete Student from DB, key is: "+key)
	}else{
		fmt.Println("删除Student成功")
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) testRangeQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	resultsIterator,err:= stub.GetStateByRange("Student:1","Student:3")
	if err!=nil{
		return shim.Error("Query by Range failed")
	}
	students,err:=getListResult(resultsIterator)
	if err!=nil{
		return shim.Error("getListResult failed")
	}
	return shim.Success(students)
}
func (t *SimpleChaincode) testRichQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	name:="Devin Zeng"
	queryString := fmt.Sprintf("{\"selector\":{\"Name\":\"%s\"}}", name)
	resultsIterator,err:= stub.GetQueryResult(queryString)
	if err!=nil{
		return shim.Error("Rich query failed")
	}
	students,err:=getListResult(resultsIterator)
	if err!=nil{
		return shim.Error("Rich query failed")
	}
	return shim.Success(students)
}
func getListResult(resultsIterator shim.StateQueryIteratorInterface) ([]byte,error){

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
	fmt.Printf("queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}
func (t *SimpleChaincode) testCompositeKey(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	cc:=ChooseCourse{"CS101",123,true}
	var key1,_= stub.CreateCompositeKey("ChooseCourse",[]string{cc.CourseNumber,strconv.Itoa(cc.StudentId)})
	var key2,_= stub.CreateCompositeKey("ChooseCourse",[]string{"PH209","123"})
	fmt.Println(key1)
	objType,attrArray,_:= stub.SplitCompositeKey(key1)
	fmt.Println("Object:"+objType+" ,Attributes:"+strings.Join(attrArray,"|"))
	fmt.Println(key2)
	stub.PutState(key1,[]byte("CS101,123"))
	stub.PutState(key2,[]byte("PH209,123"))
	 resultsIterator,_:= stub.GetStateByPartialCompositeKey("ChooseCourse",[]string{"CS101"})
	result,_:= getListResult(resultsIterator)
	return shim.Success(result)
}

type ChooseCourse struct {
	CourseNumber string //开课编号
	StudentId int //学生ID
	Confirm bool //是否确认
}
func (t *SimpleChaincode) testHistoryQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	student1:=Student{1,"Devin Zeng"}
	key:="Student:"+strconv.Itoa(student1.Id)
	it,err:= stub.GetHistoryForKey(key)
	if err!=nil{
		return shim.Error(err.Error())
	}
	var result,_= getHistoryListResult(it)
	return shim.Success(result)
}
func getHistoryListResult(resultsIterator shim.HistoryQueryIteratorInterface) ([]byte,error){

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
		item,_:= json.Marshal( queryResponse)
		buffer.Write(item)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}
func (t *SimpleChaincode) testInvokeChainCode(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	trans:=[][]byte{[]byte("invoke"),[]byte("a"),[]byte("b"),[]byte("11")}
	response:= stub.InvokeChaincode("mycc",trans,"mychannel")
	fmt.Println(response.Message)
	return shim.Success([]byte( response.Message))

}
func (t *SimpleChaincode) testGetProposal(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	var student,_= stub.GetState("Student:1")
	fmt.Println(student)
	stub.PutState("ChooseCourse1",[]byte("CS101,123"))
	p,_:=stub.GetSignedProposal()
	fmt.Println(p.String())
	//var pr pb.Proposal
	// err:= proto.Unmarshal( p.ProposalBytes,&pr)
	//if(err!=nil){
	//	return shim.Error(err.Error())
	//}
	//prJson,err:=json.Marshal(pr)
	//if err!=nil{
	//	return shim.Error(err.Error())
	//}
	return shim.Success([]byte(p.String()))
}
func (t *SimpleChaincode) testEvent(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	tosend := "Event send data is here!"
	err := stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}