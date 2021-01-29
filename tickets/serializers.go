package tickets

import (
	"time"

	"github.com/gin-gonic/gin"
)

// import (
// 	"github.com/gin-gonic/gin"
// )

type TicketListResponse struct {
	TripID                uint      `json:"trainNumber"`
	StartStationNo        uint      `json:"startStationNo"`
	EndStationNo          uint      `json:"endStationNo"`
	StartStation          string    `json:"startStation"`
	EndStation            string    `json:"endStation"`
	StartTime             time.Time `json:"startTime"`
	ArrivalTime           time.Time `json:"arrivalTime"`
	StartStationType      uint      `json:"startStationType"`
	EndStationType        uint      `json:"endStationType"`
	BusinessSeatsNumber   uint      `json:"businessSeatsNumber"`
	FirstSeatsNumber      uint      `json:"firstSeatsNumber"`
	SecondSeatsNumber     uint      `json:"secondSeatsNumber"`
	HardSeatsNumber       uint      `json:"hardSeatsNumber"`
	HardBerthNumber       uint      `json:"hardBerthNumber"`
	SoftBerthNumber       uint      `json:"softBerthNumber"`
	SeniorSoftBerthNumber uint      `json:"seniorSoftBerthNumber"`
}
type TicketListSerializer struct {
	C       *gin.Context
	Tickets []TripSeries
}

func (s *TicketListSerializer) Response() []TicketListResponse {
	response := []TicketListResponse{}
	for _, ticket := range s.Tickets {
		seatNumbers := ticket.getRemainSeats()
		temp := TicketListResponse{
			TripID:         ticket.TripID,
			StartStationNo: ticket.StartStationNo,
			EndStationNo:   ticket.EndStationNo,
			// StartStation:          ticket.StartStation,
			// EndStation:            ticket.EndStation,
			BusinessSeatsNumber:   seatNumbers.BusinessSeats,
			FirstSeatsNumber:      seatNumbers.FirstSeats,
			SecondSeatsNumber:     seatNumbers.SecondSeats,
			HardSeatsNumber:       seatNumbers.HardBerth,
			HardBerthNumber:       seatNumbers.HardBerth,
			SoftBerthNumber:       seatNumbers.SoftBerth,
			SeniorSoftBerthNumber: seatNumbers.SeniorSoftBerth,
		}
		response = append(response, temp)
	}
	return response
}
