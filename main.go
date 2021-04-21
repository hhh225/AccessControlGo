package main

import (
	"awesomeProject/bsnChainCode"

	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

func main() {
	err := shim.Start(new(bsnChainCode.BsnChainCode))
	if err != nil {
		fmt.Println("error starting")
	}
}
