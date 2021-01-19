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

func TrainList(c *gin.Context) {       //查找余票，实际上是查询车次信息
		startCity := c.Query("startCity")
		endCity := c.Query("endCity")
		date := c.Query("date")
		catogory := c.Query("type")

		TrainsModel, err := FindTrainList(startCity,endCity,date,catogory)   //找到对应的车次
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": "Invalid param"})
			return
		}
		serializer := TrainsSerializer{c, TrainsModel{}}  //新建序列化器
		c.JSON(http.StatusOK, gin.H{"list": TrainsSerializer{}.Response(),)  //返回结果
	}
}
