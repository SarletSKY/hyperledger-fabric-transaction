package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
)

// TODO: 查看信用，判断是否已经过了还款时间，之后在做。
const (
	timeLayout = "2006-01-02T15:04:05Z07:00"
)

type SmartContract struct{}

var (
	enumStatus                       = newStatus()
	enumRoles                        = newRoles()
	ApprovalOrdersColumnByFiAndBuyer = []string{enumStatus.ToApplyForPaid, enumStatus.UnPayment, enumStatus.PaymentNoVerified}        //审批订单[采购商]
	ApprovalOrdersColumnBySeller     = []string{enumStatus.NewConfirmed, enumStatus.NewUnconfirmed}                                   // 审批订单[供货商/金融机构]
	QueryOrderBySeller               = []string{enumStatus.RefundVerified, enumStatus.RefundNoVerified}                               // 查询全部订单
	SellerFinishOrder                = []string{enumStatus.DeliveredVerified, enumStatus.RefundVerified, enumStatus.RefundNoVerified} // 供货商已完成订单状态
	UpdateStatusRangeList            = []string{                                                                                      // 更新状态
		enumStatus.PaymentVerified,
		enumStatus.DeliveredNoVerified,
		enumStatus.DeliveredVerified,
		enumStatus.RefundVerified,
		enumStatus.RefundNoVerified,
	}
)

// 角色枚举
func newRoles() *Role {
	return &Role{
		Buyer:                "采购商",
		Seller:               "供货商",
		FinancialInstitution: "金融机构",
	}
}

// 状态枚举
func newStatus() *Status {
	return &Status{
		NewUnconfirmed: "新建(待确认)",
		NewConfirmed:   "新建(已确认)",

		RefundVerified:    "已还款(已核实)",
		RefundNoVerified:  "已还款(未核实)",
		DeliveredVerified: "已发货(已核实)",

		ToApplyForPaid:      "申请代付",
		UnPayment:           "代付被拒",
		PaymentNoVerified:   "已付款(未核实)",
		PaymentVerified:     "已付款(已核实)",
		DeliveredNoVerified: "已发货(未核实)",
	}
}

// TODO: 还要功能没做，每条继续看pdf
// Init函数
func (t *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// 定义变量名
	var commodityName = [3]string{"国光", "红星", "红富士"}
	// 帐号列表
	var userList []*User
	var roleList = [3]string{
		enumRoles.Buyer,
		enumRoles.Seller,
		enumRoles.FinancialInstitution,
	}
	// 创建三个帐号,并初始化三个角色
	for i, val := range roleList {
		id := strconv.Itoa(i + 1)
		user := User{
			Id:               id,
			UserName:         "admin" + id,
			Password:         "admin" + id,
			PayPassword:      "00000" + id,
			Phone:            "13790822374",
			Email:            "10086@qq.com",
			Enterprise:       "广东泽诚有限公司",
			EnterpriseOrigin: "广州",
			EnterpriseCredit: 100,
			Balance:          1000,
			Role:             val,
		}
		// 序列化
		userBytes, err := json.Marshal(user)
		if err != nil {
			return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
		}
		// 创建复合键
		var key string
		if val, err := stub.CreateCompositeKey("user", []string{user.Id}); err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("create key error %s", err), Payload: nil}
		} else {
			key = val
		}
		// 存入账本
		if err = stub.PutState(key, userBytes); err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("put statue failed: %s", err), Payload: nil}
		}
		// 存进帐号列表
		userList = append(userList, &user)
	}

	var ids = [3]string{
		"88efd7ea-bec6-4994-8ed1-f3f7b6f8cac7",
		"36bf5c7f-4cf7-4926-b0f6-0c5c18515752",
		"d9ce807b-e308-11e8-a47c-3e1591a6f5bb",
	}
	// 为卖家创建商品
	for i, val := range commodityName {
		// 商品使用uuid生成
		//id := uuid.Must(uuid.NewV4()).String()
		price := 6.00 + float64(i)
		//fmt.Println("id", id)
		commodity := Commodity{
			User:      userList[1],
			Id:        ids[i],
			Name:      val,
			Origin:    "中国",
			Price:     price,
			Introduce: val + "这个商品属于" + userList[1].Id + "用户",
		}
		// 序列化
		commodityBytes, err := json.Marshal(commodity)
		if err != nil {
			return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
		}
		// 创建复合键
		var key string
		if val, err := stub.CreateCompositeKey("commodity", []string{commodity.Id}); err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("create key error %s", err), Payload: nil}
		} else {
			key = val
		}
		// 存入账本
		if err = stub.PutState(key, commodityBytes); err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("put statue failed: %s", err), Payload: nil}
		}
	}
	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: nil,
	}
}

// Invoke函数
func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "createCommodity" { // 创建商品
		return t.createCommodity(stub, args)
	} else if function == "queryCommodityList" { // 查询商品列表
		return t.queryCommodityList(stub, args)
	} else if function == "createOrder" { // 创建订单列表
		return t.createOrder(stub, args)
	} else if function == "queryOrderList" { // 查询全部订单列表[全部订单栏]
		return t.queryOrderList(stub, args)
	} else if function == "queryOrderCompletionList" { // 查看已完成的订单
		return t.queryOrderCompletionList(stub, args)
	} else if function == "queryApprovalOrder" { // [供货商/金融机构/采购商]查询审批订单/申请代付订单[审批订单栏/申请代付栏]
		return t.queryApprovalOrder(stub, args)
	} else if function == "queryFinancialApprovalList" { // [采购商](获取金融机构列表)  当采购商点击新建订单(已确认)之后就改状态，并且进入申请代付操作[就是这个函数了]
		return t.queryFinancialApprovalList(stub, args)
	} else if function == "queryUserInfo" { // 查询用户信息
		return t.queryUserInfo(stub, args)
	} else if function == "confirmOrder" { // 供货商的审核订单[确认订单]
		return t.confirmOrder(stub, args)
	} else if function == "toApplyForPaidOrder" { // 申请代付订单
		return t.toApplyForPaidOrder(stub, args)
	} else if function == "financialApproval" { // 金融机构审核订单[代付]
		return t.financialApproval(stub, args)
	} else if function == "confirmUserPwd" { // 确认帐号密码
		return t.confirmUserPwd(stub, args)
	} else if function == "updateOrderStatus" { // 更改订单状态
		return t.updateOrderStatus(stub, args)
	} else if function == "createUser" { // 注册帐号
		return t.createUser(stub, args)
	} else if function == "queryBuyerHistoricalTransactions" {
		return t.queryBuyerHistoricalTransactions(stub, args)
	} else {
		return peer.Response{Status: 500, Message: "函数不对", Payload: nil}
	}
}

// 创建商品 ✔
func (t *SmartContract) createCommodity(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 判断参数是否为空
	if len(args) != 6 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}
	// 设置变量
	userId := args[0]
	// TODO: 这里的id可能要自己处理成uuid
	id := args[1]
	name := args[2]
	origin := args[3]
	price := args[4]
	introduce := args[5]

	if userId == "" || id == "" || name == "" || origin == "" || price == "" || introduce == "" {
		return peer.Response{Status: 5002, Message: "参数不能为空", Payload: nil}
	}

	// 获取单个用户信息
	resp, user := getUserInfo(stub, []string{userId})
	if resp.Status != int32(shim.OK) {
		return resp
	}

	// 判断该帐号是否为供货商，不是不同意创建
	if user.Role != enumRoles.Seller {
		return peer.Response{Status: 500, Message: "无法创建商品", Payload: nil}
	}

	// 数据格式转换
	var formattedPrice float64
	if val, err := strconv.ParseFloat(price, 64); err != nil {
		return peer.Response{Status: 500, Message: "format Price error", Payload: nil}
	} else {
		formattedPrice = val
	}

	// 创建商品对象
	commodity := &Commodity{
		User:      user,
		Id:        id,
		Name:      name,
		Origin:    origin,
		Price:     formattedPrice,
		Introduce: introduce,
	}

	// 创建主键
	var key string
	if val, err := stub.CreateCompositeKey("commodity", []string{id}); err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("create key error %s", err), Payload: nil}
	} else {
		key = val
	}

	// 判断有没有该商品
	if commodityBytes, err := stub.GetState(key); err == nil && len(commodityBytes) != 0 {
		return peer.Response{Status: 500, Message: "commodity already exist", Payload: nil}
	}

	// 序列化
	commodityBytes, err := json.Marshal(commodity)
	if err != nil {
		return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
	}

	// 存入账本
	if err = stub.PutState(key, commodityBytes); err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("put statue failed: %s", err), Payload: nil}
	}
	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: nil,
	}
}

// 查询商品 ✔
func (t *SmartContract) queryCommodityList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 0 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}

	// 获取key
	keys := make([]string, 0)
	result, err := stub.GetStateByPartialCompositeKey("commodity", keys)
	if err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("Get commodity error %s", err), Payload: nil}
	}
	defer result.Close()

	commodityList := make([]*Commodity, 0)
	for result.HasNext() {
		val, err := result.Next()
		if err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("Get commodity error %s", err), Payload: nil}
		}
		commodity := new(Commodity)
		if err := json.Unmarshal(val.GetValue(), commodity); err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("Commodity failed to convert from bytes, error %s", err), Payload: nil}
		}
		commodityList = append(commodityList, commodity)
	}

	// 序列化之后在返回前端
	commodityBytes, err := json.Marshal(commodityList)
	if err != nil {
		return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
	}

	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: commodityBytes,
	}
}

// 创建订单 ✔
func (t *SmartContract) createOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 7 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}
	id := args[0]
	commodityId := args[1]
	buyer := args[2]
	seller := args[3]
	quantity := args[4]
	deliverAddress := args[5]
	orderTime := args[6]

	if id == "" || commodityId == "" || buyer == "" || seller == "" || quantity == "" || deliverAddress == "" || orderTime == "" {
		return peer.Response{Status: 5002, Message: "参数不能为空", Payload: nil}
	}

	// 判断商品存不存在
	result, err := stub.GetStateByPartialCompositeKey("commodity", []string{commodityId})
	if err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("Get commodity error %s", err), Payload: nil}
	}
	defer result.Close()

	commodity := new(Commodity)
	for result.HasNext() {
		val, err := result.Next()
		if err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("Get commodity error %s", err), Payload: nil}
		}

		if err := json.Unmarshal(val.GetValue(), commodity); err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("commodity failed to convert from bytes, error %s", err), Payload: nil}
		}
	}

	// 格式转换
	var formattedQuantity float64
	if val, err := strconv.ParseFloat(quantity, 64); err != nil {
		return peer.Response{Status: 500, Message: "format Price error", Payload: nil}
	} else {
		formattedQuantity = val
	}

	var formattedOrderTime time.Time
	if val, err := time.Parse(timeLayout, orderTime); err != nil {
		return peer.Response{Status: 500, Message: "format OrderTime error", Payload: nil}
	} else {
		formattedOrderTime = val
	}

	//创建order对象
	order := Order{
		Commodity:        commodity,
		Id:               id,
		Buyer:            buyer,
		Seller:           seller, // TODO: 这里可以换掉commodity.User.Id
		Quantity:         formattedQuantity,
		OrderTotalAmount: formattedQuantity * commodity.Price,
		DeliverAddress:   deliverAddress,
		Status:           enumStatus.NewUnconfirmed,
		OrderTime:        formattedOrderTime,
	}

	// 序列化对象
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
	}

	//创建主键
	var key string
	if val, err := stub.CreateCompositeKey("order", []string{id}); err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("create key error %s", err), Payload: nil}
	} else {
		key = val
	}

	//写入区块链账本
	if err := stub.PutState(key, orderBytes); err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("put statue failed: %s", err), Payload: nil}
	}
	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: nil,
	}
}

// 查询单一订单/全部订单 ✔/✔
func (t *SmartContract) queryOrderList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) > 2 || len(args) < 1 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}
	userId := args[0]

	if userId == "" {
		return peer.Response{Status: 5002, Message: "参数不能为空", Payload: nil}
	}

	keys := make([]string, 0)
	if len(args) == 2 {
		orderId := args[1]
		if orderId == "" {
			return peer.Response{Status: 5002, Message: "参数不能为空", Payload: nil}
		}
		keys = append(keys, orderId)
	}

	// 获取单个帐号信息
	resp, user := getUserInfo(stub, []string{userId})
	if resp.Status != int32(shim.OK) {
		return resp
	}

	// role主要为了区分显示全部订单列表
	orderResult, err := stub.GetStateByPartialCompositeKey("order", keys)
	if err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("query order error: %s", err), Payload: nil}
	}
	defer orderResult.Close()

	orderList := make([]*Order, 0)
	for orderResult.HasNext() {
		val, err := orderResult.Next()
		if err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("Get order error %s", err), Payload: nil}
		}

		order := new(Order)
		if err := json.Unmarshal(val.GetValue(), order); err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("Order failed to convert from bytes, error %s", err), Payload: nil}
		}

		// 是否存在某帐号
		if !IsExistItem(user.Id, []string{order.Buyer, order.Seller, order.FinancialInstitution}) {
			continue
		}

		// 判断状态，区分下供货商。 如果帐号为供货商
		if user.Role == enumRoles.Seller {
			// 如果该订单的状态为已还款状态，则改成已发货状态返回给前端
			if IsExistItem(order.Status, QueryOrderBySeller) {
				order.Status = enumStatus.DeliveredVerified
			}
		}
		orderList = append(orderList, order)
	}

	// 序列化返回
	orderBytes, err := json.Marshal(orderList)
	if err != nil {
		return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
	}

	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: orderBytes,
	}
}

// 查询单一帐号 ✔
func (t *SmartContract) queryUserInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}

	userId := args[0]
	if userId == "" {
		return peer.Response{Status: 5002, Message: "参数不能为空", Payload: nil}
	}

	resp, user := getUserInfo(stub, []string{userId})
	if resp.Status != int32(shim.OK) {
		return resp
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
	}

	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: userBytes,
	}
}

// [供货商]审核订单栏[确认订单] ✔
func (t *SmartContract) confirmOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}

	orderId := args[0]
	role := args[1]
	status := args[2]

	fmt.Println("结果", orderId, role, status)
	if orderId == "" || status == "" || role == "" {
		return peer.Response{Status: 5002, Message: fmt.Sprintf("参数不能为空: %v %v %v", orderId, role, status), Payload: nil}
	}

	// 判断帐号是否为供货商
	if role != enumRoles.Seller {
		return peer.Response{Status: 500, Message: "该用户更改订单状态失败", Payload: nil}
	}

	resp, order := getOrderInfo(stub, []string{orderId})
	if resp.Status != int32(shim.OK) {
		return resp
	}
	// 后期可以删除
	if !IsExistItem(status, []string{enumStatus.NewConfirmed}) {
		return peer.Response{Status: 500, Message: "操作失败", Payload: nil}
	}

	order.Status = status

	// 序列化
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
	}

	var key string
	if val, err := stub.CreateCompositeKey("order", []string{orderId}); err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("create key error %s", err), Payload: nil}
	} else {
		key = val
	}
	if err = stub.PutState(key, orderBytes); err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("put statue failed: %s", err), Payload: nil}
	}

	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: nil,
	}
}

// [金融机构]审核订单栏[审核订单] ✔
func (t *SmartContract) financialApproval(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}

	orderId := args[0]
	userId := args[1]
	status := args[2]
	payPwd := args[3]
	promiseRepaymentTime := args[4]

	if orderId == "" || status == "" || payPwd == "" || userId == "" || promiseRepaymentTime == "" {
		return peer.Response{Status: 5002, Message: "参数不能为空", Payload: nil}
	}

	// 数据格式转换
	var formattedPromiseRepaymentTime time.Time
	if val, err := time.Parse(timeLayout, promiseRepaymentTime); err != nil {
		return peer.Response{Status: 500, Message: "format formattedPromiseRepaymentTime error", Payload: nil}
	} else {
		formattedPromiseRepaymentTime = val
	}

	// 获取单个订单
	resp, order := getOrderInfo(stub, []string{orderId})
	if resp.Status != int32(shim.OK) {
		return resp
	}

	// 判断金融机构是否接受贷款
	if status == enumStatus.UnPayment {
		// 将代款改成代付被拒 TODO:之后可以优化采购商再次申请，而金融机构可以增加拉黑功能
		order.Status = enumStatus.UnPayment
	} else if status == enumStatus.PaymentNoVerified {
		// 同意代付,填写支付密码,通过则生成生成付款单号
		resp, user := getUserInfo(stub, []string{userId})
		if resp.Status != int32(shim.OK) {
			return resp
		}
		// 判断用户是否为金融机构
		if user.Role != enumRoles.FinancialInstitution {
			return peer.Response{Status: 500, Message: "操作失败", Payload: nil}
		}
		if user.PayPassword != payPwd {
			return peer.Response{Status: 500, Message: "支付密码输入不正确", Payload: nil}
		}
		// 生成付款单号 = 时间戳+订单id+金融帐号 [sha256加密]
		now := time.Now().Format(timeLayout)
		PayOrderNo := Sha256([]byte(now + orderId + user.UserName))
		order.PayOrderNo = PayOrderNo
		order.Status = enumStatus.PaymentNoVerified
		order.PromiseRepaymentTime = formattedPromiseRepaymentTime
	} else {
		return peer.Response{Status: 500, Message: "状态不正确", Payload: nil}
	}

	// 序列化
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
	}

	var key string
	if val, err := stub.CreateCompositeKey("order", []string{orderId}); err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("create key error %s", err), Payload: nil}
	} else {
		key = val
	}
	if err = stub.PutState(key, orderBytes); err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("put statue failed: %s", err), Payload: nil}
	}

	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: nil,
	}
}

// 确认帐号密码 ✔
func (t *SmartContract) confirmUserPwd(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}

	userId := args[0]
	username := args[1]
	pwd := args[2]

	if userId == "" || username == "" || pwd == "" {
		return peer.Response{Status: 5002, Message: "参数不能为空", Payload: nil}
	}

	resp, user := getUserInfo(stub, []string{userId})
	if resp.Status != int32(shim.OK) {
		return resp
	}

	// 判断当前用户与输入的帐号密码一致性 TODO: 密码加密
	if username != user.UserName || pwd != user.Password {
		return peer.Response{Status: 500, Message: "用户名或者密码错误", Payload: nil}
	}

	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: nil,
	}
}

// 查询审批订单栏/申请代付栏 ✔
func (t *SmartContract) queryApprovalOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}

	userId := args[0]

	if userId == "" {
		return peer.Response{Status: 5002, Message: "参数不能为空", Payload: nil}
	}

	// 获取帐号信息
	userResp, user := getUserInfo(stub, []string{userId})
	if userResp.Status != shim.OK {
		return userResp
	}

	var keys []string
	// 获取全部订单信息
	orderResult, err := stub.GetStateByPartialCompositeKey("order", keys)
	if err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("query order error: %s", err), Payload: nil}
	}
	defer orderResult.Close()

	// [供货商/金融机构]列表
	sellerOrderList := make([]*Order, 0)
	fiAndBuyerOrderList := make([]*Order, 0)
	for orderResult.HasNext() {
		val, err := orderResult.Next()
		if err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("Get order error %s", err), Payload: nil}
		}

		order := new(Order)
		if err := json.Unmarshal(val.GetValue(), order); err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("Order failed to convert from bytes, error %s", err), Payload: nil}
		}

		// 如果根改帐号无关的不进行显示在页面[过滤]
		if !IsExistItem(user.Id, []string{order.Buyer, order.Seller, order.FinancialInstitution}) {
			continue
		}

		// 只选择对应的状态
		if IsExistItem(order.Status, ApprovalOrdersColumnBySeller) { // 供货商审批栏两个状态，因为新建已确认之后，还要采购商点击确认订单状态才会改变成申请审批
			sellerOrderList = append(sellerOrderList, order)
		} else if IsExistItem(order.Status, ApprovalOrdersColumnByFiAndBuyer) {
			fiAndBuyerOrderList = append(fiAndBuyerOrderList, order)
		}
	}

	var orderBytes []byte
	// 判断用户的角色是[供货商/金融机构/]
	if user.Role == enumRoles.Seller {
		orderBytes, _ = json.Marshal(sellerOrderList)
	} else if user.Role == enumRoles.FinancialInstitution || user.Role == enumRoles.Buyer {
		orderBytes, _ = json.Marshal(fiAndBuyerOrderList)
	}
	// 返回数据
	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: orderBytes,
	}
}

// 查询金融机构列表[获取金融机构列表] ✔
func (t *SmartContract) queryFinancialApprovalList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 0 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}

	// 获取所有用户
	var keys []string
	userResult, err := stub.GetStateByPartialCompositeKey("user", keys)
	if err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("query user error: %s", err), Payload: nil}
	}
	defer userResult.Close()

	userList := make([]*User, 0)
	for userResult.HasNext() {
		val, err := userResult.Next()
		if err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("Get user error %s", err), Payload: nil}
		}

		user := new(User)
		if err := json.Unmarshal(val.GetValue(), user); err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("User failed to convert from bytes, error %s", err), Payload: nil}
		}
		if user.Role == enumRoles.FinancialInstitution {
			userList = append(userList, user)
		}
	}
	// 序列化
	marshal, err := json.Marshal(userList)
	if err != nil {
		return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
	}

	// 返回数据
	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: marshal,
	}
}

// [采购商]申请代付订单[给订单添加金融机构] ✔
func (t *SmartContract) toApplyForPaidOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}
	orderId := args[0]
	userId := args[1]

	if orderId == "" || userId == "" {
		return peer.Response{Status: 5002, Message: "参数不能为空", Payload: nil}
	}

	// 获取单个订单
	orderResp, order := getOrderInfo(stub, []string{orderId})
	if orderResp.Status != int32(shim.OK) {
		return orderResp
	}

	if order.Status != enumStatus.NewConfirmed {
		return peer.Response{Status: 500, Message: "该状态不正确", Payload: nil}
	}

	// 获取金融机构帐号信息
	userResp, user := getUserInfo(stub, []string{userId})
	if userResp.Status != int32(shim.OK) {
		return userResp
	}

	// 补充：如果获取的不是金融机构帐号，直接报错[当然一般不会这样，因为前端已经是获取果金融机构列表的]
	if user.Role != enumRoles.FinancialInstitution {
		return peer.Response{Status: 500, Message: "数据出错", Payload: nil}
	}

	// 修改订单状态为申请代付
	order.Status = enumStatus.ToApplyForPaid
	// 给订单添加金融机构id
	order.FinancialInstitution = user.Id

	// 创建复合键
	var key string
	if val, err := stub.CreateCompositeKey("order", []string{orderId}); err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("create key error %s", err), Payload: nil}
	} else {
		key = val
	}

	// 存入账本
	marshal, err := json.Marshal(order)
	if err != nil {
		return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
	}

	if err = stub.PutState(key, marshal); err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("put statue failed: %s", err), Payload: nil}
	}

	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: nil,
	}
}

// 创建帐号 ✔
func (t *SmartContract) createUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 11 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}

	id := args[0]
	username := args[1]
	password := args[2]
	payPassword := args[3]
	phone := args[4]
	email := args[5]
	enterprise := args[6]
	enterpriseOrigin := args[7]
	enterpriseCredit := args[8]
	balance := args[9]
	role := args[10]

	if id == "" || username == "" || password == "" || payPassword == "" || phone == "" || email == "" || enterprise == "" || enterpriseOrigin == "" || enterpriseCredit == "" || balance == "" || role == "" {
		return peer.Response{Status: 5002, Message: "参数不能为空", Payload: nil}
	}

	// 支付密码一定要6位
	if len(payPassword) != 6 {
		return peer.Response{Status: 500, Message: "支付密码不正确", Payload: nil}
	}

	// 数据格式转换
	var formattedEnterpriseCredit int64
	if val, err := strconv.ParseInt(enterpriseCredit, 10, 64); err != nil {
		return peer.Response{Status: 500, Message: "format enterpriseCredit error", Payload: nil}
	} else {
		formattedEnterpriseCredit = val
	}

	var formattedBalance float64
	if val, err := strconv.ParseFloat(balance, 64); err != nil {
		return peer.Response{Status: 500, Message: "format balance error", Payload: nil}
	} else {
		formattedBalance = val
	}

	user := &User{
		Id:               id,
		UserName:         username,
		Password:         password,
		PayPassword:      payPassword,
		Phone:            phone,
		Email:            email,
		Enterprise:       enterprise,
		EnterpriseOrigin: enterpriseOrigin,
		EnterpriseCredit: formattedEnterpriseCredit,
		Balance:          formattedBalance,
		Role:             role,
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
	}

	var key string
	if val, err := stub.CreateCompositeKey("user", []string{user.Id}); err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("create key error %s", err), Payload: nil}
	} else {
		key = val
	}

	//判断该用户是否已经被注册
	if bytes, err := stub.GetState(key); err == nil && len(bytes) != 0 {
		return peer.Response{Status: 500, Message: fmt.Sprintf("create key error %s", err), Payload: nil}
	}

	if err = stub.PutState(key, userBytes); err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("put statue failed: %s", err), Payload: nil}
	}

	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: nil,
	}
}

// 更新订单状态 ✔
func (t *SmartContract) updateOrderStatus(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}

	orderId := args[0]
	status := args[1]

	var formattedCurTime time.Time
	if val, err := time.Parse(timeLayout, time.Now().Format(timeLayout)); err != nil {
		return peer.Response{Status: 500, Message: "format curTime error", Payload: nil}
	} else {
		formattedCurTime = val
	}

	// 值允许这些状态修改，当然前端会自己判断，这里只是方面测试
	// 利用反射判断数组中有没有改状态
	if !IsExistItem(status, UpdateStatusRangeList) {
		return peer.Response{Status: 500, Message: "状态不正确", Payload: nil}
	}

	if orderId == "" || status == "" {
		return peer.Response{Status: 5002, Message: "参数不能为空", Payload: nil}
	}

	resp, order := getOrderInfo(stub, []string{orderId})
	if resp.Status != int32(shim.OK) {
		return resp
	}
	order.Status = status
	// 如果是用户点击确认收货，将订单的还款时间输入
	if order.Status == enumStatus.RefundNoVerified {
		order.RepaymentTime = formattedCurTime
	}

	var key string
	if val, err := stub.CreateCompositeKey("order", []string{orderId}); err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("create key error %s", err), Payload: nil}
	} else {
		key = val
	}

	marshal, err := json.Marshal(order)
	if err != nil {
		return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
	}

	if err = stub.PutState(key, marshal); err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("put statue failed: %s", err), Payload: nil}
	}
	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: nil,
	}
}

// 已完成订单 ✔
func (t *SmartContract) queryOrderCompletionList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}
	userId := args[0]

	if userId == "" {
		return peer.Response{Status: 5002, Message: "参数不能为空", Payload: nil}
	}

	keys := make([]string, 0)
	// 获取单个帐号信息
	resp, user := getUserInfo(stub, []string{userId})
	if resp.Status != int32(shim.OK) {
		return resp
	}

	// 获取所有订单
	orderResult, err := stub.GetStateByPartialCompositeKey("order", keys)
	if err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("query order error: %s", err), Payload: nil}
	}
	defer orderResult.Close()

	orderList := make([]*Order, 0)
	// 1. 先进行判断什么角色区分帐号的角色显示分支
	if user.Role == enumRoles.Seller {
		resp, orderList = getSpecifyStatusOrder(orderResult, SellerFinishOrder, userId) //  发货
	} else {
		resp, orderList = getSpecifyStatusOrder(orderResult, []string{enumStatus.RefundVerified}, userId) // 还款
	}

	if resp.Status != int32(shim.OK) {
		return resp
	}

	// 序列化返回
	orderBytes, err := json.Marshal(orderList)
	if err != nil {
		return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
	}

	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: orderBytes,
	}
}

// 查询采购商的历史交易 ✔
func (t *SmartContract) queryBuyerHistoricalTransactions(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return peer.Response{Status: 5001, Message: "输入的参数个数不正确", Payload: nil}
	}

	orderId := args[0]
	status := args[1]

	if orderId == "" || status == "" {
		return peer.Response{Status: 5002, Message: "参数不能为空", Payload: nil}
	}

	// 申请代付的订单
	resp, order := getOrderInfo(stub, []string{orderId})
	if resp.Status != int32(shim.OK) {
		return resp
	}

	// 通过订单获取改采购商的信息
	resp, user := getUserInfo(stub, []string{order.Buyer})
	if resp.Status != int32(shim.OK) {
		return resp
	}

	keys := make([]string, 0)
	orderResult, err := stub.GetStateByPartialCompositeKey("order", keys)
	if err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("query order error: %s", err), Payload: nil}
	}
	defer orderResult.Close()

	orderList := make([]*Order, 0)
	var BuyerNoPayResp peer.Response
	var BuyerPayResp peer.Response
	// 获取指定用户的指定某些订单的状态
	if IsExistItem(status, []string{enumStatus.DeliveredVerified}) {
		BuyerNoPayResp, orderList = getSpecifyStatusOrder(orderResult, []string{status}, user.Id) // 待还款
		if BuyerNoPayResp.Status != int32(shim.OK) {
			return BuyerNoPayResp
		}
	} else if IsExistItem(status, []string{enumStatus.RefundVerified, enumStatus.RefundNoVerified}) {
		BuyerPayResp, orderList = getSpecifyStatusOrder(orderResult, []string{status}, user.Id) // 待还款
		if BuyerPayResp.Status != int32(shim.OK) {
			return BuyerPayResp
		}
	} else {
		return peer.Response{Status: 500, Message: "状态不正确", Payload: nil}
	}

	// 序列化返回
	orderBytes, err := json.Marshal(orderList)
	if err != nil {
		return peer.Response{Status: 500, Message: "marshal failed", Payload: nil}
	}

	return peer.Response{
		Status:  200,
		Message: "success message",
		Payload: orderBytes,
	}
}

// 获取指定状态订单 // TODO: 其实后续可以升级传入多个状态，进行筛选，弄成列表slice或者map，有底层源码进行查找
func getSpecifyStatusOrder(orderResult shim.StateQueryIteratorInterface, status []string, roleId string) (peer.Response, []*Order) {
	orderList := make([]*Order, 0)
	for orderResult.HasNext() {
		val, err := orderResult.Next()
		if err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("Get order error %s", err), Payload: nil}, orderList
		}

		order := new(Order)
		if err := json.Unmarshal(val.GetValue(), order); err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("Order failed to convert from bytes, error %s", err), Payload: nil}, orderList
		}

		// 如果为空，则只是单纯查指定某状态全部订单
		if roleId != "" {
			// 1.判断该订单是否是属于该帐号应该显示的订单
			if !IsExistItem(roleId, []string{order.Buyer, order.Seller, order.FinancialInstitution}) {
				continue
			}
		}

		// 2. 判断该订单是否是指定的状态为范围内
		if !IsExistItem(order.Status, status) {
			continue
		}

		// 3. 判断什么角色进行处理例外的供货商显示(就是处理已发货的显示问题)
		if roleId == order.Seller && IsExistItem(order.Status, SellerFinishOrder) {
			order.Status = enumStatus.DeliveredVerified
		}

		orderList = append(orderList, order)
	}
	return peer.Response{Status: 200, Message: "success message", Payload: nil}, orderList
}

// 获取单个用户信息
func getUserInfo(stub shim.ChaincodeStubInterface, keys []string) (peer.Response, *User) {

	// 获取帐号信息
	userResult, err := stub.GetStateByPartialCompositeKey("user", keys)
	if err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("query user error: %s", err), Payload: nil}, nil
	}
	defer userResult.Close()

	user := new(User)
	for userResult.HasNext() {
		val, err := userResult.Next()
		if err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("Get user error %s", err), Payload: nil}, nil
		}

		if err := json.Unmarshal(val.GetValue(), user); err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("User failed to convert from bytes, error %s", err), Payload: nil}, nil
		}
	}

	// 补充：在这里查看订单帐号不能为空，当然前端登录后，userId一直存在，不可能user的数据是空的
	if user.Id == "" {
		return peer.Response{Status: 500, Message: "获取不到用户", Payload: nil}, nil
	}

	return peer.Response{Status: 200, Message: "success message", Payload: nil}, user
}

// 获取单个订单信息
func getOrderInfo(stub shim.ChaincodeStubInterface, keys []string) (peer.Response, *Order) {
	// 获取帐号信息
	orderResult, err := stub.GetStateByPartialCompositeKey("order", keys)
	if err != nil {
		return peer.Response{Status: 500, Message: fmt.Sprintf("query order error: %s", err), Payload: nil}, nil
	}
	defer orderResult.Close()

	order := new(Order)
	for orderResult.HasNext() {
		val, err := orderResult.Next()
		if err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("Get order error %s", err), Payload: nil}, nil
		}

		if err := json.Unmarshal(val.GetValue(), order); err != nil {
			return peer.Response{Status: 500, Message: fmt.Sprintf("Order failed to convert from bytes, error %s", err), Payload: nil}, nil
		}
	}
	return peer.Response{Status: 200, Message: "success message", Payload: nil}, order
}

//
func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Sprintln("chaincode start failed", err)
	}
}
