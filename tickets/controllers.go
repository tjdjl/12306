package trains

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// {
//     "startCity":"", // 城市名或站名
//     "endCity":"",
//     "date":"",
//     "type":"" // 0 全类, 1高铁动车票
// }

func TicketList(c *gin.Context) {       //查找余票，实际上是查询车次信息
		startCity := c.Query("startCity")
		endCity := c.Query("endCity")
		date := c.Query("date")
		category := c.Query("type")
		var TicketsModel []Ticket
		var err error
		if category == "1"{
			TicketsModel,err = FindTicketList(startCity,endCity,date)   //找到对应的车次
		} else{
		TicketsModel, err = FindHishSpeedTicketList(startCity,endCity,date)   //找到对应的车次
		}

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": "Invalid param"})
			return
		}
		serializer := TicketsSerializer{c, TicketsModel }  //新建序列化器
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "查找成功", "data": serializer.Response()})  //返回结果
}

