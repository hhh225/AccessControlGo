package bsnChainCode

import (
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

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
