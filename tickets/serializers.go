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
	TripStaionPairs []TripStaionPair
}

func (s *TicketListSerializer) Response() []TicketListResponse {
	response := []TicketListResponse{}
	for _, tripStaionPair := range s.TripStaionPairs {
		leaveTime, _ := time.ParseInLocation("2006-01-02 15:04:05", tripStaionPair.StartTime, time.Local)
		arriveTime, _ := time.ParseInLocation("2006-01-02 15:04:05", tripStaionPair.EndTime, time.Local)
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
		tripID := Trip{tripStaionPair.TripID}
		temp := TicketListResponse{
			TrainNumber:         tripStaionPair.TrainNumber,
			LeaveStation:        tripStaionPair.StartStationName,
			ArriveStation:       tripStaionPair.EndStationName,
			LeaveTime:           tripStaionPair.StartTime,
			ArrivalTime:         tripStaionPair.EndTime,
			TravelTime:          arriveTime.Sub(leaveTime).String(),
			LeaveStationType:    leaveStationType,
			ArriveStationType:   arriveStationType,
			TrainType:           tripStaionPair.TrainType,
			BusinessSeatsNumber: tripID.getRemainSeats(tripStaionPair.StartStationNo, tripStaionPair.EndStationNo, "BusinessSeat"),
			// FirstSeatsNumber:      tripID.getRemainSeats(tripStaionPair.StartStationNo, tripStaionPair.EndStationNo, "FirstSeat"),
			// SecondSeatsNumber:     tripID.getRemainSeats(tripStaionPair.StartStationNo, tripStaionPair.EndStationNo, "SecondSeat"),
			// NoSeatsNumber:         tripID.getRemainSeats(tripStaionPair.StartStationNo, tripStaionPair.EndStationNo, "NoSeat"),
			// HardSeatsNumber:       tripID.getRemainSeats(tripStaionPair.StartStationNo, tripStaionPair.EndStationNo, "HardSeats"),
			// HardBerthNumber:       tripID.getRemainSeats(tripStaionPair.StartStationNo, tripStaionPair.EndStationNo, "HardBerth"),
			// SoftBerthNumber:       tripID.getRemainSeats(tripStaionPair.StartStationNo, tripStaionPair.EndStationNo, "SoftBerth"),
			// SeniorSoftBerthNumber: tripID.getRemainSeats(tripStaionPair.StartStationNo, tripStaionPair.EndStationNo, "SeniorSoftBerth"),
		}
		response = append(response, temp)
	}
	return response
}
