package utils

import "time"

// GetTimeDifference retorna a diferença de tempo em horas entre dois timestamps ISO-8601
func GetTimeDifference(start, end string) (float64, error) {
	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return 0, err
	}

	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return 0, err
	}

	diff := endTime.Sub(startTime).Hours()
	return diff, nil
}
