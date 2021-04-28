package trains

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

//Train 对应trains表
type Train struct {
	ID string
}

//TrainStaion 对应train_station表，车次经过的某个站点
type TrainStaion struct {
	StationNo          uint   `gorm:"column:station_no"`
	StationName        string `gorm:"column:station_name"`
	ArriveTime         string `gorm:"column:arrive_time"`
	DepartureTime      string `gorm:"column:departure_time"`
	TodayArriveTime    string `gorm:"column:today_arrive_time"`
	TodayDepartureTime string `gorm:"column:today_departure_time"`
}

//TableName 指定表名
func (TrainStaion) TableName() string {
	return "train_station"
}

//TrainStaionPair 对应TrainStaionPair表
type TrainStaionPair struct {
	StartStationNo   uint   `gorm:"column:start_station_no"` //StationNo为1代表是该车次的起始站
	EndStationNo     uint   `gorm:"column:end_station_no"`
	StartStationName string `gorm:"column:start_station_name"`
	EndStationName   string `gorm:"column:end_station_name"`
	StartTime        string `gorm:"column:start_time"`
	EndTime          string `gorm:"column:end_time"`
	//与trains表关联
	TrainID          string `gorm:"column:train_id"`
	TrainType        string `gorm:"column:train_type"`
	TrainStationNums uint   `gorm:"column:train_staion_nums"`
}

//Trip 对应trip表
type Trip struct {
	TripID string
}

//TripSegment 对应TripSegment表
type TripSegment struct {
	ID uint `gorm:"primary_key;column:id"` //StationNo为1代表是该车次的起始站
	// SegmentNo    uint    `gorm:"primary_key;column:segment_no"`
	SeatCatogory string  `gorm:"column:seat_catogory"`
	SeatBytes    []uint8 `gorm:"column:seat_bytes"`
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

//TableName 实现TableName接口，以达到结构体和表对应，如果不实现该接口，并未设置全局表名禁用复数，gorm会自动扩展表名为结构体+s
func (TripSegment) TableName() string {
	return "trip_segment"
}

//getTrainStaions 获取车次经过的所有站点
func (s *Train) getTrainStaions() ([]TrainStaion, error) {
	repository := NewTicketRepository()
	models, err := repository.FindTrainStaions(s.ID)
	return models, err
}

//ListTrainStaionPair 根据相应的startCity, endCity,date,isFast 条件，返回TrainStaionPair列表
func ListTrainStaionPair(startCity, endCity, date string, isFast bool) ([]TrainStaionPair, error) {
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
	// fmt.Println(id)
	// 判断是否是今天
	nowDate := time.Now().Format("2006-01-02")
	isToday := date == nowDate
	models, err := repository.FindTrainStationPairs(id, date, isToday, isFast)
	return models, err
}

//getRemainSeats 返回座位余量
func (s *Trip) getRemainSeats(startStationNo, endStationNo uint) *map[string]uint {
	var resMap map[string]uint
	resMap = make(map[string]uint)
	repository := NewTicketRepository()
	//repository找到对应的TripSegment记录
	seats, err := repository.FindTripSegments(s.TripID, startStationNo, endStationNo)
	if err != nil {
		return &resMap
	}
	// fmt.Println("查询到的原始座位位图：")
	// for i := 0; i < len(seats); i++ {
	// 	if i > 1 && seats[i].SeatCatogory != seats[i-1].SeatCatogory {
	// 		fmt.Printf("\n")
	// 	}
	// 	if i == 0 || (i > 1 && seats[i].SeatCatogory != seats[i-1].SeatCatogory) {
	// 		fmt.Printf("\n%s", seats[i].SeatCatogory)
	// 	}
	// 	fmt.Printf("%b", seats[i].SeatBytes)
	// }
	//对TripSegment记录进行计算
	res := calculasRemainSeats(seats)
	// fmt.Printf("计算得到的余量d%\n", res)
	return res
}

//OrderOneSeat 对于给定的TripStartNoAndEndNo和座位类型，找到一个有效的座位号并下订单
func (s *Trip) orderOneSeat(startStationNo, endStationNo uint, catogory string) error {
	repository := NewTicketRepositoryTX()
	//1.repository找到座位信息
	fmt.Println("查询前", time.Now())
	seats, err := repository.FindTripSegments(s.TripID, startStationNo, endStationNo, catogory)
	fmt.Println("查询后", time.Now())

	if err != nil {
		return err
	}
	// fmt.Println("查询到的原始座位位图：")
	// for i := 0; i < len(seats); i++ {
	// 	fmt.Printf("	区间%d:%8b\n", i+1, seats[i].SeatBytes)
	// }
	//2计算出一个有效的座位号
	validSeatNo, err := calculasValidSeatNo(seats)
	if err != nil {
		repository.Rollback()
		return err
	}
	// fmt.Println("经过计算选中的座位号", validSeatNo)
	//3 修改座位信息
	setZero(seats, validSeatNo)
	// fmt.Println("修改后的座位位图：")
	// for i := 0; i < len(seats); i++ {
	// 	fmt.Printf("	区间%d:%8b\n", i+1, seats[i].SeatBytes)
	// }

	//4.repository写回修改座位信息
	fmt.Println("插入后", time.Now())
	err = repository.UpdateTripSegment(seats)
	fmt.Println("插入后", time.Now())

	if err != nil {
		repository.Rollback()
		return err
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	//5.repository下订单
	//UserID，借助中间件.
	order := Order{UserID: 1, TripID: s.TripID, StartStationNo: startStationNo, EndStationNo: endStationNo, SeatNo: validSeatNo, SeatCatogory: catogory, Date: time.Now(), Status: "未支付"}
	// fmt.Println("生成订单:", order)
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
		// fmt.Print("退钱给用户")
	}

	// 3.修改座位信息
	seats, err := repository.FindTripSegments(order.TripID, order.StartStationNo, order.EndStationNo, order.SeatCatogory)
	if err != nil {
		repository.Rollback()
		return err
	}
	// fmt.Println("座位修改前：")
	// for i := 0; i < len(seats); i++ {
	// 	fmt.Printf("	区间%d:%8b\n", i+1, seats[i].SeatBytes)
	// }
	if len(seats) == 0 {
		return errors.New("座位信息错误")
	}
	setOne(seats, order.SeatNo)
	// fmt.Println("座位修改后：")
	// for i := 0; i < len(seats); i++ {
	// 	fmt.Printf("	区间%d:%8b\n", i+1, seats[i].SeatBytes)
	// }
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

func (s *Trip) changeOrder(orderID uint, startStationNo, endStationNo uint, catogory string) error {
	repository := NewTicketRepositoryTX()
	// 1.取得合法订单信息
	userID := uint(1)
	order, err := repository.FindValidOrder(orderID, userID)
	if err != nil {
		return err
	}
	//2.找到并修改新座位信息
	newSeats, err := repository.FindTripSegments(s.TripID, startStationNo, endStationNo, catogory)
	if err != nil {
		return err
	}
	validSeatNo, err := calculasValidSeatNo(newSeats)
	setZero(newSeats, validSeatNo)
	err = repository.UpdateTripSegment(newSeats)

	// 3.修改旧座位信息
	seats, err := repository.FindTripSegments(order.TripID, order.StartStationNo, order.EndStationNo, order.SeatCatogory)
	if err != nil {
		repository.Rollback()
		return err
	}
	if len(seats) == 0 {
		return errors.New("座位信息错误")
	}
	setOne(seats, order.SeatNo)
	err = repository.UpdateTripSegment(seats)
	if err != nil {
		repository.Rollback()
		return err
	}
	// 4.更新订单状态
	newMap := map[string]interface{}{"status": "已改票", "trip_id": s.TripID, "start_station_no": startStationNo, "end_station_no": endStationNo, "seat_catogory": catogory, "seat_no": validSeatNo}
	err = repository.UpdateOrder(&order, newMap)
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

// // 寻找有没有合适的候补，有的话，更改x表、and座位表。
// func (s *Trip) handelCandidate(startStationNo, endStationNo uint) {

// }
