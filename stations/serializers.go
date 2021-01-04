package stations

import (
	"github.com/gin-gonic/gin"
)

// "hot_cities": [{//热门城市列表
// 	"id": 1,
// 	"spell": "beijing",
// 	"name": "北京"
//   }],
//  "cities": {//城市列表
// 	 "A": [{//首字母
// 		 "id": 56,
// 		 "spell": "aba",
// 		 "name": "阿坝"
// 	  },{
// 		 "id": 57,
// 		 "spell": "akesu",
// 		 "name": "阿克苏"
// 		 },],
// 	"B":[{}]
//
//还没按字母分；

// StationListResponse
type StationListResponse struct {
	HotCities []Station `json:"hot_cities"`
	Cities    []Station `json:"cities"`
}

// StationListSerializer
type StationListSerializer struct {
	C           *gin.Context
	Stations    []Station
	HotStations []Station
}

//Response
func (s *StationListSerializer) Response() StationListResponse {
	response := StationListResponse{}
	response.Cities = s.Stations
	response.HotCities = s.HotStations
	return response
}
