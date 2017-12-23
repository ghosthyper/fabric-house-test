package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

/*Struct*/
type House struct {
	HouseNum string `json:"housenum"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Owner    string `json:"owner"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	if function == "getFirstHouse" {
		return s.QueryHouse(APIstub, args)
	}
}

func (s *SmartContract) QueryHouse(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Wrong")
	}
	houseAsBytes, err := APIstub.GetState(args[0])
	if(err!= nil){
		reutn shim.Error(err)
	}
	return shim.Success(houseAsBytes)
}

func (s *SmartContract) GetAllHouse(APIstub shim.ChaincodeStubInterface) sc.Response{
	
	houseStart:="House0"
	houseEnd:="House9999"
	allHouseIterator,err:=APIstub.GetStateByRange(houseStart,houseEnd)

	//This will be called, when this function will reach to the point to return 
	defer allHouseIterator.Close()

	var buffer bytes.buffer
	buffer.WriteString("[")

	for allHouseIterator.HasNext(){
		queryRespone,err:=allHouseIterator.Next()
		buffer.WriteString("Key : "+queryRespone.Key)
		buffer.WriteString("Value is : "+ string(queryRespone.Value))
	}

	return shim.Success(buffer.bytes)
}


func (s *SmartContract) CreateHouse(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	if len(args)!=4{
		return shim.Error("Not enough values to create house object")
	}
	var houseObj = House{ HouseNum:args[0], Name:args[1], Address:args[2], Owner:args[3] }
	houseAsBytes,_:=json.Marshal(houseObj)
	APIstub.PutState(houseObj.HouseNum, houseAsBytes)
	return shim.Success(nil)
}

func (s *SmartContract) DeleteHouse(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	if len(args)!=1{
		return shim.Error("Error")
	}
	APIstub.DelState(args[0])
	return shim.Success(nil)
}

func(s *SmartContract) ChangeOwner(APIstub shim.ChaincodeStubInterface, args []string)sc.Response{
	if len(args)!=2{
		return shim.Error("Too less parameters")
	}

	//Get the current asset first
	houseAsBytes,error:=APIstub.GetState(args[0])
	if error!=nil{
		return shim.Error(error)
	}

	house:=House{}
	json.Unmarshal(houseAsBytes, &house)
	house.Owner=args[1]
	houseAsBytes,_:=json.Marshal(house)
	APIstub.PutState(args[0],houseAsBytes)
	return shim.Success(nil)


}