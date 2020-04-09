package utils

import (
	"errors"
	"log"
	"server/model/join"
	"server/model/table"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "2006-01-02 03:04"
)

func FormatTimeStamp(date string) (time.Time, error) {
	var (
		t     time.Time
		times []int
	)
	loc, _ := time.LoadLocation("Asia/Tokyo")

	t, err := time.ParseInLocation(dateFormat, date, loc)
	log.Println(err)
	log.Println(t)
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
	t = time.Date(now.Year(), now.Month(), now.Day(), times[0], times[1], times[2], 0, loc)

	return t, nil
}

func SplitTimeOption(message []string) (time.Time, bool) {
	var date time.Time

	for i, msg := range message {
		if msg == "-t" && len(message) >= i+2 {
			if len(message) >= i+3 {
				if !strings.HasPrefix(message[i+2], "-") {
					message[i+1] += " " + message[i+2]
					log.Println(message[i+1])
				}
			}
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

func SplitWorkInfo(workInfo []join.WorkInfos, user table.Users) (string, error) {
	var message string
	log.Println(workInfo)

	message += "作業時間\n"
	for _, w := range workInfo {
		message += w.Content
		if w.Supplement != "" {
			message += w.Supplement
		}
		message += "\n"

		message += "開始時間 " + w.StartedAt.Format(dateFormat) + "\n" +
			"終了時間" + w.FinishedAt.Format(dateFormat) + "\n"

		message += "休憩時間\n"
		for _, r := range w.WorkRests {
			if r.WorkTimeId == w.Id {
				message += "開始時間 " + r.StartedAt.Format(dateFormat) + "\n" +
					"終了時間 " + r.FinishedAt.Format(dateFormat) + "\n"
			}
		}
	}

	return user.Name + " 作業記録\n" + message, nil
}
