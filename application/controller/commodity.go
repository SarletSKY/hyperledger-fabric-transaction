package controller

import (
	"accurchain.com/perishable-food/application/blockchain"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type commodityRequest struct {
	UserId    string  `json:"user_id" form:"user_id" binding:"required"`     // 用户
	Name      string  `json:"name" form:"name" binding:"required"`           // 名称
	Id        string  `json:"id" form:"id" binding:"required"`               // 商品Id
	Origin    string  `json:"origin" form:"origin" binding:"required"`       // 商品生产地
	Price     float64 `json:"price" form:"price" binding:"required"`         // 价格
	Introduce string  `json:"introduce" form:"introduce" binding:"required"` // 商品简介
}

// 创建商品
func CreateCommodity(ctx *gin.Context) {
	// 解析请求体
	req := new(commodityRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 打印请求体
	fmt.Println("请求参数：")
	marshal, _ := json.Marshal(req)
	fmt.Println(string(marshal))

	// 将请求体参数转换byte数组，发送给区块链
	resp, err := blockchain.ChannelExecute("createCommodity", [][]byte{
		[]byte(req.UserId),
		[]byte(req.Id),
		[]byte(req.Name),
		[]byte(req.Origin),
		[]byte(fmt.Sprintf("%v", req.Price)),
		[]byte(req.Introduce),
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// 查询商品
func QueryCommodityList(ctx *gin.Context) {
	// 解析请求体
	resp, err := blockchain.ChannelQuery("queryCommodityList", [][]byte{})

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
