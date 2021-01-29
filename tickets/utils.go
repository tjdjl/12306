package tickets

import (
	"errors"
	"fmt"
)

//calculasRemainSeats 帮助计算余量，将TripSegmentbytes转化为RemainSeats
func calculasRemainSeats(info []TripSegment) uint {
	if len(info) == 0 {
		return 0
	}
	resBytes := info[0].SeatBytes
	//对每个区间作与运算
	for i := 1; i < len(info); i++ {
		for j := 0; j < len(resBytes); j++ {
			resBytes[j] = resBytes[j] & info[i].SeatBytes[j]
		}
		fmt.Println("与运算后bytes", resBytes)
	}
	res := countOne(resBytes)
	fmt.Println("最终结果：", res)
	return res
}
func calculasValidSeatNo(info []TripSegment) (uint, error) {
	if len(info) == 0 {
		return 0, errors.New("没有余票了")
	}
	resBytes := info[0].SeatBytes
	//对每个区间作与运算
	for i := 1; i < len(info); i++ {
		//计算Seats 经过与运算后的位图
		for j := 0; j < len(resBytes); j++ {
			resBytes[j] = resBytes[j] & info[i].SeatBytes[j]
		}
		fmt.Println("与运算后bytes", resBytes)
	}
	return FirstOne(resBytes)
}
func setZero(info []TripSegment, validSeatNo uint) {
	if len(info) == 0 {
		return
	}
	index := validSeatNo - 1>>3 //100号
	pos := (validSeatNo - 1) % 8
	//对每个区间更改
	for i := 0; i < len(info); i++ {
		info[i].SeatBytes[index] = info[i].SeatBytes[index] & ^(1 << pos)
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
		temp := num[i]
		for j := 7; j <= 0; j-- {
			if (1<<j)&temp != 0 {
				return count, nil
			}
			count++
		}
	}
	return count, errors.New("没有余票了")
}