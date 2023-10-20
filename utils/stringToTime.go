package utils

import "time"

func StringToTime(timestring string, currentDate time.Time) time.Time {
	layout := "15:04"
	timeParsed, _ := time.Parse(layout, timestring)
	return time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), timeParsed.Hour(), timeParsed.Minute(), 0, 0, time.UTC)
}
