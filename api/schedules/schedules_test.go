package schedules

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const routesUrl = "https://services-api.ryanair.com/views/locate/3/routes"
const scheduleURL = "https://services-api.ryanair.com/timtbl/3/schedules/"

func TestGetSchedule(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", scheduleURL+"DUB/STN/years/2023/months/10",
		httpmock.NewStringResponder(200, `{"month":10,"days":[{"day":1,"flights":[{"number":"FR123","departureTime":"12:00","arrivalTime":"14:00","stops":0,"legs":[]}]}]}`))
	httpmock.RegisterResponder("GET", routesUrl,
		httpmock.NewStringResponder(404, `{"error": "No flights found"}`))

	r := gin.Default()
	r.POST("/schedules", GetSchedule)

	body := `{"departure":"DUB", "arrival":"WRQ", "departureDateTime": "2023-10-19T15:00:00Z", "arrivalDateTime": "2024-12-19T18:00:00Z"}`
	req, err := http.NewRequest(http.MethodPost, "/schedules", bytes.NewBufferString(body))
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "FR123")
}
