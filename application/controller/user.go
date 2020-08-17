package controller

import (
	"accurchain.com/perishable-food/application/blockchain"
	"accurchain.com/perishable-food/application/lib"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createUserRequest struct {
	Id               string  `json:"id" form:"id" binding:"required"`                               // 用户Id
	UserName         string  `json:"user_name" form:"user_name" binding:"required"`                 // 帐号名
	Password         string  `json:"password" form:"password" binding:"required"`                   // 密码
	PayPassword      string  `json:"pay_password" form:"pay_password" binding:"required"`           // 支付密码
	Phone            string  `json:"phone" form:"phone" binding:"required"`                         // 手机
	Email            string  `json:"email" form:"email" binding:"required"`                         // 电子邮箱
	Enterprise       string  `json:"enterprise" form:"enterprise" binding:"required"`               // 企业
	EnterpriseOrigin string  `json:"enterprise_origin" form:"enterprise_origin" binding:"required"` // 企业地址
	EnterpriseCredit int64   `json:"enterprise_credit" form:"enterprise_credit" binding:"required"` // 企业信用
	Balance          float64 `json:"balance" form:"balance" binding:"required"`                     // 资金
	Role             string  `json:"role" form:"role" binding:"required"`                           // 充当角色
}

func CreateUser(ctx *gin.Context) {
	req := new(createUserRequest)

	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 打印请求体
	fmt.Println("请求参数：")
	marshal, _ := json.Marshal(req)
	fmt.Println(string(marshal))

	resp, err := blockchain.ChannelExecute("createUser", [][]byte{
		[]byte(req.Id),
		[]byte(req.UserName),
		[]byte(req.Password),
		[]byte(req.PayPassword),
		[]byte(req.Phone),
		[]byte(req.Email),
		[]byte(req.Enterprise),
		[]byte(req.EnterpriseOrigin),
		[]byte(fmt.Sprintf("%d", req.EnterpriseCredit)),
		[]byte(fmt.Sprintf("%v", req.Balance)),
		[]byte(req.Role),
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

type queryUserRequest struct {
	Id string `json:"id" form:"id" binding:"required"` // 用户Id
}

func QueryUserInfo(ctx *gin.Context) {
	// 解析请求体
	req := new(queryUserRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp, err := blockchain.ChannelQuery("queryUserInfo", [][]byte{
		[]byte(req.Id),
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 将结果返回
	ctx.String(http.StatusOK, bytes.NewBuffer(resp.Payload).String())
}

func QueryApprovalOrder(ctx *gin.Context) {
	// 解析请求体
	req := new(queryUserRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp, err := blockchain.ChannelQuery("queryApprovalOrder", [][]byte{
		[]byte(req.Id),
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var data []map[string]interface{}
	json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data)

	// 将结果返回
	ctx.JSON(http.StatusOK, data)
}

// 查询金融机构列表
func QueryFinancialApprovalList(ctx *gin.Context) {
	resp, err := blockchain.ChannelQuery("queryFinancialApprovalList", [][]byte{})

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

// 确认密码
func ConfirmUserPwd(ctx *gin.Context) {
	req := new(lib.User)

	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 打印请求体
	fmt.Println("请求参数：")
	marshal, _ := json.Marshal(req)
	fmt.Println(string(marshal))

	resp, err := blockchain.ChannelExecute("confirmUserPwd", [][]byte{
		[]byte(req.Id),
		[]byte(req.UserName),
		[]byte(req.Password),
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
