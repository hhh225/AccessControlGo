package bsnChainCode

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

var AggregatorOwner string

var reputationAddr string

var nResponse int

var nRequestOracles int

var respondingOracles []string = make([]string, 0)

var hashes []string = make([]string, 0)

var user string

var OraclesToDevices map[string][]string = make(map[string][]string)
var OracleList []string = make([]string, 0)

func Constructor(ownerAddr string, repAddr string) peer.Response {
	if AggregatorOwner != "" {
		return shim.Error("该合约已创建过了")
	}
	AggregatorOwner = ownerAddr
	reputationAddr = repAddr
	return shim.Success([]byte("成功创建合约"))
}

func AddOracle(oracleAddr string, devices []string) peer.Response {
	OracleList = append(OracleList, oracleAddr)
	OraclesToDevices[oracleAddr] = devices
	return shim.Success([]byte("成功添加预言机"))
}

//给预言机发送请求
func SendDataRequest(stub shim.ChaincodeStubInterface, userAddr string, iotAddr string, nOracles int) {
	nRequestOracles = nOracles
	user = userAddr
	counter := 0
	for i := 0; i < len(OracleList); i++ { //对于每个预言机
		devicesListTemp := OraclesToDevices[OracleList[i]]
		for j := 0; j < len(devicesListTemp); j++ {
			if iotAddr == devicesListTemp[j] {
				counter++
				//给oracle合约发送请求

				if counter == nRequestOracles {
					break
				}
			}

		}
	}
}

//被预言机合约调用
func OracleResponse(oracleAddr string, oracleresponseData string) {
	nResponse++
	respondingOracles = append(respondingOracles, oracleAddr)
	hashes = append(hashes, oracleresponseData)

	if nResponse == nRequestOracles {
		//调用
		ReportScore(respondingOracles, hashes, nRequestOracles)
		nResponse = 0
		respondingOracles = []string{}
		hashes = []string{}
	}

}
