package tickets

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//票：   某个车次的某段区间的某个座位；
//查票： 给定出发站、终点站、时间，查找合适的 车次-某段区间，和这些车次-某段区间 对应的各个座位类型的座位余量
//买票： 给定车次-某段区间和座位类型，购买其中的一个座位

//TicketList 查找余票
func TicketList(c *gin.Context) {
	startCity := c.Query("startCity")
	endCity := c.Query("endCity")
	date := c.Query("date")
	category := c.Query("type")
	var tripSeries []TripSeries
	var err error
	if category == "1" {
		tripSeries, err = FindTripSeriesList(startCity, endCity, date) //找到对应的车次
	} else {
		tripSeries, err = FindHishSpeedTripSeriesList(startCity, endCity, date) //找到对应的车次
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": "Invalid param"})
		return
	}
	serializer := TicketListSerializer{c, tripSeries}                                     //新建序列化器
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "查找成功", "data": serializer.Response()}) //返回结果
}

// TicketBuy 买票
func TicketBuy(c *gin.Context) {
	tripID, _ := strconv.Atoi(c.Query("tripID"))
	startStationNo, _ := strconv.Atoi(c.Query("startStationNo"))

	endStationNo, _ := strconv.Atoi(c.Query("endStationNo"))
	seatCategory := c.Query("seatCategory")
	tripSegment := TripSeries{uint(tripID), uint(startStationNo), uint(endStationNo)}
	err := tripSegment.orderOneSeat(seatCategory) //找到空闲的座位号；
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": "Invalid param"})
		return
	}
	// serializer := TicketsSerializer{c, TicketsModel}                                      //新建序列化器
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "成功"}) //返回结果
}
