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
		temp := TicketListResponse{
			TripID:         ticket.TripID,
			StartStationNo: ticket.StartStationNo,
			EndStationNo:   ticket.EndStationNo,
			// StartStation:          ticket.StartStation,
			// EndStation:            ticket.EndStation,
			BusinessSeatsNumber:   ticket.getRemainSeats("BusinessSeat"),
			FirstSeatsNumber:      ticket.getRemainSeats("FirstSeat"),
			SecondSeatsNumber:     ticket.getRemainSeats("SecondSeat"),
			HardSeatsNumber:       ticket.getRemainSeats("HardSeats"),
			HardBerthNumber:       ticket.getRemainSeats("HardBerth"),
			SoftBerthNumber:       ticket.getRemainSeats("SoftBerth"),
			SeniorSoftBerthNumber: ticket.getRemainSeats("SeniorSoftBerth"),
		}
		response = append(response, temp)
	}
	return response
}
