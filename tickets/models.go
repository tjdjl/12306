package trains

import "12306.com/12306/common"

import (
	_ "fmt"

)

type Ticket struct { //票模型：明确票的实质是车次的多个连续区间
	ID        uint `gorm:"primary_key"`
	TripID    uint `gorm:"column:trip_id"`
	//Trip    Trip ``
	StartStationNo uint `gorm:"column:start_station_no"`
	EndStationNo uint `gorm:"column:end_station_no"`
	StartStation           string     `gorm:"column:start_station"`
	EndStation    string              `gorm:"column:end_station"`
	SeatInfo TicketSeatInfo
	SeatRawInfo []TicketSeatRawInfo
}


//type Trip struct {
//	ID        uint `gorm:"primary_key"`
//	TripID          string           `json:"trainNumber"`
//	TripNo          string           `json:"trainNumber"`
//}
type TicketSeatInfo struct {
	businessSeatsNumber    uint              `json:"businessSeatsNumber"`
	firstSeatsNumber       bool                  `json:"firstSeatsNumber"`
	secondSeatsNumber uint                  `json:"secondSeatsNumber"`
	hardSeatsNumber uint                  `json:"hardSeatsNumber"`
	hardBerthNumber uint                  `json:"hardBerthNumber"`
	softBerthNumber uint                  `json:"softBerthNumber"`
	seniorSoftBerthNumber uint                  `json:"seniorSoftBerthNumber"`
}


func FindTicketList(startCity,endCity,date string) ([]Ticket, error) { //从数据库找到所有符合条件的车次
	db := common.GetDB()
	var models []Ticket
	//SELECT trip_id ,A.sequence as start_station_no,B.sequence as end_station_no
	//FROM (SELECT trip_id,station_name,sequence FROM trip_station WHERE station_name = )as A and as ()B
	//WHERE A.sequence < B.sequence

	//返回
	return models,err
}

func FindHishSpeedTicketList(startCity,endCity,date string) ([]Ticket, error) { //从数据库找到所有符合条件的车次
	db := common.GetDB()
	var models []Ticket
	//返回
	return models,err
}


//func (self *Ticket) getTrainDetail() error { //获取车次对应的列车的信息
//
//}


type TicketSeatRawInfo struct {
	ID        uint `gorm:"primary_key"`
	TripID    uint `gorm:"column:name"`
	businessSeatsNumber    uint              `json:"businessSeatsNumber"`
	firstSeatsNumber       bool                  `json:"firstSeatsNumber"`
	secondSeatsNumber uint                  `json:"secondSeatsNumber"`
	hardSeatsNumber uint                  `json:"hardSeatsNumber"`
	hardBerthNumber uint                  `json:"hardBerthNumber"`
	softBerthNumber uint                  `json:"softBerthNumber"`
	seniorSoftBerthNumber uint                  `json:"seniorSoftBerthNumber"`
}
// //getSeats get the seats data of trainId
func (s *Ticket) getSeatDetail() (TicketSeatInfo, error) { //获取票的座位余量信息
	db := common.GetDB()
	tx := db.Begin()
	var info []TicketSeatRawInfo  // 下标是区间
	err:= db.Where() //

	//合并TicketSeatRawInfo区间

	//转化为TicketSeatInfo
	var res TicketSeatInfo
	return res,err
}


