/*
 * Copyright IBM Corp All Rights Reserved
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"strconv"
	"time"
)

type FileData struct {
	Id                  string `json:"id"`
	Filename            string `json:"filename"`
	FileDetail          string `json:"filedetail"`
	FileHash            string `json:"fileHash"`
	CreateDate          string `json:"createDate"`
	Username            string `json:"username"`
	Card                string `json:"card"`
	User                []string `json:"user"`
	Filepath            string `json:"filepath"`
	IpfsHash            string `json:"ipfsHash"`
	Time                string `json:"time"`
	Date                string `json:"date"`
	Money               string `json:"money"`
	TransferDate        string `json:"transferdate"`
	Flag                string `json:"flag"`
	Tmp                string `json:"tmp"`
}
type Transaction struct {
	Id                  string `json:"id"`
	Filename            string `json:"filename"`
	Date                string `json:"date"`
	Username            string `json:"username"`
	User                string `json:"user"`
	Username1           string `json:"username1"`
	User1               string `json:"user1"`
	Money               string `json:"money"`
	FileHash               string `json:"fileHash"`
	Flag                string `json:"flag"`
}
// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "setFileData" {
		err = setFileData(stub, args)
	} else if fn == "getFileData" {
		result, err = getFileData(stub, args)
	} else if fn == "getFileDetail" {
		result, err = getFileDetail(stub, args)
	} else if fn == "setTransaction" {
		result, err = setTransaction(stub, args)
	} else if fn == "setMoney" {
		result, err = setMoney(stub, args)
	} else if fn == "getMoney" {
		result, err = getMoney(stub, args)
	} else if fn == "getTransaction" {
		result, err = getTransaction(stub, args)
	} else if fn == "set" {
		result, err = set(stub, args)
	} else if fn == "delete" {
		result, err = delete(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
// 将对象序列化后保存至账本中
func setFileData(stub shim.ChaincodeStubInterface, args []string) (error) {
	var fileData FileData
	var fileDatas []FileData
	err := json.Unmarshal([]byte(args[0]), &fileData)
	if err != nil {
		return fmt.Errorf("Failed to json asset: %s", err)
	}
	assetJSON, err := stub.GetState("FileData")
	if err != nil {
		return fmt.Errorf("Failed to set asset: %s", err)
	}
	if assetJSON == nil {
		fileDatas = append(fileDatas,fileData)
	}else {
		err = json.Unmarshal(assetJSON, &fileDatas)
		if err != nil {
			return fmt.Errorf("Failed to json asset: %s", err)
		}
		fileDatas = append(fileDatas,fileData)
	}
	putData, err := json.Marshal(fileDatas)
	if err != nil {
		return fmt.Errorf("Failed to json asset: %s", err)
	}
	err = stub.PutState("FileData", putData)
	if err != nil {
		return fmt.Errorf("Failed to set asset: %s", err)
	}
	return nil
}
// 查询对应的值
func getFileDetail(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	id := args[0]
	var fileDatas []FileData
	assetJSON, err := stub.GetState("FileData")
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s", err)
	}
	err = json.Unmarshal(assetJSON, &fileDatas)
	if err != nil {
		return "",fmt.Errorf("Failed to json asset: %s", err)
	}
	for _, fd := range fileDatas{
		if id == fd.Id {
			putData, err := json.Marshal(fd)
			if err != nil {
				return "", fmt.Errorf("Failed to json asset: %s", err)
			}
			return string(putData), nil
		}
	}
	return "", fmt.Errorf("Failed to get asset: %s", err)
}
func getFileData(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	assetJSON, err := stub.GetState("FileData")
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s", err)
	}
	return string(assetJSON), nil
}
func setTransaction(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	id := args[1]
	var fileDatas []FileData
	var newFileDatas []FileData
	var transactions []Transaction
	var transaction Transaction
	assetJSON, err := stub.GetState("FileData")
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s", err)
	}
	err = json.Unmarshal(assetJSON, &fileDatas)
	if err != nil {
		return "",fmt.Errorf("Failed to json1 asset: %s", err)
	}
	for _, fd := range fileDatas{
		if id == fd.Id {
			//transaction
			l := len(fd.User)
			transaction.Id = fd.Id
			transaction.Filename = fd.Filename
			transaction.Username = fd.Username
			transaction.User = fd.User[l-1]
			transaction.Money = fd.Money
			transaction.FileHash = fd.FileHash
			//fd
			fd.Username = args[2]
			fd.User = append(fd.User, args[0])
			fd.Flag = "已交易"
			fd.Date = args[3]
			fd.Card = args[4]
		}
		newFileDatas = append(newFileDatas, fd)
	}
	putData, err := json.Marshal(newFileDatas)
	if err != nil {
		return "", fmt.Errorf("Failed to json2 asset: %s", err)
	}
	err = stub.PutState("FileData", putData)
	if err != nil {
		return "",fmt.Errorf("Failed to set asset: %s", err)
	}

	transaction.Date = args[3]
	transaction.Username1 = args[2]
	transaction.User1 = args[0]
	assetJSON, err = stub.GetState(id)
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s", err)
	}
	if string(assetJSON) == "" {
		transactions = append(transactions,transaction)
	}else {
		err = json.Unmarshal(assetJSON, &transactions)
		if err != nil {
			return "", fmt.Errorf("Failed to json3 asset: %s", err)
		}
		transactions = append(transactions,transaction)
	}
	putData, err = json.Marshal(transactions)
	if err != nil {
		return "", fmt.Errorf("Failed to json asset: %s", err)
	}
	err = stub.PutState(id, putData)
	if err != nil {
		return "",fmt.Errorf("Failed to set asset: %s", err)
	}

	a := []string{transaction.User1, transaction.Money, "down"}
	money, err := setMoney(stub, a)
	if err != nil {
		return "",fmt.Errorf("Failed to setMoney asset: %s", err)
	}
	return money, nil
}


func setMoney(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if args[2] == "init"{
		err := stub.PutState(args[0], []byte(args[1]))
		if err != nil {
			return "",fmt.Errorf("Failed to set asset: %s", err)
		}
		return "init success", nil
	}
	value, err := stub.GetState(args[0])
	if err != nil {
		return "",fmt.Errorf("Failed to get asset: %s", err)
	}
	time.Sleep(10)
	int1 ,err := strconv.Atoi(string(value))
	int2 ,err := strconv.Atoi(args[1])
	int3 := int1 + int2
	if args[2] == "down"{
		int3 = int1 - int2
	}
	money := strconv.Itoa(int3)
	err = stub.PutState(args[0], []byte(money))
	if err != nil {
		return "",fmt.Errorf("Failed to set asset: %s", err)
	}
	return money, nil
}
func getMoney(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	value, err := stub.GetState(args[0])
	if err != nil {
		return "",fmt.Errorf("Failed to get asset: %s", err)
	}
	int1 ,err := strconv.Atoi(string(value))
	int2 ,err := strconv.Atoi(args[1])
	if int1 < int2 {
		return string(value), fmt.Errorf("not enough")
	}
	return string(value), nil
}
func getTransaction(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	assetJSON, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s", err)
	}
	return string(assetJSON), nil
}
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	id := args[0]
	var fileDatas []FileData
	var nfileDatas []FileData
	assetJSON, err := stub.GetState("FileData")
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s", err)
	}
	err = json.Unmarshal(assetJSON, &fileDatas)
	if err != nil {
		return "",fmt.Errorf("Failed to json asset: %s", err)
	}
	for _, fd := range fileDatas{
		if id == fd.Id {
			fd.Money = args[1]
			fd.Tmp = fd.Flag
			fd.Flag = "待交易"
			fd.TransferDate = args[2]
		}
		nfileDatas = append(nfileDatas, fd)
	}
	putData, err := json.Marshal(nfileDatas)
	if err != nil {
		return "", fmt.Errorf("Failed to json asset: %s", err)
	}
	err = stub.PutState("FileData", putData)
	if err != nil {
		return "",fmt.Errorf("Failed to set asset: %s", err)
	}
	return "success", nil
}
func delete(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	id := args[0]
	var fileDatas []FileData
	var nfileDatas []FileData
	assetJSON, err := stub.GetState("FileData")
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s", err)
	}
	err = json.Unmarshal(assetJSON, &fileDatas)
	if err != nil {
		return "",fmt.Errorf("Failed to json asset: %s", err)
	}
	for _, fd := range fileDatas{
		if id == fd.Id {
			fd.Money = ""
			fd.Flag = fd.Tmp
			fd.TransferDate = ""
		}
		nfileDatas = append(nfileDatas, fd)
	}
	putData, err := json.Marshal(nfileDatas)
	if err != nil {
		return "", fmt.Errorf("Failed to json asset: %s", err)
	}
	err = stub.PutState("FileData", putData)
	if err != nil {
		return "",fmt.Errorf("Failed to set asset: %s", err)
	}
	return "success", nil
}
// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
