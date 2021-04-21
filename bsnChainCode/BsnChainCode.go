package bsnChainCode

import (
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

var Count = 0
var Owner = ""
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
		response := stub.InvokeChaincode("", toChaincodeArgs2(addrUser, addrDevice, nOracle), "")
		//返回预言机地址
		return shim.Success(response.Payload)
	}
}

func toChaincodeArgs2(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

type BsnChainCode struct {
}

func (t *BsnChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {

	dbBaseModel := models.DBBaseModel{BaseKey: "cc_key_", BaseInfo: "Welcome to use ChainCode "}
	reqJsonValue, err := json.Marshal(&dbBaseModel)
	if err != nil {
		return shim.Error(fmt.Sprintf("数据转换失败:%s", err.Error()))
	}
	err = stub.PutState(dbBaseModel.BaseKey, reqJsonValue)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *BsnChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(function)
	fmt.Println(args)

	switch function {
	case "Constructor": // 保存
		return Construct(args[0], args[1])
	case "AddDevice":
		return AddDevice(args[0])
	case "AddAdmin":
		return AddAdmin(args[0])
	case "AddUserDevice":
		return AddUserDevice(args[0], args[1])
	case "DelAdmin":
		return DelAdmin(args[0])
	case "DelUser":
		return DelUser(args[0])
	case "RequestUserToAccessDevices":
		//nOra,_:=strconv.ParseInt(args[2],10,0)
		return RequestUserToAccessDevicex(stub, args[0], args[1], args[2])

	default:
		//SetLogger("无效的方法")
		break
	}

	return shim.Success([]byte("fuck u"))
}
