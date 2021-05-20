package trains

import (
	"errors"
	"fmt"

	"12306.com/12306/common/bytes"
)

type TripSegments struct {
	info []TripSegment
}

func (s *TripSegments) printBytes() {
	for i := 0; i < len(s.info); i++ {
		if i > 1 && s.info[i].SeatCatogory != s.info[i-1].SeatCatogory {
			fmt.Printf("\n")
		}
		if i == 0 || (i > 1 && s.info[i].SeatCatogory != s.info[i-1].SeatCatogory) {
			fmt.Printf("\n%s", s.info[i].SeatCatogory)
		}
		fmt.Printf("%b", s.info[i].SeatBytes)
	}
}
func (s *TripSegments) printBytes1() {
	for i := 0; i < len(s.info); i++ {
		fmt.Printf("	区间%d:%8b\n", i+1, s.info[i].SeatBytes)
	}
}
func (s *TripSegments) unionBytes() []uint8 {
	resBytes := make([]uint8, len(s.info[0].SeatBytes))
	copy(resBytes, s.info[0].SeatBytes)
	//对每个区间作与运算
	for i := 1; i < len(s.info); i++ {
		//计算Seats 经过与运算后的位图
		for j := 0; j < len(resBytes); j++ {
			resBytes[j] = resBytes[j] & s.info[i].SeatBytes[j]
		}
	}
	return resBytes
}

//calculasRemainSeats 计算余量
func (s *TripSegments) calculasRemainSeats() *map[string]uint {
	var resMap map[string]uint
	resMap = make(map[string]uint)
	if len(s.info) == 0 {
		return &resMap
	}
	var bytesMap map[string][][]uint8
	bytesMap = make(map[string][][]uint8)
	for i := 0; i < len(s.info); i++ {
		bytesMap[s.info[i].SeatCatogory] = append(bytesMap[s.info[i].SeatCatogory], s.info[i].SeatBytes)
	}
	for k, v := range bytesMap {
		res := calculasRemainCatorySeats(v)
		resMap[k] = res
	}
	return &resMap
}

func (s *TripSegments) calculasOneValidSeatNo() (int64, error) {
	if len(s.info) == 0 {
		return 0, errors.New("没有余票了")
	}
	resBytes := s.unionBytes()
	return bytes.FirstOne(resBytes)
}
func (s *TripSegments) calculasValidSeatNos(count int32) ([]int64, bool) {
	if len(s.info) == 0 {
		return nil, false
	}
	resBytes := s.unionBytes()
	return bytes.FindCount(resBytes, count)
}
func (s *TripSegments) discountSeats(validSeatNos []int64) {
	for i := 0; i < len(validSeatNos); i++ {
		s.discountOneSeat(validSeatNos[i])
	}
}
func (s *TripSegments) discountOneSeat(validSeatNo int64) {
	index := (validSeatNo - 1) >> 3
	pos := (validSeatNo - 1) % 8
	//对每个区间更改
	for i := 0; i < len(s.info); i++ {
		s.info[i].SeatBytes[index] = s.info[i].SeatBytes[index] & ^(128 >> pos)
	}
}
func (s *TripSegments) addOneSeat(validSeatNo int64) {
	if len(s.info) == 0 {
		return
	}
	index := (validSeatNo - 1) >> 3
	pos := (validSeatNo - 1) % 8
	//对每个区间更改
	for i := 0; i < len(s.info); i++ {
		s.info[i].SeatBytes[index] = s.info[i].SeatBytes[index] | (128 >> pos)
	}
}

// u_int8_t x = 0x80;
//     //cout<<(int)x<<endl;
//     u_int8_t count=0;
func calculasRemainCatorySeats(info [][]uint8) uint {
	if len(info) == 0 {
		return 0
	}
	resBytes := info[0]
	for i := 1; i < len(info); i++ {
		for j := 0; j < len(resBytes); j++ {
			resBytes[j] = resBytes[j] & info[i][j]
		}
		// fmt.Println("与运算后bytes", resBytes)
	}
	res := bytes.CountOne(resBytes)
	return res
}

//OrderOneSeat 对于给定的TripStartNoAndEndNo和座位类型，找到一个有效的座位号并下订单
// func (s *Trip) orderOneSeat(startStationNo, endStationNo uint, catogory string) error {
// 	repository := NewTicketRepositoryTX()
// 	//1.repository找到座位信息
// 	fmt.Println("查询前", time.Now())
// 	seats, err := repository.FindTripSegments(s.TripID, startStationNo, endStationNo, catogory)
// 	fmt.Println("查询后", time.Now())

// 	if err != nil {
// 		return err
// 	}
// 	// fmt.Println("查询到的原始座位位图：")
// 	// for i := 0; i < len(seats); i++ {
// 	// 	fmt.Printf("	区间%d:%8b\n", i+1, seats[i].SeatBytes)
// 	// }
// 	//2计算出一个有效的座位号
// 	tripSegments := TripSegments{seats}
// 	validSeatNo, err := tripSegments.calculasOneValidSeatNo()
// 	if err != nil {
// 		repository.Rollback()
// 		return err
// 	}
// 	// fmt.Println("经过计算选中的座位号", validSeatNo)
// 	//3 修改座位信息
// 	setZero(seats, validSeatNo)
// 	// fmt.Println("修改后的座位位图：")
// 	// for i := 0; i < len(seats); i++ {
// 	// 	fmt.Printf("	区间%d:%8b\n", i+1, seats[i].SeatBytes)
// 	// }

// 	//4.repository写回修改座位信息
// 	fmt.Println("插入后", time.Now())
// 	err = repository.UpdateTripSegment(seats)
// 	fmt.Println("插入后", time.Now())

// 	if err != nil {
// 		repository.Rollback()
// 		return err
// 	}
// 	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
// 	//5.repository下订单
// 	//UserID，借助中间件.
// 	order := Order{UserID: 1, TripID: s.TripID, StartStationNo: startStationNo, EndStationNo: endStationNo, SeatNo: validSeatNo, SeatCatogory: catogory, Date: time.Now(), Status: "未支付"}
// 	// fmt.Println("生成订单:", order)
// 	err = repository.CreateOrder(&order)
// 	if err != nil {
// 		repository.Rollback()
// 		return err
// 	}
// 	//6.commit
// 	err = repository.Commit()
// 	if err != nil {
// 		repository.Rollback()
// 		return err
// 	}
// 	return nil
// }

// if err != nil {
// 	c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": err.Error()})
// 	return
// }
