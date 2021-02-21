package trains

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TrainStationList 查询某趟车次经过的站点
func TrainStationList(c *gin.Context) {
	//获取参数
	trainID := c.Query("train_no")
	isToday := true
	//查询该车次经过的站点
	train := Train{trainID}
	models, err := train.getTrainStaions()
	//响应
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": err.Error()})
		return
	}
	serializer := TrainStaionSerializer{C: c, TrainStaions: models, IsToday: isToday}     //新建序列化器
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "查找成功", "data": serializer.Response()}) //返回结果
}

//票：   某个车次的某段区间的某个座位；
//查票： 给定出发站、终点站、时间，查找合适的 车次-站点对，和这些车次-站点对 对应的各个座位类型的座位余量
//买票： 给定车次-站点对和座位类型，购买其中的一个座位

//TicketList 查找余票，本质是查找车次站点对及其余票
func TicketList(c *gin.Context) {
	//获取参数
	startCity := c.Query("startCity")
	endCity := c.Query("endCity")
	date := c.Query("date")
	category := c.Query("type")
	//检验date是否是合法日期，并把date转换成2006-02-03这种格式
	// _, err := time.ParseInLocation("2006-01-02", date, time.Local)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": err.Error()})
	// 	return
	// }

	//查找车次站点对
	var models []TrainStaionPair
	models, err := ListTrainStaionPair(startCity, endCity, date, category == "1")
	//响应
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": err.Error()})
		return
	}
	serializer := TicketsSerializer{C: c, TripStaionPairs: models, Date: date}            //新建序列化器
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "查找成功", "data": serializer.Response()}) //返回结果
}

//TicketBuy 买票
func TicketBuy(c *gin.Context) {
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

//TicketCancel 退票
func TicketCancel(c *gin.Context) {
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

//TicketChange 改签
func TicketChange(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Query("ticket_outside_id"), 10, 32)
	tripID := c.Query("tripID")
	startStationNo, err := strconv.ParseUint(c.Query("startStationNo"), 10, 32)
	endStationNo, err := strconv.ParseUint(c.Query("endStationNo"), 10, 32)
	seatCategory := c.Query("seatCategory")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": "Invalid id"})
		return
	}
	newTrip := Trip{tripID}
	err = newTrip.changeOrder(uint(id64), uint(startStationNo), uint(endStationNo), seatCategory)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "改签成功"})
}
