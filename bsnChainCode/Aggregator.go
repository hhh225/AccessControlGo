package bsnChainCode

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

var AggregatorOwner string

var reputationAddr string

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

//给预言机发送请求，并接收结果返回
func SendDataRequest(userAddr string, iotAddr string, nOracles int) string {
	//var nResponse int

	var nRequestOracles int

	var respondingOracles []string = make([]string, 0)

	var hashes []string = make([]string, 0)

	//var user string
	nRequestOracles = nOracles
	//user = userAddr
	counter := 0
	var resultOracle string = ""
	for i := 0; i < len(OracleList); i++ { //对于每个预言机
		devicesListTemp := OraclesToDevices[OracleList[i]]
		for j := 0; j < len(devicesListTemp); j++ { //对于这个预言机所能访问的所有物联网设备
			if iotAddr == devicesListTemp[j] { //如果有某个物联网设备是我们所要的物联网设备
				counter++
				//给oracle合约发送请求，接收结果
				oracleRespondingData := Query(iotAddr)

				respondingOracles = append(respondingOracles, OracleList[i])
				hashes = append(hashes, oracleRespondingData)
				if counter == nRequestOracles {
					resultOracle = ReportScore(respondingOracles, hashes, nRequestOracles)
					//nResponse = 0

					goto End
				}
			}

		}
	}
End:
	return resultOracle
}
