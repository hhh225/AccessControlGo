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
		//调用预言机智能合约

	}
}

var Aggrement map[string]int = make(map[string]int)                  //查看每个回复的重复次数
var DuplicateOracles map[string][]string = make(map[string][]string) //这个回复的所有预言机
var TrueOracles []string = make([]string, 0)                         //所有这个回复的预言机

/*
respondOracles代表所有有回复的预言机
dataHashes代表每个respondOracles的回复
*/
func ReportScore(respondOracles []string, dataHashes []string) {
	var nMatches int = int(nRequestOracles) / 2
	var correctHash string
	var points []int

	for i := 0; i < len(dataHashes); i++ {
		if Aggrement[dataHashes[i]] != 0 {
			Aggrement[dataHashes[i]]++
			DuplicateOracles[dataHashes[i]] = append(DuplicateOracles[dataHashes[i]], respondOracles[i])
		} else {
			Aggrement[dataHashes[i]] = 1
			DuplicateOracles[dataHashes[i]] = []string{respondOracles[i]}
		}
	}
	for i := 0; i < len(dataHashes); i++ {
		if Aggrement[dataHashes[i]] >= nMatches {
			correctHash = dataHashes[i]
			for j := 0; j < len(DuplicateOracles[correctHash]); j++ {
				TrueOracles = append(TrueOracles, DuplicateOracles[correctHash][j])
			}
			break
		}
	}
	for i := 0; i < len(dataHashes); i++ {
		for j := 0; j < len(TrueOracles); j++ {
			if respondOracles[i] == TrueOracles[j] {
				points = append(points, 100)
				break
			} else {
				points = append(points, 0)
			}
		}
	}
}
