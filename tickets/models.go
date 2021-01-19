package trains

import "12306.com/12306/common"

import (
	_ "fmt"

)

type Ticket struct { //票模型：明确票的实质是车次的多个连续区间
	ID        uint `gorm:"primary_key"`
	TripID    uint `gorm:"column:name"`
	SegmentFromNo uint `gorm:"column:name"`
	SegmentToNo uint `gorm:"column:name"`SS
}


func FindTicketList(data interface{}) ([]Ticket, error) { //从数据库找到所有符合条件的车次
	db := common.GetDB()
	var models []Ticket
	//返回
	return
}



func (self *Ticket) getTrainDetail() error { //获取车次对应的列车的信息

}

// //getSeats get the seats data of trainId
func (s *Ticket) getSeatDetail() ([]Ticket, error) { //获取票的座位余量信息

}