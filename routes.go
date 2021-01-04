package main

import (
	"12306.com/12306/orders"
	"12306.com/12306/stations"
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

	//orders
	r.PUT("/cancel/order/:id", orders.CancelTicket)
	// r.PUT("/change/order/", orders.ChangeTicket)

	//trains

	return r
}
