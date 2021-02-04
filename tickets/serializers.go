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
	// TripID                uint      `json:"trip_id"`
	// StartStationNo        uint      `json:"leave_station_no"`
	// EndStationNo          uint      `json:"arrival_station_no"`
	TrainNumber           string        `json:"train_number"` // 列车字符串
	LeaveStation          string        `json:"leave_station"`
	ArriveStation         string        `json:"arrival_station"`
	LeaveTime             time.Time     `json:"leave_time"`
	ArrivalTime           time.Time     `json:"arrival_time"`
	TravelTime            time.Duration `json:"travel_time"` //耗时,"7小时13分"
	LeaveStationType      string        `json:"leave_station_type"`
	ArriveStationType     string        `json:"arrival_station_type"`
	TrainType             string        `json:"train_type"` // 列车类型
	BusinessSeatsNumber   uint          `json:"business_seats_number"`
	FirstSeatsNumber      uint          `json:"first_seats_number"`
	SecondSeatsNumber     uint          `json:"second_seats_number"`
	NoSeatsNumber         uint          `json:"no_seats_number"`
	HardSeatsNumber       uint          `json:"hard_seats_number"`
	HardBerthNumber       uint          `json:"hard_berth_number"`
	SoftBerthNumber       uint          `json:"soft_berth_number"`
	SeniorSoftBerthNumber uint          `json:"senior_soft_berth_number"`
}
type TicketListSerializer struct {
	C                    *gin.Context
	TripStartAndEndSlice []TripStartNoAndEndNo
}

func (s *TicketListSerializer) Response() []TicketListResponse {
	response := []TicketListResponse{}
	for _, tripStartAndEnd := range s.TripStartAndEndSlice {
		train := tripStartAndEnd.getTrainDetail()
		station := tripStartAndEnd.getStationDetail()
		var leaveStationType, arriveStationType string
		if tripStartAndEnd.StartStationNo == uint(1) {
			leaveStationType = "始"
		} else {
			leaveStationType = "过"
		}
		if tripStartAndEnd.EndStationNo == train.Length {
			arriveStationType = "终"
		} else {
			arriveStationType = "过"
		}
		temp := TicketListResponse{
			TrainNumber:           train.TrainNumber,
			LeaveStation:          station[0].StationName,
			ArriveStation:         station[1].StationName,
			LeaveTime:             station[0].StationTime,
			ArrivalTime:           station[1].StationTime,
			TravelTime:            station[1].StationTime.Sub(station[0].StationTime),
			LeaveStationType:      leaveStationType,
			ArriveStationType:     arriveStationType,
			TrainType:             train.Catogory,
			BusinessSeatsNumber:   tripStartAndEnd.getRemainSeats("BusinessSeat"),
			FirstSeatsNumber:      tripStartAndEnd.getRemainSeats("FirstSeat"),
			SecondSeatsNumber:     tripStartAndEnd.getRemainSeats("SecondSeat"),
			NoSeatsNumber:         tripStartAndEnd.getRemainSeats("NoSeat"),
			HardSeatsNumber:       tripStartAndEnd.getRemainSeats("HardSeats"),
			HardBerthNumber:       tripStartAndEnd.getRemainSeats("HardBerth"),
			SoftBerthNumber:       tripStartAndEnd.getRemainSeats("SoftBerth"),
			SeniorSoftBerthNumber: tripStartAndEnd.getRemainSeats("SeniorSoftBerth"),
		}
		response = append(response, temp)
	}
	return response
}
