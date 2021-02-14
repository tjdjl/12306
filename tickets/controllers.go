package tickets

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//票：   某个车次的某段区间的某个座位；
//查票： 给定出发站、终点站、时间，查找合适的 车次-某段区间，和这些车次-某段区间 对应的各个座位类型的座位余量
//买票： 给定车次-某段区间和座位类型，购买其中的一个座位

//ListTicket 查找余票
func ListTicket(c *gin.Context) {
	//获取参数
	startCity := c.Query("startCity")
	endCity := c.Query("endCity")
	date := c.Query("date")
	category := c.Query("type")
	//检验合法日期，并转换成2006-02-03这种格式
	// _, err := time.ParseInLocation("2006-01-02", date, time.Local)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": err.Error()})
	// 	return
	// }
	var err error
	//查找车次
	var models []TrainStaionPair
	if category == "1" {
		models, err = FindTrainStaionPairList(startCity, endCity, date, false)
	} else {
		models, err = FindTrainStaionPairList(startCity, endCity, date, true)
	}
	//响应
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": err.Error()})
		return
	}
	serializer := TicketListSerializer{c, models, date}                                   //新建序列化器
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "查找成功", "data": serializer.Response()}) //返回结果
}

// BuyTicket 买票
func BuyTicket(c *gin.Context) {
	//获取参数
	tripID := c.Query("tripID")
	startStationNo, err := strconv.ParseUint(c.Query("startStationNo"), 10, 32)
	endStationNo, err := strconv.ParseUint(c.Query("endStationNo"), 10, 32)
	seatCategory := c.Query("seatCategory")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": err.Error()})
		return
	}
	//买票
	trip := Trip{tripID}
	err = trip.orderOneSeat(uint(startStationNo), uint(endStationNo), seatCategory) //找到空闲的座位号；
	//响应
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "买票成功"}) //返回结果
}

//CancelTicket 退票
func CancelTicket(c *gin.Context) {
	//获取参数
	id64, err := strconv.ParseUint(c.Query("ticket_outside_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": "Invalid id"})
		return
	}
	//退票
	trip := Trip{}
	err = trip.cancleOrder(uint(id64))
	//响应
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "退票成功"})
}

//ChangeTicket 改签
func ChangeTicket(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Query("ticket_outside_id"), 10, 32)
	id := uint(id64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": "Invalid id"})
		return
	}
	tripID64, err := strconv.ParseUint(c.Query("tripID"), 10, 32)
	tripID := uint(tripID64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": "Invalid trip id"})
		return
	}
	seatCategory := c.Query("seatCategory")
	trip:=Trip{}
	err=trip.changeOrder(id,tripID,seatCategory)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "改签成功"})
}
