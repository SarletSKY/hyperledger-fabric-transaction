package controller

import (
	"accurchain.com/perishable-food/application/blockchain"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	timeLayout = "2006-01-02T15:04:05Z07:00"
)

// 订单
type orderRequest struct {
	CommodityId          string    `json:"commodity_id" form:"commodity_id"`             // 商品名称
	Id                   string    `json:"id" form:"id"`                                 // 订单id
	PayOrderNo           string    `json:"pay_order_no" form:"pay_order_no"`             // 付款单号
	Buyer                string    `json:"buyer" form:"buyer"`                           // 采购商
	FinancialInstitution string    `json:"financial_institution"`                        // 金融机构
	Seller               string    `json:"seller" form:"seller"`                         // 供应商
	Quantity             float64   `json:"quantity" form:"quantity"`                     // 数量
	OrderTotalAmount     float64   `json:"order_total_amount" form:"order_total_amount"` // 订单总金额
	DeliverAddress       string    `json:"deliver_address" form:"deliver_address"`       // 配送地址
	Status               string    `json:"status" form:"status"`                         // 订单状态
	OrderTime            time.Time `json:"order_time" form:"order_time"`                 // 下单时间
	PromiseRepaymentTime time.Time `json:"promise_time" form:"promise_time"`             // 承诺还款时间
	RepaymentTime        time.Time `json:"repayment_time" form:"repayment_time"`         // 实际还款时间
	Role                 string    `json:"role" form:"role"`                             // 角色
	User                 string    `json:"user" form:"user"`                             // 帐号
	PayPassword          string    `json:"pay_password" form:"pay_password"`             // 支付密码
}

// 创建订单
func CreateOrder(ctx *gin.Context) {
	// 解析请求体
	req := new(orderRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 打印请求体
	fmt.Println("请求参数：")
	marshal, _ := json.Marshal(req)
	fmt.Println(string(marshal))

	resp, err := blockchain.ChannelExecute("createOrder", [][]byte{
		[]byte(req.Id),
		[]byte(req.CommodityId),
		[]byte(req.Buyer),
		[]byte(req.Seller),
		[]byte(fmt.Sprintf("%v", req.Quantity)),
		[]byte(req.DeliverAddress),
		[]byte(time.Now().Format(timeLayout)),
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// 确认订单
func ConfirmOrder(ctx *gin.Context) {
	req := new(orderRequest)

	// 打印请求体
	fmt.Println("请求参数：")
	marshal, _ := json.Marshal(req)
	fmt.Println(string(marshal))

	if err := ctx.BindJSON(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp, err := blockchain.ChannelExecute("confirmOrder", [][]byte{
		[]byte(req.Id),
		[]byte(req.Role),
		[]byte(req.Status),
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 将结果返回
	ctx.JSON(http.StatusOK, resp)
}

// 审批订单
func FinancialApproval(ctx *gin.Context) {
	req := new(orderRequest)
	if err := ctx.BindJSON(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 打印请求体
	fmt.Println("请求参数：")
	marshal, _ := json.Marshal(req)
	fmt.Println(string(marshal))

	resp, err := blockchain.ChannelExecute("financialApproval", [][]byte{
		[]byte(req.Id),
		[]byte(req.User),
		[]byte(req.Status),
		[]byte(req.PayPassword),
		[]byte(time.Now().Format(timeLayout)),
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// 申请代付订单
func ToApplyForPaidOrder(ctx *gin.Context) {
	req := new(orderRequest)
	if err := ctx.BindJSON(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 打印请求体
	fmt.Println("请求参数：")
	marshal, _ := json.Marshal(req)
	fmt.Println(string(marshal))

	resp, err := blockchain.ChannelExecute("toApplyForPaidOrder", [][]byte{
		[]byte(req.Id),
		[]byte(req.User),
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// 更新订单状态
func UpdateOrderStatus(ctx *gin.Context) {
	// 解析请求体
	req := new(orderRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 打印请求体
	fmt.Println("请求参数：")
	marshal, _ := json.Marshal(req)
	fmt.Println(string(marshal))

	resp, err := blockchain.ChannelExecute("updateOrderStatus", [][]byte{
		[]byte(req.Id),
		[]byte(req.Status),
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// 查询订单
func QueryOrderList(ctx *gin.Context) {
	// 解析请求体
	req := new(orderRequest)

	orderId := ctx.Request.FormValue("orderId")
	userId := ctx.Request.FormValue("userId")
	var args [][]byte
	args = append(args, []byte(userId))

	if orderId != "" {
		args = append(args, []byte(orderId))
	}

	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp, err := blockchain.ChannelQuery("queryOrderList", args)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data)

	// 将结果返回
	ctx.JSON(http.StatusOK, data)
}

// 查询已完成订单
func QueryOrderCompletionList(ctx *gin.Context) {
	userId := ctx.Request.FormValue("userId")

	resp, err := blockchain.ChannelQuery("queryOrderCompletionList", [][]byte{
		[]byte(userId),
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 反序列化json
	var data []map[string]interface{}
	json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data)

	// 将结果返回
	ctx.JSON(http.StatusOK, data)
}

// 查询采购商的历史订单
func QueryBuyerHistoricalTransactions(ctx *gin.Context) {

	req := new(orderRequest)

	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp, err := blockchain.ChannelQuery("queryBuyerHistoricalTransactions", [][]byte{
		[]byte(req.Id),
		[]byte(req.Status),
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 反序列化json
	var data []map[string]interface{}
	json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data)

	// 将结果返回
	ctx.JSON(http.StatusOK, data)
}
