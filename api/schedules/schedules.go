package schedules

import (
	"github.com/gin-gonic/gin"
	"github.com/markgerald/flyroutes/domain"
	"github.com/markgerald/flyroutes/models"
	"github.com/markgerald/flyroutes/services"
	"net/http"
)

func GetSchedule(c *gin.Context) {
	var requestData models.RequestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request data"})
		return
	}
	schedule, err := domain.GetScheduleData(requestData.Departure, requestData.Arrival)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	allRoutes, err := domain.GetRoutesData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	flightService := services.FlightService{}
	filteredFlights := flightService.FilterFlights(schedule, requestData, allRoutes)
	if filteredFlights == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No flights found"})
		return
	}

	c.JSON(http.StatusOK, filteredFlights)
}
