package schedules

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/markgerald/flyroutes/models"
	"github.com/markgerald/flyroutes/services"
	"io"
	"net/http"
	"strconv"
	"time"
)

var (
	scheduleURL = "https://services-api.ryanair.com/timtbl/3/schedules/"
	routesUrl   = "https://services-api.ryanair.com/views/locate/3/routes"
)

func GetSchedule(c *gin.Context) {
	var requestData models.RequestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request data"})
		return
	}

	year, month, _ := time.Now().Date()
	url := scheduleURL +
		requestData.Departure + "/" +
		requestData.Arrival + "/years/" +
		strconv.Itoa(year) + "/months/" + strconv.Itoa(int(month))
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var schedule models.Schedule
	err = json.Unmarshal(body, &schedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON response"})
		return
	}

	respRoutes, err := http.Get(routesUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer respRoutes.Body.Close()
	bodyRoutes, err := io.ReadAll(respRoutes.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read routes response body"})
		return
	}

	var allRoutes []models.Route
	err = json.Unmarshal(bodyRoutes, &allRoutes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse routes JSON response"})
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
