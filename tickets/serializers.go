package tickets

import (
	"time"

	"github.com/gin-gonic/gin"
)

// import (
// 	"github.com/gin-gonic/gin"
// )

type TicketResponse struct {
	TripID                uint      `json:"trainNumber"`
	StartStationNo        uint      `json:"startStationNo"`
	EndStationNo          uint      `json:"endStationNo"`
	startStation          string    `json:"startStation"`
	endStation            string    `json:"endStation"`
	startTime             time.Time `json:"startTime"`
	arrivalTime           time.Time `json:"arrivalTime"`
	startStationType      uint      `json:"startStationType"`
	endStationType        uint      `json:"endStationType"`
	businessSeatsNumber   uint      `json:"businessSeatsNumber"`
	firstSeatsNumber      bool      `json:"firstSeatsNumber"`
	secondSeatsNumber     uint      `json:"secondSeatsNumber"`
	hardSeatsNumber       uint      `json:"hardSeatsNumber"`
	hardBerthNumber       uint      `json:"hardBerthNumber"`
	softBerthNumber       uint      `json:"softBerthNumber"`
	seniorSoftBerthNumber uint      `json:"seniorSoftBerthNumber"`
}
type TicketsSerializer struct {
	C       *gin.Context
	Tickets []Ticket
}

func (s *TicketsSerializer) Response() []TicketResponse {
	response := []TicketResponse{}
	for _, ticket := range s.Tickets {
		temp := TicketResponse{
			TripID:         ticket.TripID,
			StartStationNo: ticket.StartStationNo,
			EndStationNo:   ticket.EndStationNo,
			startStation:   ticket.StartStation,
			endStation:     ticket.EndStation,
		}
		response = append(response, temp)
	}
	return response
}
