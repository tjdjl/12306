package trains

// import (
// 	"github.com/gin-gonic/gin"
// )
// "list": [
// 	{
// 		"train_number":"", // 列次 字符串
// 		 "start_station":"",
// 		 "end_station":"",
// 		 "start_time":"", // 出发时间
// 		 "arrival_time":"", // 达到时间
// 		 "start_station_type":"", //起始站类型, 是否始发站还是过站, 0, 1
// 		 "end_station_type":"", //到达站类型, 是否是过站还是终点站, 1,2
// 		 "train_type":"", // 列车类型,
// 		 "business_seats_number": , // 商务座余数
// 		 "first_seats_number": , // 一等座
// 		 "second_seats_number": , // 二等座
// 		 "no_seats_number": , //无座
// 		 "hard_seats_number": , // 硬座数量
// 		 "hard_berth_number": ,// 硬卧
// 		 "soft_berth_number": ,// 软卧
// 		 "senior_soft_berth_number": , // 高软
// 	}
//  ]
// type TrainValidator struct { //用来接收
// 	Article struct {
// 		StartCity string `form:"startCity" json:"startCity" binding:""`
// 		EndCity   string `form:"endCity" json:"endCity" binding:""`
// 		Date      string `form:"date" json:"date" binding:""`
// 		Type      string `form:"type" json:"type"`
// 	} `json:"article"`
// 	TrainModel TrainModel `json:"-"`
// }

// //找到车次id

// func NewArticleModelValidator() TrainValidator {
// 	return TrainValidator{}
// }

// func (s *TrainValidator) Bind(c *gin.Context) error {

// 	return nil
// }
