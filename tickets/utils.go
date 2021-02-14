package tickets

import (
	"errors"
	"fmt"
)

//calculasRemainSeats 计算余量
func calculasRemainSeats(info []TripSegment) *map[string]uint {
	var resMap map[string]uint
	resMap = make(map[string]uint)
	if len(info) == 0 {
		return &resMap
	}
	var bytesMap map[string][][]uint8
	bytesMap = make(map[string][][]uint8)
	for i := 0; i < len(info); i++ {
		bytesMap[info[i].SeatCatogory] = append(bytesMap[info[i].SeatCatogory], info[i].SeatBytes)
	}
	for k, v := range bytesMap {
		res := calculasRemainCatorySeats(v)
		resMap[k] = res
	}
	return &resMap
}
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
	res := countOne(resBytes)
	// fmt.Println("最终结果：", res)
	return res
}
func calculasValidSeatNo(info []TripSegment) (uint, error) {
	if len(info) == 0 {
		return 0, errors.New("没有余票了")
	}
	resBytes := make([]uint8, len(info[0].SeatBytes))
	copy(resBytes, info[0].SeatBytes)
	fmt.Println(resBytes)
	//对每个区间作与运算
	for i := 1; i < len(info); i++ {
		//计算Seats 经过与运算后的位图
		for j := 0; j < len(resBytes); j++ {
			resBytes[j] = resBytes[j] & info[i].SeatBytes[j]
		}
		// fmt.Println("与运算后bytes", resBytes)
	}
	return FirstOne(resBytes)
}
func setZero(info []TripSegment, validSeatNo uint) {
	if len(info) == 0 {
		return
	}
	index := (validSeatNo - 1) >> 3
	pos := (validSeatNo - 1) % 8
	//对每个区间更改
	for i := 0; i < len(info); i++ {
		info[i].SeatBytes[index] = info[i].SeatBytes[index] & ^(128 >> pos)
	}
}
func setOne(info []TripSegment, validSeatNo uint) {
	if len(info) == 0 {
		return
	}
	index := (validSeatNo - 1) >> 3
	pos := (validSeatNo - 1) % 8
	//对每个区间更改
	for i := 0; i < len(info); i++ {
		info[i].SeatBytes[index] = info[i].SeatBytes[index] | (128 >> pos)
	}
}

//countOne 帮助计算余量 可以优化成hashmap查询
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

//FirstOne
func FirstOne(num []uint8) (uint, error) {
	count := uint(1)
	for i := 0; i < len(num); i++ {
		x := uint8(128)
		temp := num[i]
		for x != 0 {
			if temp&x != 0 {
				return count, nil
			}
			count++
			x = x >> 1
		}
	}
	return count, errors.New("没有余票了")
}

// u_int8_t x = 0x80;
//     //cout<<(int)x<<endl;
//     u_int8_t count=0;
