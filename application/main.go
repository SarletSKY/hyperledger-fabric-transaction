package main

import (
	"accurchain.com/perishable-food/application/blockchain"
	"accurchain.com/perishable-food/application/controller"
	"github.com/gin-gonic/gin"
)

// 路由配置
func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/queryUserInfo", controller.QueryUserInfo)                           // 查询帐号
	router.GET("/queryOrderList", controller.QueryOrderList)                         // 查询订单
	router.GET("/queryCommodityList", controller.QueryCommodityList)                 // 查询商品
	router.GET("/queryApprovalOrder", controller.QueryApprovalOrder)                 // 查询审批订单栏
	router.GET("/queryFinancialApprovalList", controller.QueryFinancialApprovalList) // 获取金融机构列表
	router.GET("/queryOrderCompletionList", controller.QueryOrderCompletionList)     // 查询已经完成订单
	router.POST("/createUser", controller.CreateUser)                                // 创建帐号
	router.POST("/createOrder", controller.CreateOrder)                              // 创建订单
	router.POST("/confirmOrder", controller.ConfirmOrder)                            // 确认订单
	router.POST("/createCommodity", controller.CreateCommodity)                      // 创建商品
	router.POST("/toApplyForPaidOrder", controller.ToApplyForPaidOrder)              // 申请代付操作
	router.POST("/financialApproval", controller.FinancialApproval)                  // 审批订单
	router.POST("/confirmUserPwd", controller.ConfirmUserPwd)                        // 确认密码
	router.POST("/updateOrderStatus", controller.UpdateOrderStatus)                  // 更新订单状态
	return router
}

func main() {
	// 初始化sdk
	blockchain.Init()
	// 运行
	router := setupRouter()
	router.Run()
}
