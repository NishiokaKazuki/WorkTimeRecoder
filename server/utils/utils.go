package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "2006/01/02 12:34"
)

func FormatTimeStamp(date string) (time.Time, error) {
	var (
		t     time.Time
		times []int
	)

	t, err := time.Parse(dateFormat, date)
	if err == nil {
		return t, nil
	}

	if !strings.Contains(date, ":") {
		return t, errors.New("Cannot format timestamp.")
	}
	slice := strings.Split(date, ":")

	for _, s := range slice {
		timeSplit, err := strconv.Atoi(s)
		if err != nil {
			return t, errors.New("Cannot format timestamp.")
		}
		times = append(times, timeSplit)
	}

	for i := len(times); i < 4; i++ {
		times = append(times, 0)
	}

	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Tokyo")
	t = time.Date(now.Year(), now.Month(), now.Day(), times[0], times[1], times[2], 0, loc)

	return t, nil
}

func SplitTimeOption(message []string) (time.Time, bool) {
	var date time.Time

	for i, msg := range message {
		if msg == "-t" && len(message) >= i+2 {
			if date, err := FormatTimeStamp(message[i+1]); err == nil {
				return date, true
			}
		}
	}

	return date, false
}

func SplitSuppleOption(message []string) (string, bool) {

	for i, msg := range message {
		if msg == "-m" && len(message) >= i+2 {
			return message[i+1], true
		}
	}

	return "", false
}
