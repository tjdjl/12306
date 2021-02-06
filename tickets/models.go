package tickets

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

//TrainStaionPair 一个trip的多个连续区间，没有直接对应的数据库表
type TrainStaionPair struct {
	StartStationNo   uint   `gorm:"column:start_station_no"` //StationNo为1代表是该车次的起始站
	EndStationNo     uint   `gorm:"column:end_station_no"`
	StartStationName string `gorm:"column:start_station_name"`
	EndStationName   string `gorm:"column:end_station_name"`
	StartTime        string `gorm:"column:start_time"`
	EndTime          string `gorm:"column:end_time"`
	TrainID          string `gorm:"column:train_id"`
	TrainType        string `gorm:"column:train_type"`
	TrainStationNums uint   `gorm:"column:train_staion_nums"`
}

//FindTrainStaionPairList 根据相应的startCity, endCity,date 条件，返回TripStaionPair列表
func FindTrainStaionPairList(startCity, endCity, date string, isFast bool) ([]TrainStaionPair, error) {
	//取得城市id
	repository := NewTicketRepository()
	startCityID, err := repository.FindCityID(startCity)
	endCityID, err := repository.FindCityID(endCity)
	if err != nil {
		return nil, err

	}
	//拼接成TripStationPairId
	s1 := strconv.FormatUint(uint64(startCityID), 10)
	s2 := strconv.FormatUint(uint64(endCityID), 10)
	id := s1 + "-" + s2 + "-" + "%"
	fmt.Println(id)
	// 判断是否是今天
	now := time.Now()
	nowDate := now.Format("2006-01-02")
	isToday := false
	if date == nowDate {
		isToday = true
	}
	models, err := repository.FindTripStationPair(id, date, isToday, isFast)
	return models, err
}

//TripSegment 对应TripSegment表
type TripSegment struct {
	ID uint `gorm:"primary_key;column:id"` //StationNo为1代表是该车次的起始站
	// SegmentNo    uint    `gorm:"primary_key;column:segment_no"`
	SeatCatogory string  `gorm:"column:seat_catogory"`
	SeatBytes    []uint8 `gorm:"column:seat_bytes"`
}

//TableName 实现TableName接口，以达到结构体和表对应，如果不实现该接口，并未设置全局表名禁用复数，gorm会自动扩展表名为结构体+s
func (TripSegment) TableName() string {
	return "trip_segment"
}

//Trip 一个trip的多个连续区间，没有直接对应的数据库表
type Trip struct {
	TripID string
}

//getRemainSeats 返回座位余量
func (s *Trip) getRemainSeats(startStationNo, endStationNo uint) *map[string]uint { //获取票的座位余量信息
	var resMap map[string]uint
	resMap = make(map[string]uint)
	repository := NewTicketRepository()
	//repository找到对应的TripSegment记录
	seats, err := repository.FindTripSegment(s.TripID, startStationNo, endStationNo)
	if err != nil {
		return &resMap
	}
	fmt.Println("查询到的原始座位位图：")
	for i := 0; i < len(seats); i++ {
		if i > 1 && seats[i].SeatCatogory != seats[i-1].SeatCatogory {
			fmt.Printf("\n")
		}
		if i == 0 || (i > 1 && seats[i].SeatCatogory != seats[i-1].SeatCatogory) {
			fmt.Printf("\n%s", seats[i].SeatCatogory)
		}
		fmt.Printf("%b", seats[i].SeatBytes)
	}
	//对TripSegment记录进行计算
	res := calculasRemainSeats(seats)
	fmt.Printf("\n")
	fmt.Println("计算得到的余量:", res)
	return res
}

//Order 对应订单表
type Order struct {
	ID             uint      `gorm:"primary_key;auto_increment" json:"id"`
	TripID         string    `json:"trip_id"`
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
	//处理候补
	// s.handelCandidate(startStationNo, endStationNo)
	return nil
}

// // 寻找有没有合适的候补，有的话，更改x表、and座位表。
// func (s *Trip) handelCandidate(startStationNo, endStationNo uint) {

// }
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
