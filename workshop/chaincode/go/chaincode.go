package main

import (
	"strings"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type MyGoChaincode struct {
}

type AppLicense struct {
	LicenseId string `json:"licenseId"`
	Name string `json:"name"`
	Organization string `json:"organization"`
	LicensedTo string `json:"licensedTo"`
}

func (t *MyGoChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("****Init Chaincode****")
	return shim.Success(nil)
}

func (t *MyGoChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("****Invoke Chaincode****")
	function, args := stub.GetFunctionAndParameters()
	if function == "create" {
		return t.createLicenses(stub)
	} else if function == "transfer" {
		return t.transferLicense(stub, args)
	} else if function == "query" {
		return t.queryLicense(stub, args)
	}

	return shim.Error("Invalid function name")
}

func (t *MyGoChaincode) createLicenses(stub shim.ChaincodeStubInterface) pb.Response {
	licenses := []AppLicense{
		AppLicense{LicenseId: "abc123", Name: "MyWordProcessor", Organization: "MicroSuave", LicensedTo: "Billy"},
		AppLicense{LicenseId: "abc456", Name: "MyWordProcessor", Organization: "MicroSuave", LicensedTo: "Billy"},
		AppLicense{LicenseId: "abc789", Name: "MyWordProcessor", Organization: "MicroSuave", LicensedTo: "Billy"},
		AppLicense{LicenseId: "abc987", Name: "MyWordProcessor", Organization: "MicroSuave", LicensedTo: "Billy"},
		AppLicense{LicenseId: "abc654", Name: "MyWordProcessor", Organization: "MicroSuave", LicensedTo: "Billy"},
		AppLicense{LicenseId: "abc321", Name: "MyWordProcessor", Organization: "MicroSuave", LicensedTo: "Billy"},
	}

	i := 0
	for i < len(licenses) {
		licenseAsBytes, err := json.Marshal(licenses[i])

		if err != nil {
			jsonResp := "{\"Error\":\"Failed to parse asset on init\"}"
			return shim.Error(jsonResp)
		}

		stub.PutState(licenses[i].LicenseId, licenseAsBytes)
		fmt.Println("License Added", licenses[i].LicenseId)
		i = i + 1
	}

	return shim.Success(nil)
}

func (t *MyGoChaincode) transferLicense(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	licenseAsBytes, err := stub.GetState(args[0])

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get current state for license " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	license := AppLicense{}

	json.Unmarshal(licenseAsBytes, &license)

	if strings.Compare(license.LicensedTo, "Billy") != 0 {
		jsonResp := "{\"Error\":\"The license has been transferred to " + license.LicensedTo + "\"}"
		return shim.Error(jsonResp)
	}

	license.LicensedTo = args[1]

	licenseAsBytes, err = json.Marshal(license)

	if err != nil {
		jsonResp := "{\"Error\":\"Error parsing to json\"}"
		return shim.Error(jsonResp)
	}

	stub.PutState(args[0], licenseAsBytes)

	return shim.Success(nil)

}

func (t *MyGoChaincode) queryLicense(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	licenseAsBytes, err := stub.GetState(args[0])
	if err != nil {
		jsonResp := "{\"Error\":\"Error getting the state for license " + args[0] +  "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Println("License Data: " + string(licenseAsBytes))

	return shim.Success(licenseAsBytes)

}

func main() {
	err := shim.Start(new(MyGoChaincode))
	if err != nil {
		fmt.Printf("Error creating new chaincode: %s", err)
	}
}