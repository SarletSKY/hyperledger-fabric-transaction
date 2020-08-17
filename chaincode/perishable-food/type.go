package main

import "time"

// 帐号信息
// 登录账号、登录密码、手机号、电子邮箱、企业名称、企业 地址、企业统一社会信用代码
type User struct {
	Id               string  `json:"id"`                // 用户Id
	UserName         string  `json:"user_name"`         // 帐号名
	Password         string  `json:"password"`          // 密码
	PayPassword      string  `json:"pay_password"`      // 支付密码
	Phone            string  `json:"phone"`             // 手机
	Email            string  `json:"email"`             // 电子邮箱
	Enterprise       string  `json:"enterprise"`        // 企业
	EnterpriseOrigin string  `json:"enterprise_origin"` // 企业地址
	EnterpriseCredit int64   `json:"enterprise_credit"` // 企业信用
	Balance          float64 `json:"balance"`           // 资金
	Role             string  `json:"role"`              // 充当角色
}

// 商品结构体
type Commodity struct {
	User      *User   `json:"user"`      // 商品属于那个帐号
	Id        string  `json:"id"`        // 商品ID
	Name      string  `json:"name"`      // 商品名称
	Origin    string  `json:"origin"`    // 生产地
	Price     float64 `json:"price"`     // 价格
	Introduce string  `json:"introduce"` // 介绍
}

// 订单的结构体
// 需填写包括：商品名称、订单金额、供应商(只能从已注册的 供应商用户下拉列表中选择填入，无法主动输入)。订单建立后，采购商能够在订单列表查看到相应的订 单详情，此时订单状态为：新建(待确认)。
type Order struct {
	Commodity            *Commodity `json:"commodity"`             // 商品名称
	Id                   string     `json:"id"`                    // 订单id
	PayOrderNo           string     `json:"pay_order_no"`          // 付款单号
	Buyer                string     `json:"buyer"`                 // 采购商
	FinancialInstitution string     `json:"financial_institution"` // 金融机构
	Seller               string     `json:"seller"`                // 供应商
	Quantity             float64    `json:"quantity"`              // 数量
	OrderTotalAmount     float64    `json:"order_total_amount"`    // 订单总金额
	DeliverAddress       string     `json:"deliver_address"`       // 配送地址
	Status               string     `json:"status"`                // 订单状态
	OrderTime            time.Time  `json:"order_time"`            // 下单时间
	PromiseRepaymentTime time.Time  `json:"promise_time"`          // 承诺还款时间
	RepaymentTime        time.Time  `json:"repayment_time"`        // 实际还款时间
}

// 角色
type Role struct {
	Buyer                string
	Seller               string
	FinancialInstitution string
}

// 订单状态
type Status struct {
	NewUnconfirmed string // 新建(待确认)
	NewConfirmed   string // 新建(已确认)

	RefundVerified    string //已还款(已核实)
	RefundNoVerified  string //已还款(未核实)
	DeliveredVerified string //已发货(已核实)

	ToApplyForPaid      string //申请代付
	UnPayment           string //代付被拒
	PaymentNoVerified   string //已付款(未核实)
	PaymentVerified     string //已付款(已核实)
	DeliveredNoVerified string //已发货(未核实)
}
