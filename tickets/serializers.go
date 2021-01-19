package trains

import (
	"github.com/gin-gonic/gin"
	"time"
)

// import (
// 	"github.com/gin-gonic/gin"
// )

type TicketResponse struct {
	ID             uint                  `json:"-"`
	TripID          string           `json:"trainNumber"`
	startStation           string                `json:"startStation"`
	endStation    string                `json:"endStation"`
	startTime            time.Time                `json:"startTime"`
	arrivalTime       time.Time                `json:"arrivalTime"`
	startStationType      uint                `json:"startStationType"`
	endStationType       uint   `json:"endStationType"`
	Tags           []string              `json:"businessSeatsNumber"`
	firstSeatsNumber       bool                  `json:"firstSeatsNumber"`
	secondSeatsNumber uint                  `json:"secondSeatsNumber"`
	hardSeatsNumber uint                  `json:"hardSeatsNumber"`
	hardBerthNumber uint                  `json:"hardBerthNumber"`
	softBerthNumber uint                  `json:"softBerthNumber"`
	seniorSoftBerthNumber uint                  `json:"seniorSoftBerthNumber"`
}
type TicketsSerializer struct {
	C *gin.Context
	Tickets []Ticket
}

func (s *TicketsSerializer) Response() []TicketResponse {
	response := []TicketResponse{}
	for _, ticket := range s.Tickets {
		temp  := TicketResponse{
			ID:          s.ID,
			Slug:        slug.Make(s.Title),
			Title:       s.Title,
			Description: s.Description,
			Body:        s.Body,
			CreatedAt:   s.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
			//UpdatedAt:      s.UpdatedAt.UTC().Format(time.RFC3339Nano),
			UpdatedAt:      s.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
			Author:         authorSerializer.Response(),
			Favorite:       s.isFavoriteBy(GetArticleUserModel(myUserModel)),
			FavoritesCount: s.favoritesCount(),
		}
		response = append(response, temp)
	}
	return response
}
