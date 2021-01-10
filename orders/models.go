package orders

import (
	"errors"
	"fmt"
	_ "fmt"
	"time"

	"12306.com/12306/common"
)

type Order struct {
	ID           uint      `gorm:"primary_key;auto_increment" json:"id"`
	OrderID      string    `json:"order_id"`
	UserID       uint      `json:"user_id"`
	TripID       uint      `json:"trip"`
	SeatID       uint      `json:"seat_id"`
	StartNo      uint      `json:"startNo"`
	EndNo        uint      `json:"endNo"`
	StartStation string    `json:"startStation"`
	EndStation   string    `json:"endStation"`
	Date         time.Time `json:"date"`
	Status       string    `json:"status"`
}

// func SaveOneOrder(data interface{}) error { //下单
// 	// 1.找到一个空闲的座位
// 	// SELECT seatid
// 	// FROM （
// 	// Select seatid where trainid = 1 and qujian =0 and status = 0 for update union
// 	// Select seatid where trainid = 1 and qujian =1 and status = 0 for update union
// 	// Select seatid where trainid = 1 and qujian =2 and status = 0 for update union
// 	// ）as seat_table
// 	// order by id asc  LIMIT 1
// 	// 2.修改它的区间状态为已占用
// 	// update status =1
// 	// where trainid = 3
// 	// and qujian between 1 and 3
// 	// and seatid  =5
// 	// status = 0;
// 	// 判断row affected;决定是否回滚；
// 	// 3. 生成订单
// }

func (order *Order) cancleOrder() error { //取消订单
	db := common.GetDB()
	tx := db.Begin()
	// 0.get order details
	err := tx.Set("gorm:query_option", "FOR UPDATE").First(&order, order.ID).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if !(order.Status == "未支付" || order.Status == "已支付") {
		tx.Rollback()
		return errors.New("订单状态错误,不支持退票")
	}
	//判断下单的用户是不是登录用户本人，借助中间件

	// 	1.update seat table
	err = tx.Table("trip_seat_segment").Where("trip_id = ? AND segment between ? AND ? and seat_id = ?", order.TripID, order.StartNo, order.EndNo-1, order.SeatID).Updates(map[string]interface{}{"status": 0}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2.return money
	if order.Status == "已支付" {
		fmt.Print("退钱给用户")
	}

	// 3.update order table
	err = tx.Model(&order).Updates(map[string]interface{}{"status": "已退票"}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 4.commit
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil

}

// func (self *Order) changeOrder(date string) error { //修改订单，只能改时间，不能改区间
// 	// 0.找到train id ,等信息,锁订单
// 	// 1.修改原来座位的区间状态为空白（首先还座位，避免还之前查不到座位还之后查得到座位的情况）
// 	// update status =0
// 	// where trainid = 2
// 	// and qujian between 1 and 3
// 	// and seatid  =5
// 	// 2.找到一个空闲的座位 参考悲观锁
// 	// 3.修改它的区间状态为已占用
// 	// update status =1
// 	// where trainid = 3
// 	// and qujian between 1 and 3
// 	// and seatid  =5
// 	// 4.修改订单状态
// }
