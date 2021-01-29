package orders

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//CancelTicket
func CancelTicket(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	id := uint(id64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": "Invalid id"})
		return
	}
	orderModel := Order{ID: id}
	err = orderModel.cancleOrder()
	fmt.Print(err)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": "cancel order wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "退票成功", "data": ""})
}

// //ChangeTicket
// func ChangeTicket(c *gin.Context) {
// 	orderModel := Order{orderId: c.Query("id")}
// 	date := c.Query("date")
// 	err := orderModel.changeOrder(date)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": "Invalid param"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "改票成功", "data": ""})
// }
