package tickets

import (
	"fmt"
	"time"

	"12306.com/12306/common"
)

//TripSeries 一个车次的多个连续区间，没有直接对应的数据库表
type TripSeries struct {
	TripID         uint `gorm:"column:trip_id"`
	StartStationNo uint `gorm:"column:start_station_no"` //StationNo为1代表是该车次的起始站
	EndStationNo   uint `gorm:"column:end_station_no"`
	// StartStation   string `gorm:"column:start_station"`
	// EndStation     string `gorm:"column:end_station"`
}

//RemainSeats 该结构体用于表示，对应给定的TripSegment结构体，它的各个座位类型的余票数
type RemainSeats struct {
	BusinessSeats   uint `json:"businessSeatsNumber"`
	FirstSeats      uint `json:"firstSeatsNumber"`
	SecondSeats     uint `json:"secondSeatsNumber"`
	HardSeats       uint `json:"hardSeatsNumber"`
	HardBerth       uint `json:"hardBerthNumber"`
	SoftBerth       uint `json:"softBerthNumber"`
	SeniorSoftBerth uint `json:"seniorSoftBerthNumber"`
}

//TripSegmentSeats 对应trip_segment的一行中的座位那一列
type TripSegmentSeats struct {
	Seats []uint8 `gorm:"column:seats"`
}

type Order struct {
	ID             uint      `gorm:"primary_key;auto_increment" json:"id"`
	TripID         uint      `json:"trip_id"`
	StartStationNo uint      `json:"start_station_no"`
	EndStationNo   uint      `json:"end_station_no"`
	SeatNo         uint      `json:"seat_no"`
	SeatCatogory   string    `json:"seat_catogory"`
	UserID         uint      `json:"user_id"`
	StartStation   string    `json:"startStation"`
	EndStation     string    `json:"endStation"`
	Date           time.Time `json:"date"`
	Status         string    `json:"status"`
}

//Trip 对应数据库表中的车次
// type Trip struct {
// 	ID     uint   `gorm:"primary_key"`
// 	TripID string `json:"trip_id"`
// 	TripNo string `json:"Trip_No"`
// }

//FindTripSeriesList 根据相应的startCity, endCity,date 条件，返回TripSeries列表
func FindTripSeriesList(startCity, endCity, date string) ([]TripSeries, error) { //从数据库找到所有符合条件的车次
	db := common.GetDB()
	var models []TripSeries
	err := db.Raw("SELECT A.trip_id AS trip_id ,A.sequence AS start_station_no,B.sequence AS end_station_no FROM (SELECT trip_id,station_name,sequence FROM trip_station WHERE station_name =? AND date(start_time)=? ) A, (SELECT trip_id,station_name,sequence FROM trip_station WHERE station_name =? AND date(start_time)=? ) B WHERE A.sequence < B.sequence AND A.trip_id = B.trip_id ", startCity, date, endCity, date).Find(&models).Error
	return models, err
}

//FindHishSpeedTripSeriesList 根据相应的startCity, endCity, date条件，返回高铁快车的TripSeries列表
func FindHishSpeedTripSeriesList(startCity, endCity, date string) ([]TripSeries, error) { //从数据库找到所有符合条件的车次
	db := common.GetDB()
	var models []TripSeries
	err := db.Raw("SELECT A.trip_id ,A.sequence ,B.sequence FROM (SELECT trip_id,station_name,sequence FROM trip_station WHERE station_name =? AND date(start_time)=? ) A, (SELECT trip_id,station_name,sequence FROM trip_station WHERE station_name =? AND date(start_time)=? ) B WHERE A.sequence < B.sequence AND A.trip_id = B.trip_id ", startCity, date, endCity, date).Find(&models).Error
	//返回
	return models, err
}

//TripSegment 对应trip_segment1的一行
type TripSegment struct {
	ID           uint    `gorm:"column:id"`
	TripID       uint    `gorm:"column:trip_id"`
	SegmentNo    uint    `gorm:"column:segment_no"`
	SeatCatogory string  `gorm:"column:seat_catogory"`
	SeatBytes    []uint8 `gorm:"column:seat_bytes"`
}

//getRemainSeats 对于给定的TripSeries，根据座位类型,返回座位余量
func (s *TripSeries) getRemainSeats(catogory string) uint { //获取票的座位余量信息
	db := common.GetDB()
	//business_seats
	var info []TripSegment
	string := "SELECT * FROM trip_segment WHERE trip_id = ? AND segment_no between ? AND ? AND seat_catogory = ?"
	db.Raw(string, s.TripID, s.StartStationNo, s.EndStationNo-1, catogory).Find(&info)
	fmt.Println("座位：", info)
	res := calculasRemainSeats(info)
	fmt.Println("余量:", res)
	return res
}

//OrderOneSeat 对于给定的TripSeries和座位类型，找到一个有效的座位号并下订单
func (s *TripSeries) orderOneSeat(catogory string) error {
	db := common.GetDB()
	tx := db.Begin()
	//1.找到座位信息
	var seats []TripSegment
	string := "SELECT * FROM trip_segment WHERE trip_id = ? AND segment_no between ? AND ? AND seat_catogory = ?"
	err := tx.Set("gorm:query_option", "FOR UPDATE").Raw(string, s.TripID, s.StartStationNo, s.EndStationNo-1, catogory).Find(&seats).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	fmt.Println("座位：", seats)
	//2计算出一个有效的座位号
	validSeatNo, _ := calculasValidSeatNo(seats)
	fmt.Println("选中的座位号", validSeatNo)
	// //3 修改座位信息
	// setZero(seats, validSeatNo)
	// err = tx.Save(seats).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// //4.下订单
	// //UserID，借助中间件.
	// order := Order{UserID: 1, TripID: s.TripID, StartStationNo: s.StartStationNo, EndStationNo: s.EndStationNo, SeatNo: validSeatNo, SeatCatogory: catogory, Date: time.Now(), Status: "未支付"}
	// fmt.Println("订单", order)
	// err = tx.Create(&order).Error
	// //5.commit
	// err = tx.Commit().Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	return nil
}

// func (s *TripSeries) cancleOrder(orderID uint) error {
// 	db := common.GetDB()
// 	tx := db.Begin()
// 	// 0.取得订单信息，并锁住订单表的该行；
// 	order := Order{}
// 	err := tx.Set("gorm:query_option", "FOR UPDATE").First(&order, orderID).Error
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}
// 	// 1.判断该订单状态是否支持退票
// 	if !(order.Status == "未支付" || order.Status == "已支付") {
// 		tx.Rollback()
// 		return errors.New("订单状态错误,不支持退票")
// 	}
// 	// 2.判断下单的用户是不是登录用户本人，借助中间件
// 	// 3.退钱给用户
// 	if order.Status == "已支付" {
// 		fmt.Print("退钱给用户")
// 	}
// 	// 3.修改座位信息
// 	var seats TripSegmentSeats
// 	for i := order.EndStationNo; i < order.StartStationNo; i++ {
// 		string := "SELECT " + order.SeatCatogory + " AS seats FROM trip_segment WHERE trip_id = ? AND segment_no ="
// 		err = tx.Set("gorm:query_option", "FOR UPDATE").Raw(string, s.TripID, i).Find(seats).Error
// 		fmt.Println("座位：", seats)
// 		newSeats, err := modify(seats, order.SeatNo)
// 		fmt.Println("新的座位：", newSeats)
// 		err = tx.Model(&order).Updates(map[string]interface{}{"status": "已退票"}).Error
// 	}
// 	// err = tx.Table("trip_seat").Where("trip_id = ? AND segment between ? AND ? and seat_id = ?", order.TripID, order.StartNo, order.EndNo-1, order.SeatID).Updates(map[string]interface{}{"status": 0}).Error
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}
// 	// 4.更新订单状态
// 	err = tx.Model(&order).Updates(map[string]interface{}{"status": "已退票"}).Error
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}
// 	// 5.commit
// 	err = tx.Commit().Error
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}
// 	return nil
// }
