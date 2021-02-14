package tickets

import (
	"errors"
	"fmt"
	"time"
)

//TripStaionPair 一个trip的多个连续区间，没有直接对应的数据库表
type TripStaionPair struct {
	StartStationNo   uint   `gorm:"column:start_station_no"` //StationNo为1代表是该车次的起始站
	EndStationNo     uint   `gorm:"column:end_station_no"`
	StartStationName string `gorm:"column:start_station_name"`
	EndStationName   string `gorm:"column:end_station_name"`
	StartTime        string `gorm:"column:start_time"`
	EndTime          string `gorm:"column:end_time"`
	TrainNumber      string `gorm:"column:train_number"`
	TrainType        string `gorm:"column:train_type"`
	TrainStationNums uint   `gorm:"column:train_staion_nums"`
	TripID           uint   `gorm:"column:trip_id"`
}

// // Train 对应数据库表中的train
// type Train struct {
// 	Catogory    string `gorm:"column:catogory"`
// 	TrainNumber string `gorm:"column:train_number"`
// 	Length      uint   `gorm:"column:length"`
// }

//FindTripStaionPairList 根据相应的startCity, endCity,date 条件，返回TripStaionPair列表
func FindTripStaionPairList(startCity, endCity, date string) ([]TripStaionPair, error) {
	repository := NewTicketRepository()
	models, err := repository.FindTripStationPair(startCity, endCity, date)
	return models, err
}

//FindFastTripStaionPairList 根据相应的startCity, endCity, date条件，返回高铁快车的TripStaionPair列表
func FindFastTripStaionPairList(startCity, endCity, date string) ([]TripStaionPair, error) {
	repository := NewTicketRepository()
	models, err := repository.FindFastTripStationPair(startCity, endCity, date)
	return models, err
}

//TripSegment 对应TripSegment表
type TripSegment struct {
	ID           uint    `gorm:"primary_key;column:id"` //StationNo为1代表是该车次的起始站
	SegmentNo    uint    `gorm:"primary_key;column:segment_no"`
	SeatCatogory string  `gorm:"column:seat_catogory"`
	SeatBytes    []uint8 `gorm:"column:seat_bytes"`
}

//TableName 实现TableName接口，以达到结构体和表对应，如果不实现该接口，并未设置全局表名禁用复数，gorm会自动扩展表名为结构体+s
func (TripSegment) TableName() string {
	return "trip_segment"
}

//Trip 一个trip的多个连续区间，没有直接对应的数据库表
type Trip struct {
	TripID uint `gorm:"column:trip_id"`
}

//getRemainSeats 对于给定的TripStartNoAndEndNo，根据座位类型,返回座位余量
func (s *Trip) getRemainSeats(startStationNo, endStationNo uint, catogory string) uint { //获取票的座位余量信息
	repository := NewTicketRepository()
	seats, err := repository.FindTripSegment(s.TripID, startStationNo, endStationNo, catogory)
	if err != nil {
		return 0
	}
	fmt.Println("查询到的原始座位位图：")
	for i := 0; i < len(seats); i++ {
		fmt.Printf("%b", seats[i].SeatBytes)
	}
	res := calculasRemainSeats(seats)
	fmt.Printf("\n")
	fmt.Println("计算得到的余量:", res)
	return res
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
func (s *Trip) orderOneSeat(startStationNo, endStationNo uint, catogory string) error {
	repository := NewTicketRepositoryTX()
	//1.repository找到座位信息
	seats, err := repository.FindTripSegment(s.TripID, startStationNo, endStationNo, catogory)
	if err != nil {
		return err
	}
	fmt.Println("查询到的原始座位位图：")
	for i := 0; i < len(seats); i++ {
		fmt.Printf("%b", seats[i].SeatBytes)
	}
	fmt.Printf("\n")
	//2计算出一个有效的座位号
	validSeatNo, err := calculasValidSeatNo(seats)
	if err != nil {
		repository.Rollback()
		return err
	}
	fmt.Println("经过计算选中的座位号", validSeatNo)
	//3 修改座位信息
	setZero(seats, validSeatNo)
	fmt.Println("修改后的座位位图：")
	for i := 0; i < len(seats); i++ {
		fmt.Printf("%b", seats[i].SeatBytes)
	}
	fmt.Printf("\n")
	//4.repository写回修改座位信息
	err = repository.UpdateTripSegment(seats)
	if err != nil {
		repository.Rollback()
		return err
	}
	//5.repository下订单
	//UserID，借助中间件.
	order := Order{UserID: 1, TripID: s.TripID, StartStationNo: startStationNo, EndStationNo: endStationNo, SeatNo: validSeatNo, SeatCatogory: catogory, Date: time.Now(), Status: "未支付"}
	fmt.Println("订单", order)
	err = repository.CreateOrder(&order)
	if err != nil {
		repository.Rollback()
		return err
	}
	//6.commit
	err = repository.Commit()
	if err != nil {
		repository.Rollback()
		return err
	}
	return nil
}

func (s *Trip) cancleOrder(orderID uint) error {
	repository := NewTicketRepositoryTX()
	// repository.
	// 1.取得合法订单信息
	userID := uint(1)
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
	seats, err := repository.FindTripSegment(order.TripID, order.StartStationNo, order.EndStationNo, order.SeatCatogory)
	if err != nil {
		repository.Rollback()
		return err
	}
	fmt.Println("座位修改前：")
	for i := 0; i < len(seats); i++ {
		fmt.Printf("%b", seats[i].SeatBytes)
	}
	fmt.Printf("\n")
	setOne(seats, order.SeatNo)
	fmt.Println("座位修改后：")
	for i := 0; i < len(seats); i++ {
		fmt.Printf("%b", seats[i].SeatBytes)
	}
	fmt.Printf("\n")
	if len(seats) == 0 {
		return errors.New("wrong")
	}
	err = repository.UpdateTripSegment(seats)
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

func (s* Trip) changeOrder(orderID uint,tripID uint, catogory string) error{
	// 1.取得合法订单信息
	repository := NewTicketRepositoryTX()
	// 2.取得订单起始到达城市
	userID := uint(1)
	oldOrder, err := repository.FindValidOrder(orderID, userID)
	if err!=nil {
		return err
	}
	//3.repository找到座位信息
	newSeats, err := repository.FindTripSegment(tripID, oldOrder.startStationNo, oldOrder.endStationNo, catogory)
	if err != nil {
		return err
	}
	//4.计算出一个有效的座位号
	newValidSeatNo, err := calculasValidSeatNo(newSeats)
	if err != nil {
		repository.Rollback()
		return err
	}
	fmt.Println("经过计算选中的座位号", newValidSeatNo)
	//5.修改座位信息
	setZero(newSeats, newValidSeatNo)
	//6.repository写回修改座位信息
	err = repository.UpdateTripSegment(newSeats)
	if err != nil {
		repository.Rollback()
		return err
	}
	//7.repository下订单
	//UserID，借助中间件.
	newOrder := Order{UserID: 1, TripID: tripID, StartStationNo: oldOrder.startStationNo, EndStationNo: oldOrder.endStationNo, SeatNo: newValidSeatNo, SeatCatogory: catogory, Date: time.Now(), Status: "未支付"}
	fmt.Println("订单", newOrder)
	err = repository.CreateOrder(&newOrder)
	if err != nil {
		repository.Rollback()
		return err
	}
	//8.commit
	err = repository.Commit()
	if err != nil {
		repository.Rollback()
		return err
	}
	// 9.退钱给用户
	if oldOrder.Status == "已支付" {
		fmt.Print("退钱给用户")
	}

	// 10.修改座位信息
	oldSeats, err := repository.FindTripSegment(oldOrder.TripID, oldOrder.StartStationNo, oldOrder.EndStationNo, oldOrder.SeatCatogory)
	if err != nil {
		repository.Rollback()
		return err
	}
	setOne(oldSeats, oldOrder.SeatNo)
	if len(oldSeats) == 0 {
		return errors.New("wrong")
	}
	err = repository.UpdateTripSegment(oldSeats)
	if err != nil {
		repository.Rollback()
		return err
	}
	// 11.更新订单状态
	err = repository.UpdateOrderStatus(&oldOrder, "已改签")
	if err != nil {
		repository.Rollback()
		return err
	}
	// 12.commit
	err = repository.Commit()
	if err != nil {
		repository.Rollback()
		return err
	}
	return nil
}
