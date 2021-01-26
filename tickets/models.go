package tickets

import (
	"fmt"

	"12306.com/12306/common"

	_ "fmt"
)

//Ticket 这里的票对应，一个车次的多个连续区间，并不存在数据库表中
type Ticket struct {
	TripID         uint   `gorm:"column:trip_id"`
	StartStationNo uint   `gorm:"column:start_station_no"`
	EndStationNo   uint   `gorm:"column:end_station_no"`
	StartStation   string `gorm:"column:start_station"`
	EndStation     string `gorm:"column:end_station"`
}

//Trip 对应数据库表中的车次
type Trip struct {
	ID     uint   `gorm:"primary_key"`
	TripID string `json:"trip_id"`
	TripNo string `json:"Trip_No"`
}

//FindTicketList 根据相应的startCity, endCity,date 条件返回ticket列表
func FindTicketList(startCity, endCity, date string) ([]Ticket, error) { //从数据库找到所有符合条件的车次
	db := common.GetDB()
	var models []Ticket
	err := db.Raw("SELECT A.trip_id AS trip_id ,A.sequence AS start_station_no,B.sequence AS end_station_no FROM (SELECT trip_id,station_name,sequence FROM trip_station WHERE station_name =? AND date(start_time)=? ) A, (SELECT trip_id,station_name,sequence FROM trip_station WHERE station_name =? AND date(start_time)=? ) B WHERE A.sequence < B.sequence AND A.trip_id = B.trip_id ", startCity, date, endCity, date).Find(&models).Error
	return models, err
}

//FindHishSpeedTicketList 根据相应的startCity, endCity, date条件返回高铁快车的ticket列表
func FindHishSpeedTicketList(startCity, endCity, date string) ([]Ticket, error) { //从数据库找到所有符合条件的车次
	db := common.GetDB()
	var models []Ticket
	err := db.Raw("SELECT A.trip_id ,A.sequence ,B.sequence FROM (SELECT trip_id,station_name,sequence FROM trip_station WHERE station_name =? AND date(start_time)=? ) A, (SELECT trip_id,station_name,sequence FROM trip_station WHERE station_name =? AND date(start_time)=? ) B WHERE A.sequence < B.sequence AND A.trip_id = B.trip_id ", startCity, date, endCity, date).Find(&models).Error
	//返回
	return models, err
}

// //getTrainDetail 对于给定的Ticket，查找对应的车次信息
// func (self *Ticket) getTrainDetail() { //获取车次对应的列车的信息
// }

//RemainSeats 该结构体用于表示，对应给定的ticket结构体，它的各个座位类型的余票数
type RemainSeats struct {
	BusinessSeats   uint `json:"businessSeatsNumber"`
	FirstSeats      uint `json:"firstSeatsNumber"`
	SecondSeats     uint `json:"secondSeatsNumber"`
	HardSeats       uint `json:"hardSeatsNumber"`
	HardBerth       uint `json:"hardBerthNumber"`
	SoftBerth       uint `json:"softBerthNumber"`
	SeniorSoftBerth uint `json:"seniorSoftBerthNumber"`
}

//SeatsBytes 该结构体用于表示，对应给定的ticket结构体，各个座位类型的原始二进制序列,直接对应数据库的二进制序列
type SeatsBytes struct {
	TripID          uint    `gorm:"column:trip_id"`
	SegmentNo       uint    `gorm:"column:segment_no"`
	BusinessSeats   []uint8 `gorm:"column:business_seats"`
	FirstSeats      []uint8 `gorm:"column:first_seats"`
	SecondSeats     []uint8 `gorm:"column:second_seats"`
	HardSeats       []uint8 `gorm:"column:hard_seats"`
	HardBerth       []uint8 `gorm:"column:hard_berth"`
	SoftBerth       []uint8 `gorm:"column:soft_berth"`
	SeniorSoftBerth []uint8 `gorm:"column:senior_soft_berth"`
}

//getSeatDetail 对于给定的Ticket，查找对应的座位余量信息
func (s *Ticket) getSeatDetail() RemainSeats { //获取票的座位余量信息
	db := common.GetDB()
	var info []SeatsBytes
	// 批量读取
	db.Raw("SELECT * FROM trip_segment WHERE trip_id = ? AND segment_no between ? AND ? ", s.TripID, s.StartStationNo, s.EndStationNo-1).Find(&info)
	fmt.Println("座位：", info)
	res := calculasRemainSeats(info)
	return res
}

//calculasRemainSeats 帮助计算余量，将TripSegmentbytes转化为RemainSeats
func calculasRemainSeats(info []SeatsBytes) RemainSeats {
	var resBytes SeatsBytes = info[0]
	//对每个区间作与运算
	for i := 1; i < len(info); i++ {
		//计算BusinessSeats 经过与运算后的位图
		for j := 0; j < len(resBytes.BusinessSeats); j++ {
			resBytes.BusinessSeats[j] = resBytes.BusinessSeats[j] & info[i].BusinessSeats[j]
		}
		for j := 0; j < len(resBytes.FirstSeats); j++ {
			resBytes.FirstSeats[j] = resBytes.FirstSeats[j] & info[i].FirstSeats[j]
		}
		for j := 0; j < len(resBytes.SecondSeats); j++ {
			resBytes.SecondSeats[j] = resBytes.SecondSeats[j] & info[i].SecondSeats[j]
		}
		for j := 0; j < len(resBytes.HardSeats); j++ {
			resBytes.HardSeats[j] = resBytes.HardSeats[j] & info[i].HardSeats[j]
		}
		for j := 0; j < len(resBytes.HardBerth); j++ {
			resBytes.HardBerth[j] = resBytes.HardBerth[j] & info[i].HardBerth[j]
		}
		for j := 0; j < len(resBytes.SoftBerth); j++ {
			resBytes.SoftBerth[j] = resBytes.SoftBerth[j] & info[i].SoftBerth[j]
		}
		for j := 0; j < len(resBytes.SeniorSoftBerth); j++ {
			resBytes.SeniorSoftBerth[j] = resBytes.SeniorSoftBerth[j] & info[i].SeniorSoftBerth[j]
		}
		fmt.Println("与运算后bytes", resBytes)
	}
	var res RemainSeats
	res.BusinessSeats = countOne(resBytes.BusinessSeats)
	res.FirstSeats = countOne(resBytes.FirstSeats)
	res.SecondSeats = countOne(resBytes.SecondSeats)
	res.HardSeats = countOne(resBytes.HardSeats)
	res.HardBerth = countOne(resBytes.HardBerth)
	res.SoftBerth = countOne(resBytes.SoftBerth)
	res.SeniorSoftBerth = countOne(resBytes.SeniorSoftBerth)
	fmt.Println("最终结果：", res)
	return res
}

//countOne 帮助计算余量
func countOne(num []uint8) uint {
	count := uint(0)
	for i := 0; i < len(num); i++ {
		temp := num[i]
		for temp != 0 {
			if temp%2 != 0 {
				count++
			}
			temp /= 2
		}
	}
	return count
}
