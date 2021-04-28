package stations

import (
	"github.com/gin-gonic/gin"
)

//还没按字母分；

type StationListResponse struct {
	HotCities []Station `json:"hot_cities"`
	Cities    []Station `json:"cities"`
}

type StationListSerializer struct {
	C           *gin.Context
	Stations    []Station
	HotStations []Station
}

func (s *StationListSerializer) Response() StationListResponse {
	response := StationListResponse{}
	response.Cities = s.Stations
	response.HotCities = s.HotStations
	return response
}
