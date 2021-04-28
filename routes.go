package main

import (
	"12306.com/12306/common/middleware"
	"12306.com/12306/stations"
	"12306.com/12306/trains"
	"12306.com/12306/users"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	//users
	//注册
	r.POST("/user/api/v1/register/", users.Register)
	//登录
	r.POST("/user/api/v1/login/", users.Login)
	//添加乘车人
	r.POST("/user/api/v1/passenger/", middleware.AuthMiddleware(), users.AddPassenger)
	//修改乘车人
	r.PUT("/user/api/v1/passenger/", middleware.AuthMiddleware(), users.UpdatePassenger)
	//查询乘车人
	r.GET("/user/api/v1/passenger/", middleware.AuthMiddleware(), users.QueryPassenger)

	//stations
	//查找所有站点
	r.GET("/search/api/v1/queryAllStations/", stations.AllStationsList)

	//trains
	//查某车次经过的站点
	r.GET("/search/api/v1/queryStation/", trains.TrainStationList)
	//查票
	r.POST("/search/api/v1/remainder/", trains.TicketList)
	//买票
	r.GET("/buy/ticket/", middleware.AuthMiddleware(), trains.TicketBuy)
	// r.GET("/buy/ticket/", middleware.AuthMiddleware(), trains.TicketBuy)
	//退票
	r.POST("/reticket/api/v1/", middleware.AuthMiddleware(), trains.TicketCancel)
	//改票
	r.PUT("/change/order/", trains.TicketChange)

	return r
}
