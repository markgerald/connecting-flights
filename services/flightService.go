package services

import (
	"encoding/json"
	"github.com/markgerald/flyroutes/models"
	"github.com/markgerald/flyroutes/utils"
	"io"
	"net/http"
	"strconv"
	"time"
)

type FlightService struct{}

func (fs *FlightService) FilterFlights(scheduleData models.Schedule, requestData models.RequestData, allRoutes []models.Route) []models.Flight {
	var filteredFlights []models.Flight

	departureDateTime, err := time.Parse(time.RFC3339, requestData.DepartureDateTime)
	if err != nil {
		return nil
	}

	arrivalDateTime, err := time.Parse(time.RFC3339, requestData.ArrivalDateTime)
	if err != nil {
		return nil
	}

	// Primeiro, lidamos com os voos diretos.
	for _, day := range scheduleData.Days {
		for _, flight := range day.Flights {
			currentDate := time.Date(
				departureDateTime.Year(),
				departureDateTime.Month(),
				day.Day,
				0, 0, 0, 0, time.UTC)
			fullDepartureDateTime := utils.StringToTime(flight.DepartureTime, currentDate)
			fullArrivalDateTime := utils.StringToTime(flight.ArrivalTime, currentDate)
			if fullDepartureDateTime.After(departureDateTime) && fullArrivalDateTime.Before(arrivalDateTime) {
				filteredFlight := models.Flight{
					Number: flight.Number,
					Stops:  0,
					Legs: []models.Leg{
						{
							DepartureAirport:  requestData.Departure,
							ArrivalAirport:    requestData.Arrival,
							DepartureDateTime: fullDepartureDateTime.Format(time.RFC3339),
							ArrivalDateTime:   fullArrivalDateTime.Format(time.RFC3339),
						},
					},
				}
				filteredFlights = append(filteredFlights, filteredFlight)
			}
		}
	}

	// Agora, lidamos com voos com conexões.
	for _, route := range allRoutes {
		if route.AirportFrom == requestData.Departure && route.ConnectingAirport != "" {
			connectingSchedule, err := fs.getConnectedFlightSchedule(route.AirportFrom, route.ConnectingAirport)
			if err != nil {
				continue
			}

			for _, connDay := range connectingSchedule.Days {
				for _, connFlight := range connDay.Flights {
					currentDate := time.Date(
						departureDateTime.Year(),
						departureDateTime.Month(),
						connDay.Day,
						0, 0, 0, 0, time.UTC)
					connDepartureTime := utils.StringToTime(connFlight.DepartureTime, currentDate)
					connArrivalTime := utils.StringToTime(connFlight.ArrivalTime, currentDate)

					if connArrivalTime.Add(2 * time.Hour).Before(connDepartureTime) {
						// Aqui é onde vamos incluir a parte de voo direto e a parte de conexão.
						filteredFlight := models.Flight{
							Number: connFlight.Number,
							Stops:  1,
							Legs: []models.Leg{
								{
									DepartureAirport:  requestData.Departure,
									ArrivalAirport:    route.ConnectingAirport,
									DepartureDateTime: connDepartureTime.Format(time.RFC3339),
									ArrivalDateTime:   connArrivalTime.Format(time.RFC3339),
								},
								{
									DepartureAirport:  route.ConnectingAirport,
									ArrivalAirport:    requestData.Arrival,                                     // Fazendo a suposição de que é a chegada final.
									DepartureDateTime: connArrivalTime.Add(2 * time.Hour).Format(time.RFC3339), // Supondo 2 horas de espera para o próximo voo.
									ArrivalDateTime:   connArrivalTime.Add(3 * time.Hour).Format(time.RFC3339), // Supondo que o próximo voo leva 1 hora.
								},
							},
						}
						filteredFlights = append(filteredFlights, filteredFlight)
					}
				}
			}
		}
	}

	return filteredFlights
}

func (fs *FlightService) getConnectedFlightSchedule(departure, arrival string) (models.Schedule, error) {
	year, month, _ := time.Now().Date()
	url := "https://services-api.ryanair.com/timtbl/3/schedules/" +
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
