package utils

import (
	"errors"
	"log"
	"math/rand"
	"server/model/join"
	"server/model/table"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "2006-01-02 03:04"
	dFormat    = "2006-01-02"
	decimal    = 10
	letters    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	location   = "Asia/Tokyo"
)

func FormatTimeStamp(date string) (time.Time, error) {
	var (
		t     time.Time
		times []int
	)
	loc, _ := time.LoadLocation(location)

	t, err := time.ParseInLocation(dateFormat, date, loc)
	if err == nil {
		return t, nil
	}

	t, err = time.ParseInLocation(dFormat, date, loc)
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

func WorkLogMessage(workTimes []table.WorkTimes) (string, error) {
	var message string

	for _, w := range workTimes {
		message += strconv.FormatUint(w.Id, decimal) + " "
		message += w.Content + "\n"
	}

	return message, nil
}

func WorkTimeMessage(sumTimes time.Duration, date time.Time) string {
	return date.Format("2006-01-02") + " の作業時間は" + sumTimes.String() + "です"
}

func CalcWorkTimes(workTimes []table.WorkTimes, workRests []table.WorkRests, date time.Time) (time.Duration, error) {
	var times time.Duration
	month := date.Month()
	year := date.Year()
	day := date.Day()
	loc, _ := time.LoadLocation(location)

	for _, w := range workTimes {

		if !(w.StartedAt.Year() == year && w.StartedAt.Month() == month && w.StartedAt.Day() == day) {
			w.StartedAt = time.Date(year, month, day, 0, 0, 0, 0, loc)
		}
		if !(w.FinishedAt.Year() == year && w.FinishedAt.Month() == month && w.FinishedAt.Day() == day) {
			w.FinishedAt = time.Date(year, month, day, 0, 0, 0, 0, loc).Add(24 * time.Hour)
		}
		times += w.FinishedAt.Sub(w.StartedAt)
	}

	for _, w := range workRests {

		if !(w.StartedAt.Year() == year && w.StartedAt.Month() == month && w.StartedAt.Day() == day) {
			w.StartedAt = time.Date(year, month, day, 0, 0, 0, 0, loc)
		}
		if !(w.FinishedAt.Year() == year && w.FinishedAt.Month() == month && w.FinishedAt.Day() == day) {
			w.FinishedAt = time.Date(year, month, day, 0, 0, 0, 0, loc).Add(24 * time.Hour)
		}
		times -= w.FinishedAt.Sub(w.StartedAt)
	}

	return times, nil
}

func SessionWorkTimesSetHash(session []table.SessionWorkTimes) []table.SessionWorkTimes {
	for i, _ := range session {
		session[i].Hash = createHash(30)
	}

	return session
}

func createHash(cnt int) string {
	b := make([]byte, cnt)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
