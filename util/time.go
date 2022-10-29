package util

import (
	"strconv"
	"strings"
)

// TimeFormat format time interval to human-readable
func TimeFormat(totalSeconds int) string {
	if totalSeconds <= 0 {
		return "0 second"
	}
	timeBuilder := strings.Builder{}

	hours := totalSeconds / 3600
	minutes := totalSeconds % 3600 / 60
	seconds := totalSeconds % 3600 % 60

	if hours > 0 {
		timeBuilder.WriteString(pluralize(hours, "hour", "hours"))
	}

	if minutes > 0 {
		if timeBuilder.Len() > 0 {
			timeBuilder.WriteString(", ")
		}
		timeBuilder.WriteString(pluralize(minutes, "minute", "minutes"))
	}

	if seconds > 0 {
		if timeBuilder.Len() > 0 {
			timeBuilder.WriteString(", ")
		}
		timeBuilder.WriteString(pluralize(seconds, "second", "seconds"))
	}

	return timeBuilder.String()
}

func pluralize(times int, label string, pluralLabel string) string {
	if times <= 0 {
		return "no " + pluralLabel
	}

	if times == 1 {
		return strconv.Itoa(times) + " " + label
	}

	return strconv.Itoa(times) + " " + pluralLabel
}
