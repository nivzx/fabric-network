package main

import (
	"encoding/json"
	"fmt"
	// "strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)

type TokenData struct {
	Token    string `json:"token"`
	Location string `json:"location"`
}

type TokenChaincode struct {
}

func (t *TokenChaincode) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (t *TokenChaincode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	switch function {
	case "initLedger":
		return t.initLedger(APIstub)
	case "writeToken":
		return t.writeToken(APIstub, args)
	case "getTokensByLocation":
		return t.getTokensByLocation(APIstub, args)
	default:
		return shim.Error("Invalid function name.")
	}
}

func (t *TokenChaincode) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	tokenData := []TokenData{
		TokenData{Token: "token1", Location: "6.089_80.073"},
		TokenData{Token: "token2", Location: "6.089_80.073"},
		TokenData{Token: "token3", Location: "6.089_80.073"},
	}

	for _, data := range tokenData {
		dataAsBytes, _ := json.Marshal(data)
		APIstub.PutState(data.Token, dataAsBytes)
	}

	return shim.Success(nil)
}

func (t *TokenChaincode) writeToken(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	token := args[0]
	location := args[1]

	// Check if token exists
	tokenDataAsBytes, err := APIstub.GetState(token)
	if err != nil {
		return shim.Error(err.Error())
	}

	var tokenData TokenData

	// If token exists, update the location
	if tokenDataAsBytes != nil {
		err = json.Unmarshal(tokenDataAsBytes, &tokenData)
		if err != nil {
			return shim.Error(err.Error())
		}
		tokenData.Location = location
	} else {
		// Create a new entry if token doesn't exist
		tokenData = TokenData{Token: token, Location: location}
	}

	tokenDataAsBytes, err = json.Marshal(tokenData)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = APIstub.PutState(token, tokenDataAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *TokenChaincode) getTokensByLocation(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	location := args[0]
	queryString := fmt.Sprintf(`{"selector":{"location":"%s"}}`, location)

	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var tokensInLocation []TokenData

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tokenData TokenData
		err = json.Unmarshal(queryResponse.Value, &tokenData)
		if err != nil {
			return shim.Error(err.Error())
		}

		tokensInLocation = append(tokensInLocation, tokenData)
	}

	tokensInLocationAsBytes, _ := json.Marshal(tokensInLocation)
	return shim.Success(tokensInLocationAsBytes)
}

func main() {
	err := shim.Start(new(TokenChaincode))
	if err != nil {
		fmt.Printf("Error creating new Token Chaincode: %s", err)
	}
}
