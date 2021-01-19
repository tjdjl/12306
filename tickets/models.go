package tickets

import (
	"fmt"

	"12306.com/12306/common"

	_ "fmt"
)

type Ticket struct { //票模型：明确票的实质是车次的多个连续区间
	TripID uint `gorm:"column:trip_id"`
	//Trip    Trip ``
	StartStationNo uint   `gorm:"column:start_station_no"`
	EndStationNo   uint   `gorm:"column:end_station_no"`
	StartStation   string `gorm:"column:start_station"`
	EndStation     string `gorm:"column:end_station"`
	SeatInfo       TicketSeats
}

//type Trip struct {
//	ID        uint `gorm:"primary_key"`
//	TripID          string           `json:"trainNumber"`
//	TripNo          string           `json:"trainNumber"`
//}

func FindTicketList(startCity, endCity, date string) ([]Ticket, error) { //从数据库找到所有符合条件的车次
	db := common.GetDB()
	var models []Ticket
	err := db.Raw("SELECT A.trip_id AS trip_id ,A.sequence AS start_station_no,B.sequence AS end_station_no FROM (SELECT trip_id,station_name,sequence FROM trip_station WHERE station_name =? AND date(start_time)=? ) A, (SELECT trip_id,station_name,sequence FROM trip_station WHERE station_name =? AND date(start_time)=? ) B WHERE A.sequence < B.sequence AND A.trip_id = B.trip_id ", startCity, date, endCity, date).Find(&models).Error

	//返回
	fmt.Println(models[0].EndStationNo)
	return models, err
}

func FindHishSpeedTicketList(startCity, endCity, date string) ([]Ticket, error) { //从数据库找到所有符合条件的车次
	db := common.GetDB()
	var models []Ticket
	err := db.Raw("SELECT A.trip_id ,A.sequence ,B.sequence FROM (SELECT trip_id,station_name,sequence FROM trip_station WHERE station_name =? AND date(start_time)=? ) A, (SELECT trip_id,station_name,sequence FROM trip_station WHERE station_name =? AND date(start_time)=? ) B WHERE A.sequence < B.sequence AND A.trip_id = B.trip_id ", startCity, date, endCity, date).Find(&models).Error
	//返回
	return models, err
}

//func (self *Ticket) getTrainDetail() error { //获取车次对应的列车的信息
//
//}

type TicketSeats struct {
	ID               uint `gorm:"primary_key"`
	TripID           uint `gorm:"column:name"`
	TripSegmentSeats []TripSegmentSeats
	remainSeatNum    RemainSeats
}
type RemainSeats struct {
	businessSeatsNumber   uint `json:"businessSeatsNumber"`
	firstSeatsNumber      bool `json:"firstSeatsNumber"`
	secondSeatsNumber     uint `json:"secondSeatsNumber"`
	hardSeatsNumber       uint `json:"hardSeatsNumber"`
	hardBerthNumber       uint `json:"hardBerthNumber"`
	softBerthNumber       uint `json:"softBerthNumber"`
	seniorSoftBerthNumber uint `json:"seniorSoftBerthNumber"`
}
type TripSegmentSeats struct {
	businessSeatsNumber   uint `json:"businessSeatsNumber"`
	firstSeatsNumber      bool `json:"firstSeatsNumber"`
	secondSeatsNumber     uint `json:"secondSeatsNumber"`
	hardSeatsNumber       uint `json:"hardSeatsNumber"`
	hardBerthNumber       uint `json:"hardBerthNumber"`
	softBerthNumber       uint `json:"softBerthNumber"`
	seniorSoftBerthNumber uint `json:"seniorSoftBerthNumber"`
}

// //getSeats get the seats data of trainId
//func (s *Ticket) getSeatDetail() (TicketSeats, error) { //获取票的座位余量信息
//	db := common.GetDB()
//	tx := db.Begin()
//	var info []TicketSeatRawInfo  // 批量读取，下标是区间
//	err := tx.Table("trip_segment").Where("trip_id = ? AND segment between ? AND ? and seat_id = ?", s.TripID, s.StartStationNo, s.EndStationNo-1).Error
//
//	//合并TicketSeatRawInfo区间
//
//	//转化为TicketSeatInfo
//	var res TicketSeatInfo
//	return res,err
//}
