package services

import (
	"github.com/markgerald/flyroutes/domain"
	"github.com/markgerald/flyroutes/models"
	"github.com/markgerald/flyroutes/utils"
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

	for _, route := range allRoutes {
		if route.AirportFrom == requestData.Departure && route.ConnectingAirport != "" {
			connectingSchedule, err := domain.GetConnectedFlightSchedule(route.AirportFrom, route.ConnectingAirport)
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
									ArrivalAirport:    requestData.Arrival,
									DepartureDateTime: connArrivalTime.Add(2 * time.Hour).Format(time.RFC3339),
									ArrivalDateTime:   connArrivalTime.Add(3 * time.Hour).Format(time.RFC3339),
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
