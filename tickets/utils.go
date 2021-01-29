package tickets

import (
	"errors"
	"fmt"
)

//calculasRemainSeats 帮助计算余量，将TripSegmentbytes转化为RemainSeats
func calculasRemainSeats(info []TripSegment) RemainSeats {
	var resBytes TripSegment = info[0]
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

func calculasValidSeatNo(info []TripSegmentSeats) (uint, error) {
	resBytes := info[0]
	//对每个区间作与运算
	for i := 1; i < len(info); i++ {
		//计算Seats 经过与运算后的位图
		for j := 0; j < len(resBytes.Seats); j++ {
			resBytes.Seats[j] = resBytes.Seats[j] & info[i].Seats[j]
		}
		fmt.Println("与运算后bytes", resBytes)
	}
	return FirstOne(resBytes.Seats)
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
	for i := len(num) - 1; i >= 0; i-- {
		temp := num[i]
		for i := 0; i < 8; i++ {
			if temp%2 != 0 {
				return count, nil
			}
			temp /= 2
			count++
		}
	}
	return count, errors.New("没有余票了")
}
