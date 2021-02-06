package tickets

import (
	"time"

	"github.com/gin-gonic/gin"
)

// import (
// 	"github.com/gin-gonic/gin"
// )

//TicketListResponse 查票的响应格式
type TicketListResponse struct {
	TrainNumber           string `json:"train_number"` // 列车字符串
	LeaveStation          string `json:"leave_station"`
	ArriveStation         string `json:"arrival_station"`
	LeaveTime             string `json:"leave_time"`
	ArrivalTime           string `json:"arrival_time"`
	TravelTime            string `json:"travel_time"` //耗时,"7小时13分"
	LeaveStationType      string `json:"leave_station_type"`
	ArriveStationType     string `json:"arrival_station_type"`
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
type TicketListSerializer struct {
	C               *gin.Context
	TripStaionPairs []TrainStaionPair
	date            string
}

func (s *TicketListSerializer) Response() []TicketListResponse {
	response := []TicketListResponse{}
	for _, tripStaionPair := range s.TripStaionPairs {
		leaveTime, _ := time.ParseInLocation("15:04:05", tripStaionPair.StartTime, time.Local)
		arriveTime, _ := time.ParseInLocation("15:04:05", tripStaionPair.EndTime, time.Local)
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

		trip := Trip{tripStaionPair.TrainID + "-" + s.date}
		remainSeats := trip.getRemainSeats(tripStaionPair.StartStationNo, tripStaionPair.EndStationNo)

		temp := TicketListResponse{
			TrainNumber:           tripStaionPair.TrainID,
			LeaveStation:          tripStaionPair.StartStationName,
			ArriveStation:         tripStaionPair.EndStationName,
			LeaveTime:             tripStaionPair.StartTime,
			ArrivalTime:           tripStaionPair.EndTime,
			TravelTime:            arriveTime.Sub(leaveTime).String(),
			LeaveStationType:      leaveStationType,
			ArriveStationType:     arriveStationType,
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
