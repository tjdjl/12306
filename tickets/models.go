package tickets

import (
	"errors"
	"fmt"
	"time"
)

//TripStartNoAndEndNo 一个trip的多个连续区间，没有直接对应的数据库表
type TripStartNoAndEndNo struct {
	TripID         uint `gorm:"column:trip_id"`
	StartStationNo uint `gorm:"column:start_station_no"` //StationNo为1代表是该车次的起始站
	EndStationNo   uint `gorm:"column:end_station_no"`
}

//FindTripStartAndEndList 根据相应的startCity, endCity,date 条件，返回TripSeries列表
func FindTripStartAndEndList(startCity, endCity, date string) ([]TripStartNoAndEndNo, error) { //从数据库找到所有符合条件的车次
	repository := NewTicketRepository()
	models, err := repository.FindTripStartNoAndEndNo(startCity, endCity, date)
	return models, err
}

//FindHighSpeedTripStartAndEndList 根据相应的startCity, endCity, date条件，返回高铁快车的TripSeries列表
func FindHighSpeedTripStartAndEndList(startCity, endCity, date string) ([]TripStartNoAndEndNo, error) { //从数据库找到所有符合条件的车次
	repository := NewTicketRepository()
	models, err := repository.FindHighSpeedTripStartNoAndEndNo(startCity, endCity, date)
	return models, err
}

// Train 对应数据库表中的train
type Train struct {
	Catogory    string `gorm:"column:catogory"`
	TrainNumber string `gorm:"column:train_number"`
	Length      uint   `gorm:"column:length"`
}

// StartStaionDetail
type StaionDetail struct {
	StationName string    `gorm:"column:station_name"`
	StationTime time.Time `gorm:"column:station_time"`
}

//TripSegment 对应TripSegment表
type TripSegment struct {
	ID           uint    `gorm:"column:id"`
	TripID       uint    `gorm:"column:trip_id"`
	SegmentNo    uint    `gorm:"column:segment_no"`
	SeatCatogory string  `gorm:"column:seat_catogory"`
	SeatBytes    []uint8 `gorm:"column:seat_bytes"`
}

//TableName 实现TableName接口，以达到结构体和表对应，如果不实现该接口，并未设置全局表名禁用复数，gorm会自动扩展表名为结构体+s
func (TripSegment) TableName() string {
	return "trip_segment"
}

//getRemainSeats 对于给定的TripStartNoAndEndNo，根据座位类型,返回座位余量
func (s *TripStartNoAndEndNo) getRemainSeats(catogory string) uint { //获取票的座位余量信息
	repository := NewTicketRepository()
	info, err := repository.FindTripSegment(s.TripID, s.StartStationNo, s.EndStationNo, catogory)
	if err != nil {
		return 0
	}
	fmt.Println("座位：", info)
	res := calculasRemainSeats(info)
	fmt.Println("余量:", res)
	return res
}

//getRemainSeats 对于给定的TripStartNoAndEndNo，根据座位类型,返回座位余量
func (s *TripStartNoAndEndNo) getTrainDetail() Train { //获取票的座位余量信息
	repository := NewTicketRepository()
	info, err := repository.FindTrain(s.TripID)
	if err != nil {
		return info
	}
	return info
}

//getRemainSeats 对于给定的TripStartNoAndEndNo，根据座位类型,返回座位余量
func (s *TripStartNoAndEndNo) getStationDetail() [2]StaionDetail { //获取票的座位余量信息
	repository := NewTicketRepository()
	var info [2]StaionDetail
	var err error
	info[0], err = repository.FindStartStaionDetail(s.TripID, s.StartStationNo)
	info[1], err = repository.FindEndStaionDetail(s.TripID, s.EndStationNo)
	if err != nil {
		return info
	}
	return info
}

//Order 对应订单表
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

//OrderOneSeat 对于给定的TripStartNoAndEndNo和座位类型，找到一个有效的座位号并下订单
func (s *TripStartNoAndEndNo) orderOneSeat(catogory string) error {
	repository := NewTicketRepositoryTX()
	//1.找到座位信息
	seats, err := repository.FindTripSegment(s.TripID, s.StartStationNo, s.EndStationNo, catogory)
	if err != nil {
		return err
	}
	fmt.Println("座位：", seats)
	//2计算出一个有效的座位号
	validSeatNo, err := calculasValidSeatNo(seats)
	if err != nil {
		repository.Rollback()
		return err
	}
	fmt.Println("选中的座位号", validSeatNo)
	//3 修改座位信息
	setZero(seats, validSeatNo)
	fmt.Println("座位修改后：", seats)
	err = repository.UpdateTripSegment(seats[0])
	if err != nil {
		repository.Rollback()
		return err
	}
	//4.下订单
	//UserID，借助中间件.
	order := Order{UserID: 1, TripID: s.TripID, StartStationNo: s.StartStationNo, EndStationNo: s.EndStationNo, SeatNo: validSeatNo, SeatCatogory: catogory, Date: time.Now(), Status: "未支付"}
	fmt.Println("订单", order)
	err = repository.CreateOrder(&order)
	if err != nil {
		repository.Rollback()
		return err
	}
	//5.commit
	err = repository.Commit()
	if err != nil {
		repository.Rollback()
		return err
	}
	return nil
}

func (s *TripStartNoAndEndNo) cancleOrder(orderID uint) error {
	repository := NewTicketRepositoryTX()
	// repository.

	userID := uint(1)
	// 1.取得合法订单信息
	order, err := repository.FindValidOrder(orderID, userID)
	if err != nil {
		repository.Rollback()
		return err
	}
	// 2.退钱给用户
	if order.Status == "已支付" {
		fmt.Print("退钱给用户")
	}

	// 3.修改座位信息
	seats, err := repository.FindTripSegment(order.TripID, order.StartStationNo, order.EndStationNo-1, order.SeatCatogory)
	if err != nil {
		repository.Rollback()
		return err
	}
	fmt.Println("座位修改前：", seats)
	setOne(seats, order.SeatNo)
	fmt.Println("座位修改后：", seats)
	if len(seats) == 0 {
		return errors.New("wrong")
	}
	err = repository.UpdateTripSegment(seats[0])
	if err != nil {
		repository.Rollback()
		return err
	}
	// 4.更新订单状态
	err = repository.UpdateOrderStatus(&order, "已退票")
	if err != nil {
		repository.Rollback()
		return err
	}
	// 5.commit
	err = repository.Commit()
	if err != nil {
		repository.Rollback()
		return err
	}
	return nil
}
