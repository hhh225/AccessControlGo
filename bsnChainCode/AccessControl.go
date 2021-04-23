package bsnChainCode

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"strconv"
)

var Count = 0
var AccessControlOwner = ""
var AggregatorSC = ""
var Admins = make([]string, 0)

type Token struct {
	User     string
	Device   string
	NOracles int
}

var Tokens = make([]Token, 0)
var Devices = make([]string, 0)
var User_Devices map[string][]string

func Construct(addrOwner string, addrAggregatorSC string) peer.Response {
	Admins = append(Admins, addrOwner)
	AggregatorSC = addrAggregatorSC
	return shim.Success([]byte("创建成功"))
}

//新增物联网服务
func AddDevice(newDevice string) peer.Response {
	Devices = append(Devices, newDevice)
	return shim.Success([]byte("新增服务成功"))
}

func AddAdmin(addrAdmin string) peer.Response {
	Admins = append(Admins, addrAdmin)
	return shim.Success([]byte("新增管理员成功"))
}

func AddUserDevice(userAddr string, deviceAddr string) peer.Response {
	deviceExist := false
	for i := 0; i < len(Devices); i++ {
		if Devices[i] == deviceAddr {
			deviceExist = true
			break
		}
	}
	if deviceExist {
		User_Devices[userAddr] = append(User_Devices[userAddr], deviceAddr)
		return shim.Success([]byte("添加设备成功"))
	} else {
		return shim.Error("添加设备失败，设备不存在")
	}
}

func DelAdmin(addrAdmin string) peer.Response {
	i := 0
	for i = 0; i < len(Admins); i++ {
		if Admins[i] == addrAdmin {
			Admins = append(Admins[:i], Admins[i+1:]...)
			return shim.Success([]byte("成功删除管理员"))
			break

		}
	}
	return shim.Error("此管理员不存在")
}

func DelUser(addrUser string) peer.Response {
	if User_Devices[addrUser] != nil {
		delete(User_Devices, addrUser)
		return shim.Success([]byte("成功删除用户"))
	}
	return shim.Error("此用户不存在")
}

func toChaincodeArgs2(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func RequestUserToAccessDevicex(stub shim.ChaincodeStubInterface, addrUser string, addrDevice string, nOracle string) peer.Response {
	i := 0
	for i = 0; i < len(Devices); i++ {
		if Devices[i] == addrDevice {

			break
		}
	}
	if i == len(Devices) {
		return shim.Error("该设备不存在")
	}
	devicesList := User_Devices[addrUser]
	if devicesList == nil {
		return shim.Error("用户不存在")
	}
	for i = 0; i < len(devicesList); i++ {
		if devicesList[i] == addrDevice {
			break
		}
	}
	if i == len(devicesList) {
		return shim.Error("用户不能访问该设备")
	} else {
		//调用aggregator
		nOracleInt, _ := strconv.ParseInt(nOracle, 10, 0)
		SendDataRequest(addrUser, addrDevice, int(nOracleInt))
		//返回预言机地址
		return shim.Success([]byte("预言机地址"))
	}
}
