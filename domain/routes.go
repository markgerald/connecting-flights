package domain

import (
	"encoding/json"
	"github.com/markgerald/flyroutes/models"
	"io"
	"net/http"
)

const RoutesUrl = "https://services-api.ryanair.com/views/locate/3/routes"

func GetRoutesData() ([]models.Route, error) {
	resp, err := http.Get(RoutesUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var allRoutes []models.Route
	err = json.Unmarshal(body, &allRoutes)
	if err != nil {
		return nil, err
	}

	return allRoutes, nil
}
