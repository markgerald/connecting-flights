package domain

import (
	"encoding/json"
	"github.com/markgerald/flyroutes/models"
	"io"
	"net/http"
	"strconv"
	"time"
)

const ScheduleURL = "https://services-api.ryanair.com/timtbl/3/schedules/"

func GetScheduleData(departure, arrival string) (models.Schedule, error) {
	year, month, _ := time.Now().Date()
	url := ScheduleURL +
		departure + "/" +
		arrival + "/years/" +
		strconv.Itoa(year) + "/months/" + strconv.Itoa(int(month))
	resp, err := http.Get(url)
	if err != nil {
		return models.Schedule{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Schedule{}, err
	}

	var schedule models.Schedule
	err = json.Unmarshal(body, &schedule)
	if err != nil {
		return models.Schedule{}, err
	}

	return schedule, nil
}

func GetConnectedFlightSchedule(departure, arrival string) (models.Schedule, error) {
	year, month, _ := time.Now().Date()
	url := ScheduleURL +
		departure + "/" +
		arrival + "/years/" +
		strconv.Itoa(year) + "/months/" + strconv.Itoa(int(month))

	resp, err := http.Get(url)
	if err != nil {
		return models.Schedule{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Schedule{}, err
	}

	var schedule models.Schedule
	err = json.Unmarshal(body, &schedule)
	if err != nil {
		return models.Schedule{}, err
	}

	return schedule, nil
}
