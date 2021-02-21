package main

import (
	"12306.com/12306/stations"
	"12306.com/12306/trains"
	"12306.com/12306/users"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	//users
	r.POST("/user/register/", users.Register)
	r.POST("/user/login/", users.Login)
	// r.GET("/api/auth/info", users.AuthMiddleware(), users.Info)

	//stations
	//查找所有站点
	r.GET("/search/api/v1/queryAllStations/", stations.AllStationsList)

	//trains
	//查某车次经过的站点
	r.GET("/search/api/v1/queryStation/", trains.TrainStationList)
	//查票
	r.POST("/search/api/v1/remainder/", trains.TicketList)
	//买票
	r.GET("/buy/ticket/", trains.TicketBuy)
	//退票
	r.POST("/reticket/api/v1/", trains.TicketCancel)
	// r.PUT("/change/order/", orders.ChangeTicket)

	return r
}
