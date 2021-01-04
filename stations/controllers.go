package stations

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//StationList
func StationList(c *gin.Context) {
	stationModels, err := FindStations()
	hotStationModels, err := FindHotStations()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": "Invalid param"})
		return
	}
	serializer := StationListSerializer{c, stationModels, hotStationModels}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "查找成功", "data": serializer.Response()})
}
