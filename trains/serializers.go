package trains

import (
	"time"

	"github.com/gin-gonic/gin"
)

// import (
// 	"github.com/gin-gonic/gin"
// )

//TicketResponse 单种票的响应格式
type TicketResponse struct {
	TrainNumber           string `json:"train_number"`
	LeaveStation          string `json:"leave_station"`
	ArrivalStation        string `json:"arrival_station"`
	LeaveTime             string `json:"leave_time"`
	ArrivalTime           string `json:"arrival_time"`
	TravelTime            string `json:"travel_time"` //耗时,"7小时13分"
	LeaveStationType      string `json:"leave_station_type"`
	ArrivalStationType    string `json:"arrival_station_type"`
	TrainType             string `json:"train_type"` // 列车类型
	BusinessSeatsNumber   uint   `json:"business_seats_number"`
	FirstSeatsNumber      uint   `json:"first_seats_number"`
	SecondSeatsNumber     uint   `json:"second_seats_number"`
	NoSeatsNumber         uint   `json:"no_seats_number"`
	HardSeatsNumber       uint   `json:"hard_seats_number"`
	HardBerthNumber       uint   `json:"hard_berth_number"`
	SoftBerthNumber       uint   `json:"soft_berth_number"`
	SeniorSoftBerthNumber uint   `json:"senior_soft_berth_number"`
}

//TicketsSerializer 所有票的序列化器
type TicketsSerializer struct {
	C               *gin.Context
	TripStaionPairs []TrainStaionPair
	Date            string
}

func (s *TicketsSerializer) Response() []TicketResponse {
	response := []TicketResponse{}
	for _, tripStaionPair := range s.TripStaionPairs {
		//用于计算travel时间
		leaveTime, _ := time.ParseInLocation("15:04:05", tripStaionPair.StartTime, time.Local)
		arrivalTime, _ := time.ParseInLocation("15:04:05", tripStaionPair.EndTime, time.Local)
		//计算过站类型
		var leaveStationType, arriveStationType string
		if tripStaionPair.StartStationNo == uint(1) {
			leaveStationType = "始"
		} else {
			leaveStationType = "过"
		}
		if tripStaionPair.EndStationNo == tripStaionPair.TrainStationNums {
			arriveStationType = "终"
		} else {
			arriveStationType = "过"
		}
		//查询余量
		trip := Trip{tripStaionPair.TrainID + "-" + s.Date}
		remainSeats := trip.getRemainSeats(tripStaionPair.StartStationNo, tripStaionPair.EndStationNo)
		//序列化每类Ticket
		temp := TicketResponse{
			TrainNumber:           tripStaionPair.TrainID,
			LeaveStation:          tripStaionPair.StartStationName,
			ArrivalStation:        tripStaionPair.EndStationName,
			LeaveTime:             tripStaionPair.StartTime,
			ArrivalTime:           tripStaionPair.EndTime,
			TravelTime:            arrivalTime.Sub(leaveTime).String(),
			LeaveStationType:      leaveStationType,
			ArrivalStationType:    arriveStationType,
			TrainType:             tripStaionPair.TrainType,
			BusinessSeatsNumber:   (*remainSeats)["BusinessSeat"],
			FirstSeatsNumber:      (*remainSeats)["FirstSeat"],
			SecondSeatsNumber:     (*remainSeats)["SecondSeat"],
			NoSeatsNumber:         (*remainSeats)["NoSeat"],
			HardSeatsNumber:       (*remainSeats)["HardSeats"],
			HardBerthNumber:       (*remainSeats)["HardBerth"],
			SoftBerthNumber:       (*remainSeats)["SoftBerth"],
			SeniorSoftBerthNumber: (*remainSeats)["SeniorSoftBerth"],
		}
		response = append(response, temp)
	}
	return response
}

//TrainStaionResponse 某个车次经过的某个站点的响应格式
type TrainStaionResponse struct {
	StationNo   uint   `json:"station_no"`
	StationName string `json:"station_name"`
	ArrivalTime string `json:"arrival_time"`
	LeaveTime   string `json:"leave_time"`
	WaitTime    string `json:"wait_time"`
}

//TrainStaionSerializer 某个车次经过的所有站点的序列化器
type TrainStaionSerializer struct {
	C            *gin.Context
	TrainStaions []TrainStaion
	IsToday      bool
}

func (s *TrainStaionSerializer) Response() []TrainStaionResponse {
	response := []TrainStaionResponse{}
	for _, trainStaion := range s.TrainStaions {
		if s.IsToday {
			trainStaion.ArriveTime = trainStaion.TodayArriveTime
			trainStaion.DepartureTime = trainStaion.TodayDepartureTime
		}
		arrivalTime, _ := time.ParseInLocation("15:04:05", trainStaion.ArriveTime, time.Local)
		leaveTime, _ := time.ParseInLocation("15:04:05", trainStaion.DepartureTime, time.Local)
		waitTime := leaveTime.Sub(arrivalTime).String()
		temp := TrainStaionResponse{
			StationNo:   trainStaion.StationNo,
			StationName: trainStaion.StationName,
			ArrivalTime: trainStaion.ArriveTime,
			LeaveTime:   trainStaion.DepartureTime,
			WaitTime:    waitTime,
		}
		response = append(response, temp)
	}
	return response
}
