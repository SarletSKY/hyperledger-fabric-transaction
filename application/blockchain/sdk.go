package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// sdk 相关配置
var (
	SDK           *fabsdk.FabricSDK             // FabricSDK
	ChannelName   = "transaction"               // 通道名称
	ChainCodeName = "food"                      // 链码名称
	Org           = "org1"                      // 组织名称
	User          = "Admin"                     // 用户
	ConfigPath    = "./application/config.yaml" // 配置文件的路径
)

// SDK初始化
func Init() {
	var err error
	// 通过配置文件初始化SDK
	SDK, err = fabsdk.New(config.FromFile(ConfigPath))
	if err != nil {
		panic(err)
	}
}

// 区块链交互
func ChannelExecute(fcn string, args [][]byte) (channel.Response, error) {
	ctx := SDK.ChannelContext(ChannelName, fabsdk.WithUser(User), fabsdk.WithOrg(Org))
	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}

	// 客户端进行交互
	resp, err := cli.Execute(channel.Request{
		ChaincodeID: ChainCodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints("word.node1.gdzc.com", "weixiong.node2.gdzc.com", "peer1.node3.gdzc.com"))
	if err != nil {
		return channel.Response{}, err
	}

	return resp, nil
}

// 区块链交互
func ChannelQuery(fcn string, args [][]byte) (channel.Response, error) {
	ctx := SDK.ChannelContext(ChannelName, fabsdk.WithUser(User), fabsdk.WithOrg(Org))
	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}

	return cli.Query(channel.Request{
		ChaincodeID: ChainCodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints("word.node1.gdzc.com", "weixiong.node2.gdzc.com", "peer1.node3.gdzc.com"))
}
