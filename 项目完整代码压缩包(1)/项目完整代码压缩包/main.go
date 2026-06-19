package main

import (
	"fmt"
	"os"
	"testwork/sdkInit"
	"testwork/web"
)
const (
	cc_name = "test"
	cc_version = "1.0"
)

var App sdkInit.Application //app
//数据上传到ipfs
func main()  {
	//org信息
	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/root/testwork/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/root/testwork/fixtures/channel-artifacts/Org2MSPanchors.tx",
		},
	}
	// 初始化info
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychanne1",
		ChannelConfig:    "/root/testwork/fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer1.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    "/root/testwork/chaincode/go/test",
		ChaincodeVersion: cc_version,
	}
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil{
		fmt.Println(">> Sdk set error ", err)
		os.Exit(-1)
	}
	fmt.Println(">> Set the chain code status through the chain code external service......")
	if err := info.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk);err != nil{
		fmt.Println(">> InitService error: ", err)
		os.Exit(-1)
	}
	App=sdkInit.Application{
		SdkEnvInfo: &info,
	}
	fmt.Println(">> Set chain code status completed")
	web.WebStart(App)

}
