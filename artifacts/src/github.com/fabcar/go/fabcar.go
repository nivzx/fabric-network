package main

import (
	"encoding/json"
	"fmt"
	// "bytes"
	// "net/http"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)

type SmartContract struct {
}

type SignalData struct {
	Location string  `json:"location"`
	Level    float64 `json:"level"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	switch function {
	case "writeLevel":
		return s.writeLevel(APIstub, args)
	case "readLevel":
		return s.readLevel(APIstub, args)
	case "getAllLevels":
		return s.getAllLevels(APIstub)
	case "initLedger":
		return s.initLedger(APIstub)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}

func (s *SmartContract) writeLevel(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	location := args[0]
	level := args[1]

	levelValue, err := strconv.ParseFloat(level, 64)
	if err != nil {
		return shim.Error("Invalid level value. " + err.Error())
	}

	// if levelValue > -50.00 {
	// 	// Create the JSON data for the HTTP request
	// 	requestData := map[string]interface{}{
	// 		"location": location,
	// 	}

	// 	requestBytes, err := json.Marshal(requestData)
	// 	if err != nil {
	// 		return shim.Error("Failed to marshal JSON request. " + err.Error())
	// 	}

	// 	// Send the POST request
	// 	_, err = http.Post("http://172.18.0.2:8080", "application/json", bytes.NewBuffer(requestBytes))
	// 	if err != nil {
	// 		return shim.Error("Failed to send POST request. " + err.Error())
	// 	}
	// }

	signalData := SignalData{Location: location, Level: levelValue}
	signalDataAsBytes, _ := json.Marshal(signalData)

	err = APIstub.PutState(location, signalDataAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) readLevel(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	location := args[0]
	signalDataAsBytes, err := APIstub.GetState(location)
	if err != nil {
		return shim.Error(err.Error())
	}
	if signalDataAsBytes == nil {
		return shim.Error("Signal data not found for location: " + location)
	}

	var signalData SignalData
	err = json.Unmarshal(signalDataAsBytes, &signalData)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(signalDataAsBytes)
}

func (s *SmartContract) getAllLevels(APIstub shim.ChaincodeStubInterface) sc.Response {
	resultsIterator, err := APIstub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var allSignalData []SignalData

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var signalData SignalData
		err = json.Unmarshal(queryResponse.Value, &signalData)
		if err != nil {
			return shim.Error(err.Error())
		}

		allSignalData = append(allSignalData, signalData)
	}

	allSignalDataAsBytes, _ := json.Marshal(allSignalData)
	return shim.Success(allSignalDataAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	signalData := []SignalData{
		SignalData{Location: "80.9365_6.0937", Level: -98.847},
		SignalData{Location: "81.0245_6.2057", Level: -75.982},
		SignalData{Location: "80.8112_6.4930", Level: -52.365},
	}

	for _, data := range signalData {
		dataAsBytes, _ := json.Marshal(data)
		APIstub.PutState(data.Location, dataAsBytes)
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
