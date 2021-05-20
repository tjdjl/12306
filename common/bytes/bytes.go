package bytes

import "errors"

//FirstOne
func FirstOne(num []uint8) (int64, error) {
	count := int64(1)
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

//findCount
func FindCount(bytes []uint8, count int32) ([]int64, bool) {
	var res []int64
	num := int64(1)
	k := int32(0)
	for i := 0; i < len(bytes); i++ {
		x := uint8(128)
		temp := bytes[i]
		for x != 0 {
			if temp&x != 0 {
				res = append(res, num)
				k++
				if k == count {
					return res, true
				}
			}
			num++
			x = x >> 1
		}
	}
	return nil, false
}

//countOne 帮助计算余量 可以优化成hashmap查询
func CountOne(num []uint8) uint {
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
