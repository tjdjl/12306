package main

import (
	"12306.com/12306/stations"
	"12306.com/12306/tickets"
	"12306.com/12306/users"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	//users
	r.POST("/user/register/", users.Register)
	r.POST("/user/login/", users.Login)
	// r.GET("/api/auth/info", users.AuthMiddleware(), users.Info)

	//stations
	r.GET("/search/allstations/", stations.StationList)

	//tickets
	//查票
	r.GET("/search/remainder/", tickets.ListTicket)
	//买票
	r.GET("/buy/ticket/", tickets.BuyTicket)
	//退票
	r.PUT("/cancel/order/:id", tickets.CancelTicket)
	// r.PUT("/change/order/", orders.ChangeTicket)

	return r
}
