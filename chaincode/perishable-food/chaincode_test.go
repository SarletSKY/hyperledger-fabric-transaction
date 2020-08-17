package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	log "github.com/wonderivan/logger"
	"strconv"
	"testing"
)

func init() {
	log.SetLogger("/home/zwx/比赛/go/src/accurchain.com/perishable-food/chaincode/config/chaincode_test_log.json")
}

/**
createOrder： 测试创建订单
参数：订单Id、商品Id、采购商、供货商、购买商品数量、采购商的地址、下单时间、运送时间
*/
func checkCreateOrderInvokeFun(stub *shim.MockStub, t *testing.T) {
	checkInvoke(stub, t, [][]byte{
		[]byte("createOrder"),
		[]byte("001"),
		[]byte("88efd7ea-bec6-4994-8ed1-f3f7b6f8cac7"),
		[]byte("1"),
		[]byte("2"),
		[]byte("5"),
		[]byte("广州市"),
		[]byte("2018-12-02T18:00:00+08:00"),
	})
	checkInvoke(stub, t, [][]byte{
		[]byte("createOrder"),
		[]byte("002"),
		[]byte("88efd7ea-bec6-4994-8ed1-f3f7b6f8cac7"),
		[]byte("1"),
		[]byte("2"),
		[]byte("5"),
		[]byte("广州市"),
		[]byte("2018-12-02T18:00:00+08:00"),
	})
}

/**
createCommodity： 测试创建商品
参数：帐号Id、商品Id、商品名称、 生产地、价格、商品介绍
*/
func checkCreateCommodityInvokeFun(stub *shim.MockStub, t *testing.T) {
	checkInvoke(stub, t, [][]byte{
		[]byte("createCommodity"),
		[]byte("2"),
		[]byte("2746eef4-7f44-4b65-a221-ca661fc0f1a2"),
		[]byte("芭蕉"),
		[]byte("中国"),
		[]byte("4"),
		[]byte("商家创建商品～芭蕉"),
	})
}

/**
createUser： 测试创建单一用户
参数：用户Id、帐号名、密码、支付密码、手机号、邮箱、企业、企业地址、 企业信用、金钱、角色
*/
func checkCreateUserInvokeFun(stub *shim.MockStub, t *testing.T) {
	checkInvoke(stub, t, [][]byte{
		[]byte("createUser"),
		[]byte("4"),
		[]byte("admin4"),
		[]byte("admin4"),
		[]byte("000004"),
		[]byte("1008611"),
		[]byte("1008611@qq.com"),
		[]byte("广东泽诚有限公司"),
		[]byte("广州"),
		[]byte("100"),
		[]byte("1000"),
		[]byte(enumRoles.FinancialInstitution),
	})
}

/**
confirmUserPwd： 测试确认帐号密码
参数：帐号Id、用户名、密码
*/
func checkConfirmUserPwdInvokeFun(stub *shim.MockStub, t *testing.T, id string, username string, password string) {
	// 调用链码就算是去测试了
	checkInvoke(stub, t, [][]byte{
		[]byte("confirmUserPwd"),
		[]byte(id),
		[]byte(username),
		[]byte(password),
	})
}

/**
queryCommodityList： 测试确认订单
参数：订单Id、角色Id、状态
*/
func checkConfirmOrderInvokeFun(stub *shim.MockStub, t *testing.T, id string, role string, status string) {
	checkInvoke(stub, t, [][]byte{
		[]byte("confirmOrder"),
		[]byte(id),
		[]byte(role),
		[]byte(status),
	})
}

/**
toApplyForPaidOrder： 测试申请代付订单
参数：订单Id、角色Id
*/
func checkToApplyForPaidOrderInvokeFunc(stub *shim.MockStub, t *testing.T, orderId string, financialInstitutionId string) {
	checkInvoke(stub, t, [][]byte{
		[]byte("toApplyForPaidOrder"),
		[]byte(orderId),
		[]byte(financialInstitutionId),
	})
}

/**
financialApproval： 金融机构审核订单栏
参数：订单Id、用户Id、角色、状态、支付密码
*/
func checkFinancialApproval(stub *shim.MockStub, t *testing.T, orderId, userId, status, payPwd string, financialApproval string) {
	checkInvoke(stub, t, [][]byte{
		[]byte("financialApproval"),
		[]byte(orderId),
		[]byte(userId),
		[]byte(status),
		[]byte(payPwd),
		[]byte(financialApproval),
	})
}

/**
updateOrderStatus： 测试更改订单状态
参数：订单Id、状态
*/
func checkUpdateOrderStatusInvokeFun(stub *shim.MockStub, t *testing.T, orderId string, status string) {
	checkInvoke(stub, t, [][]byte{
		[]byte("updateOrderStatus"),
		[]byte(orderId),
		[]byte(status),
	})
}

/**
queryCommodityList： 测试查找全部商品/附带检查
参数：nil
*/
func checkQueryCommodityInvokeFun(stub *shim.MockStub, t *testing.T, name string) {
	res := stub.MockInvoke("1", [][]byte{
		[]byte("queryCommodityList"),
	})

	// 判断调用结果状态
	if res.Status != shim.OK {
		log.Error("Query failed", string(res.Message))
		t.FailNow()
	}

	// 判断查询的结果数据是否是空？
	if res.Payload == nil {
		log.Error("Query failed to get value")
		t.FailNow()
	}

	// 将查询的结果全部返回
	var list []Commodity
	if err := json.Unmarshal(res.Payload, &list); err != nil {
		log.Error("Commodity list failed to convert from bytes")
		t.FailNow()
	}

	// 简单的判断一下
	if list[0].Name != name {
		log.Error("Query value", name, "was not as expected")
		t.FailNow()
	}
}

/**
queryOrderList： 测试查询订单列表
参数：帐号Id、订单Id
*/
// TODO: 后期修改可以传一个/两个参数，性质不一样
func checkQueryOrderInvokeFun(stub *shim.MockStub, t *testing.T, id string, userId string) {
	res := stub.MockInvoke("1", [][]byte{
		[]byte("queryOrderList"),
		[]byte(userId), // 哪个角色进行查看订单
		[]byte(id),
	})

	// 判断调用结果状态
	if res.Status != shim.OK {
		log.Error("Query failed", string(res.Message))
		t.FailNow()
	}

	// 判断查询的结果数据是否是空？
	if res.Payload == nil {
		log.Info("没有查询到任何订单")
		t.FailNow()
	}

	// 将查询的结果全部返回
	var list []Order
	if err := json.Unmarshal(res.Payload, &list); err != nil {
		log.Error("Order list failed to convert from bytes")
		t.FailNow()
	}
	// 查看内容
	for i, val := range list {
		log.Info(fmt.Sprintf("[全部]查询到的订单第%+v个结果为%+v\n", i+1, val))
	}
	// 简单的判断一下
	if list[0].Id != id {
		log.Error("Query value", id, "was not as expected")
		t.FailNow()
	}
}

/**
createOrder： 测试查询单一用户
参数：用户Id
*/
func checkQueryUserInvokeFun(stub *shim.MockStub, t *testing.T, id string) {
	res := stub.MockInvoke("1", [][]byte{
		[]byte("queryUserInfo"),
		[]byte(id),
	})

	if res.Status != shim.OK {
		log.Error("Query failed", string(res.Message))
		t.FailNow()
	}

	if res.Payload == nil {
		log.Error("Query failed to get value")
		t.FailNow()
	}

	var list User
	if err := json.Unmarshal(res.Payload, &list); err != nil {
		log.Error("User list failed to convert from bytes", err)
		t.FailNow()
	}

	// TODO: 这里值需要获取单一用户，获取所有用户暂时业务上不需要
	if list.Id != id {
		log.Error("Query value", id, "was not as expected")
		t.FailNow()
	}
	log.Info("查询到帐号%+v\n", list)
}

/**
queryApprovalOrder： 测试查询申请代付栏/审批订单栏
参数：帐号Id
*/
func checkQueryApprovalOrder(stub *shim.MockStub, t *testing.T, userId string) {
	res := stub.MockInvoke("1", [][]byte{
		[]byte("queryApprovalOrder"),
		[]byte(userId),
	})

	if res.Status != shim.OK {
		fmt.Println("Query failed", string(res.Message))
		t.FailNow()
	}

	if res.Payload == nil {
		fmt.Println("Query failed to get value")
		t.FailNow()
	}

	var list []Order
	if err := json.Unmarshal(res.Payload, &list); err != nil {
		fmt.Println("Order list failed to convert from bytes")
		t.FailNow()
	}
	if len(list) == 0 {
		log.Info(fmt.Sprintf("userId为[%+v]查询到不到结果", userId))
	}
	for i, val := range list {
		log.Info(fmt.Sprintf("userId为[%+v]查询到%+v个结果，申请代付栏/审批订单栏结果为：%+v", userId, i+1, val))
	}
}

/**
queryFinancialApproval： 测试查询金融机构列表
参数：nil
*/
func checkQueryFinancialApprovalInvokeFun(stub *shim.MockStub, t *testing.T) string {
	res := stub.MockInvoke("1", [][]byte{
		[]byte("queryFinancialApprovalList"),
	})

	// 判断调用结果状态
	if res.Status != shim.OK {
		log.Error("Query failed", string(res.Message))
		t.FailNow()
	}

	if res.Payload == nil {
		log.Error("Query failed to get value")
		t.FailNow()
	}

	var list []User
	if err := json.Unmarshal(res.Payload, &list); err != nil {
		fmt.Println("User list failed to convert from bytes")
		t.FailNow()
	}

	// 查看内容
	for i, val := range list {
		log.Info(fmt.Sprintf("查询到第%+v个结果，金融机构列表结果为：%+v", i+1, val))
	}
	return list[0].Id
}

/**
queryFinancialApproval： 测试查询金融机构列表
参数：nil
*/
func checkQueryFinishOrderInvokeFun(stub *shim.MockStub, t *testing.T, userId string) {
	res := stub.MockInvoke("1", [][]byte{
		[]byte("queryOrderCompletionList"),
		[]byte(userId),
	})

	if res.Status != shim.OK {
		log.Error("Query failed", string(res.Message))
		t.FailNow()
	}

	if res.Payload == nil {
		log.Info("没有已完成的订单")
		t.FailNow()
	}

	// 将查询的结果全部返回
	var list []Order
	if err := json.Unmarshal(res.Payload, &list); err != nil {
		log.Error("Order list failed to convert from bytes")
		t.FailNow()
	}

	// 如果反序列化之后为空
	if len(list) == 0 {
		log.Info("没有已完成的订单")
		t.FailNow()
	}
	// 查看内容
	for i, val := range list {
		log.Info(fmt.Sprintf("查询到已完成的订单第%+v个结果为%+v\n", i+1, val))
	}
}

/**
queryBuyerHistoricalTransactions： 测试查询采购商的历史交易
参数：帐号Id、状态
*/
func checkQueryBuyerHisTraInvokeFun(stub *shim.MockStub, t *testing.T, orderId string, status string) {
	res := stub.MockInvoke("1", [][]byte{
		[]byte("queryBuyerHistoricalTransactions"),
		[]byte(orderId),
		[]byte(status),
	})

	if res.Status != shim.OK {
		log.Error("Query failed", string(res.Message))
		t.FailNow()
	}

	if res.Payload == nil {
		log.Error("Query failed to get value")
		t.FailNow()
	}

	var list []Order
	if err := json.Unmarshal(res.Payload, &list); err != nil {
		log.Error("Order list failed to convert from bytes", err)
		t.FailNow()
	}
	// 查看内容
	for i, val := range list {
		log.Info(fmt.Sprintf("查询该采购商的第%v条历史交易记录%+v\n", i+1, val))
	}
}

// 检查初始化
func checkInit(stub *shim.MockStub, t *testing.T, args [][]byte) {
	log.Debug("初始化链码")
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		log.Error("Init failed", string(res.Message))
		t.FailNow()
	}
}

// 检查调用
func checkInvoke(stub *shim.MockStub, t *testing.T, args [][]byte) {
	log.Debug("开始调用invoke函数")
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		log.Error("Invoke", args, "failed", string(res.Message))
	}
}

// 检查商品获取
func checkCommodityResult(stub *shim.MockStub, t *testing.T, key string, value string) {
	keys := make([]string, 0)
	keys = append(keys, key)
	// 获取商品进行迭代
	result, err := stub.GetStateByPartialCompositeKey("commodity", keys)
	if err != nil {
		log.Error("Statue", key, "failed to get value")
	}
	defer result.Close()

	var commodityList []*Commodity
	for result.HasNext() {
		val, err := result.Next()
		if err != nil {
			t.FailNow()
		}

		commodity := new(Commodity)
		if err := json.Unmarshal(val.GetValue(), commodity); err != nil {
			log.Error("Commodity", key, "failed to convert from bytes")
			t.FailNow()
		}

		// 反序列化之化，对测试传入的value进行判断
		if commodity.Name != value {
			log.Error("Commodity name", key, "was not", value, "as expected")
			t.FailNow()
		}
		commodityList = append(commodityList, commodity)
	}
	for i, val := range commodityList {
		log.Info(fmt.Sprintf("查询到第%+v个商品信息结果为: %+v", i+1, val))
	}
}

// 检查订单获取
func checkOrderResult(stub *shim.MockStub, t *testing.T, key string, value string) {
	keys := make([]string, 0)
	keys = append(keys, key)
	result, err := stub.GetStateByPartialCompositeKey("order", keys)
	if err != nil {
		log.Error("Statue", key, "failed to get value")
		t.FailNow()
	}

	var orderList []*Order
	defer result.Close()
	for result.HasNext() {
		val, err := result.Next()
		if err != nil {
			t.FailNow()
		}

		order := new(Order)
		if err := json.Unmarshal(val.GetValue(), order); err != nil {
			log.Error("Order", key, "Failed to convert from bytes")
			t.FailNow()
		}

		if order.Id != value {
			log.Error("Order id", key, "was not", value, "as expected")
			t.FailNow()
		}
		orderList = append(orderList, order)
	}
	for i, val := range orderList {
		log.Info(fmt.Sprintf("查询到第%+v个订单信息结果为: %+v", i+1, val))
	}
}

// 检查帐号获取
func checkUserResult(stub *shim.MockStub, t *testing.T, key string, value string) {
	keys := make([]string, 0)
	keys = append(keys, key)
	result, err := stub.GetStateByPartialCompositeKey("user", keys)
	if err != nil {
		log.Error("Statue", key, "failed to get value")
		t.FailNow()
	}
	defer result.Close()
	var userList []*User
	for result.HasNext() {
		val, err := result.Next()
		if err != nil {
			t.FailNow()
		}

		user := new(User)
		if err := json.Unmarshal(val.GetValue(), user); err != nil {
			log.Error("User", key, "failed to convert from bytes")
			t.FailNow()
		}

		// 反序列化之化，对测试传入的value进行判断
		if user.Role != value {
			log.Error("User name", key, "was not", value, "as expected")
			t.FailNow()
		}
		userList = append(userList, user)
	}

	for i, val := range userList {
		log.Info(fmt.Sprintf("查询到第%+v个帐号信息结果为: %+v", i+1, val))
	}
}

// 检查订单状态获取
func checkOrderStatusResult(stub *shim.MockStub, t *testing.T, key string, value string) {
	keys := make([]string, 0)
	keys = append(keys, key)
	result, err := stub.GetStateByPartialCompositeKey("order", keys)
	if err != nil {
		log.Error("Statue", key, "failed to get value")
		t.FailNow()
	}

	defer result.Close()
	for result.HasNext() {
		val, err := result.Next()
		if err != nil {
			t.FailNow()
		}

		order := new(Order)
		if err := json.Unmarshal(val.GetValue(), order); err != nil {
			log.Error("Order", key, "Failed to convert from bytes")
			t.FailNow()
		}
		log.Info(fmt.Sprintf("查询到订单状态信息结果为: %+v", order))
		if order.Status != value {
			log.Error("Order id", key, "was not", value, "as expected")
			t.FailNow()
		}
	}
}

// 测试初始化、获取帐号/商品
func TestSmartContract_Init(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试初始化、获取帐号/商品", scc)

	checkInit(stub, t, nil)

	var roleList = [3]string{enumRoles.Buyer, enumRoles.Seller, enumRoles.FinancialInstitution}

	log.Debug("检查帐号")
	// 检查帐号列表
	for i, val := range roleList {
		checkUserResult(stub, t, strconv.Itoa(i+1), val)
	}

	var commodityName = [3]string{"国光", "红星", "红富士"}
	var ids = [3]string{
		"88efd7ea-bec6-4994-8ed1-f3f7b6f8cac7",
		"36bf5c7f-4cf7-4926-b0f6-0c5c18515752",
		"d9ce807b-e308-11e8-a47c-3e1591a6f5bb",
	}

	// 检查商品列表
	log.Debug("检查商品")
	for i, val := range commodityName {
		checkCommodityResult(stub, t, ids[i], val)
	}
}

// 测试创建商品
func TestCreateCommodity(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试创建商品", scc)

	checkInit(stub, t, nil)

	checkCreateCommodityInvokeFun(stub, t)

	// 检查创建商品是否添加成功
	checkCommodityResult(stub, t, "2746eef4-7f44-4b65-a221-ca661fc0f1a2", "芭蕉")
}

// 测试创建订单
func TestCreateOrder(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试创建订单", scc)

	checkInit(stub, t, nil)

	// 检查Invoke创建订单函数
	checkCreateOrderInvokeFun(stub, t)

	// 检查创建订单是否添加成功
	checkOrderResult(stub, t, "001", "001")
}

// 注册帐号
func TestCreateUser(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试注册帐号", scc)

	checkInit(stub, t, nil)

	checkCreateUserInvokeFun(stub, t)
	checkQueryUserInvokeFun(stub, t, "4")
}

// 测试查找商品
func TestQueryCommodity(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试检查商品", scc)

	checkInit(stub, t, nil)
	checkCreateCommodityInvokeFun(stub, t)

	checkQueryCommodityInvokeFun(stub, t, "芭蕉")
}

// 测试查找订单
func TestQueryOrder(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试查找订单", scc)

	checkInit(stub, t, nil)
	checkCreateOrderInvokeFun(stub, t)
	// 1. 不传参，就是简简单单的测试上一个创建订单有没有成功
	// 2. 传参来判断不同角色显示一不一致
	// role: 1～采购商、2～供货商、3～金融机构

	// TODO: 这里后续增加实现更改状态[当然直接修改不现实，正常是前端有严格要求！,在之后如果有空，将它改成从头到尾正常的修改状态],然后在查状态，这里我先手动新建订单不是新建(未确认)就行测试[功能是成功的]
	checkQueryOrderInvokeFun(stub, t, "002", "1")
	/*checkQueryOrderInvokeFun(stub, t, "001", "2")
	checkQueryOrderInvokeFun(stub, t, "001", "3")*/
}

// 测试用户查询
func TestUserQuery(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试用户查询", scc)

	checkInit(stub, t, nil)

	// TODO: 注册用户还没做
	checkQueryUserInvokeFun(stub, t, "4")
}

// 测试帐号密码正确性
func TestConfirmUserPwd(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试帐号密码正确性", scc)

	checkInit(stub, t, nil)

	checkConfirmUserPwdInvokeFun(stub, t, "1", "admin1", "admin1")
	checkConfirmUserPwdInvokeFun(stub, t, "2", "admin2", "admin2")
	checkConfirmUserPwdInvokeFun(stub, t, "3", "admin3", "admin3")
}

// 测试确认订单
func TestConfirmOrder(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试确认订单", scc)

	checkInit(stub, t, nil)

	checkCreateOrderInvokeFun(stub, t)

	// 检查确认订单
	checkConfirmOrderInvokeFun(stub, t, "001", enumRoles.Seller, enumStatus.NewConfirmed)

	// 查询订单状态结果 key~orderId value~status
	checkOrderStatusResult(stub, t, "001", enumStatus.NewConfirmed)
}

// 测试查询金融机构列表
func TestQueryFinancialApprovalList(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试查询金融机构列表", scc)

	checkInit(stub, t, nil)

	checkQueryFinancialApprovalInvokeFun(stub, t)
}

// 测试申请代付订单
func TestToApplyForPaidOrder(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试申请代付订单", scc)

	checkInit(stub, t, nil)

	checkCreateOrderInvokeFun(stub, t)
	checkQueryOrderInvokeFun(stub, t, "001", "1")
	checkConfirmOrderInvokeFun(stub, t, "001", enumRoles.Seller, enumStatus.NewConfirmed)
	checkQueryOrderInvokeFun(stub, t, "001", "1")
	// 这里接收金融机构列表的Id
	financialInstitutionId := checkQueryFinancialApprovalInvokeFun(stub, t)

	checkToApplyForPaidOrderInvokeFunc(stub, t, "001", financialInstitutionId)
	// 检查结果
	checkOrderStatusResult(stub, t, "001", enumStatus.ToApplyForPaid)
}

// 测试申请代付栏/审批订单栏
func TestQueryApprovalOrder(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试查询申请代付栏/审批订单栏", scc)

	checkInit(stub, t, nil)

	checkCreateOrderInvokeFun(stub, t)
	checkConfirmOrderInvokeFun(stub, t, "001", enumRoles.Seller, enumStatus.NewConfirmed)
	checkQueryApprovalOrder(stub, t, "2")
	financialInstitutionId := checkQueryFinancialApprovalInvokeFun(stub, t)
	checkToApplyForPaidOrderInvokeFunc(stub, t, "001", financialInstitutionId)

	checkQueryApprovalOrder(stub, t, "1")
	checkQueryApprovalOrder(stub, t, "3")
}

// 测试更新状态
func TestUpdateOrderStatus(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试查询申请代付栏/审批订单栏", scc)

	checkInit(stub, t, nil)

	checkCreateOrderInvokeFun(stub, t)
	checkConfirmOrderInvokeFun(stub, t, "001", enumRoles.Seller, enumStatus.NewConfirmed)
	financialInstitutionId := checkQueryFinancialApprovalInvokeFun(stub, t)
	checkToApplyForPaidOrderInvokeFunc(stub, t, "001", financialInstitutionId)
	checkFinancialApproval(stub, t, "001", "3", enumStatus.PaymentNoVerified, "000003", "2018-12-02T18:00:00+08:00")

	checkUpdateOrderStatusInvokeFun(stub, t, "001", enumStatus.PaymentVerified)
	checkOrderStatusResult(stub, t, "001", enumStatus.PaymentVerified)
}

// 测试审核订单状态
func TestFinancialApproval(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试查询申请代付栏/审批订单栏", scc)

	checkInit(stub, t, nil)

	checkCreateOrderInvokeFun(stub, t)
	checkConfirmOrderInvokeFun(stub, t, "001", enumRoles.Seller, enumStatus.NewConfirmed)
	financialInstitutionId := checkQueryFinancialApprovalInvokeFun(stub, t)
	checkToApplyForPaidOrderInvokeFunc(stub, t, "001", financialInstitutionId)

	checkFinancialApproval(stub, t, "001", "3", enumStatus.PaymentNoVerified, "000003", "2018-12-02T18:00:00+08:00")
	checkQueryOrderInvokeFun(stub, t, "001", "3")
}

// 测试查询已完成的订单
func TestQueryOrderCompletionList(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试查询申请代付栏/审批订单栏", scc)

	checkInit(stub, t, nil)

	checkCreateOrderInvokeFun(stub, t)
	checkConfirmOrderInvokeFun(stub, t, "001", enumRoles.Seller, enumStatus.NewConfirmed)
	financialInstitutionId := checkQueryFinancialApprovalInvokeFun(stub, t)
	checkToApplyForPaidOrderInvokeFunc(stub, t, "001", financialInstitutionId)
	checkFinancialApproval(stub, t, "001", "3", enumStatus.PaymentNoVerified, "000003", "2018-12-02T18:00:00+08:00")
	// 这里直接修改成最后状态
	checkUpdateOrderStatusInvokeFun(stub, t, "001", enumStatus.DeliveredVerified)
	checkQueryFinishOrderInvokeFun(stub, t, "2")
	checkCreateUserInvokeFun(stub, t)
	checkQueryUserInvokeFun(stub, t, "4")
}

// 测试查询采购商的历史交易
func TestQueryBuyerHistoricalTransactions(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("测试查询采购商的历史交易", scc)

	checkInit(stub, t, nil)

	checkCreateOrderInvokeFun(stub, t)
	checkConfirmOrderInvokeFun(stub, t, "001", enumRoles.Seller, enumStatus.NewConfirmed)
	financialInstitutionId := checkQueryFinancialApprovalInvokeFun(stub, t)
	checkToApplyForPaidOrderInvokeFunc(stub, t, "001", financialInstitutionId)

	// 根据金融机构查看某订单两个选项，来决定那个状态[就两个状态]严格来说是三个
	checkQueryBuyerHisTraInvokeFun(stub, t, "001", enumStatus.DeliveredVerified)
	checkQueryBuyerHisTraInvokeFun(stub, t, "001", enumStatus.RefundVerified)
	checkQueryBuyerHisTraInvokeFun(stub, t, "001", enumStatus.RefundNoVerified)
}

// 总体操作一遍
func TestSmartContract_Invoke(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("总体测试整个流程", scc)

	checkInit(stub, t, nil)

	checkCreateOrderInvokeFun(stub, t)
	checkQueryOrderInvokeFun(stub, t, "001", "1")

	checkConfirmOrderInvokeFun(stub, t, "001", enumRoles.Seller, enumStatus.NewConfirmed)
	checkQueryOrderInvokeFun(stub, t, "001", "1")

	financialInstitutionId := checkQueryFinancialApprovalInvokeFun(stub, t)
	checkToApplyForPaidOrderInvokeFunc(stub, t, "001", financialInstitutionId)
	checkQueryOrderInvokeFun(stub, t, "001", "1")

	checkFinancialApproval(stub, t, "001", "3", enumStatus.PaymentNoVerified, "000003", "2018-12-02T18:00:00+08:00")
	checkQueryOrderInvokeFun(stub, t, "001", "1")

	checkUpdateOrderStatusInvokeFun(stub, t, "001", enumStatus.PaymentVerified)
	checkQueryOrderInvokeFun(stub, t, "001", "1")

	checkUpdateOrderStatusInvokeFun(stub, t, "001", enumStatus.DeliveredNoVerified)
	checkQueryOrderInvokeFun(stub, t, "001", "1")

	checkUpdateOrderStatusInvokeFun(stub, t, "001", enumStatus.DeliveredVerified)
	checkQueryOrderInvokeFun(stub, t, "001", "1")

	checkUpdateOrderStatusInvokeFun(stub, t, "001", enumStatus.RefundNoVerified)
	checkQueryOrderInvokeFun(stub, t, "001", "1")

	checkUpdateOrderStatusInvokeFun(stub, t, "001", enumStatus.RefundVerified)
	checkQueryFinishOrderInvokeFun(stub, t, "1")
}
